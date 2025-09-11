package kafka

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
)

type (
	handlerFunc func(ctx context.Context, msg Message) error
	handler     struct {
		executor handlerFunc
		name     string
	}
	Message struct {
		Data      []byte
		Topic     string
		Offset    int64
		Partition int32
		Timestamp time.Time
		Id        string
	}
	Consumer struct {
		appname         string
		c               *kafka.Consumer
		l               *slog.Logger
		latency         *prometheus.HistogramVec // Tracks how long it takes to process a message.
		connectionState *prometheus.GaugeVec     // Reports connection state: 1 = connected, 0 = disconnected
		reconnectCount  *prometheus.CounterVec   // Counts reconnection attempts
		handlers        map[string]handler
		servers         []string
		username        string
		password        string
		mu              sync.RWMutex
		connected       bool
		maxRetries      int           // Maximum number of reconnection attempts (0 = infinite)
		interval        time.Duration // Delay between reconnection attempts
	}
)

func NewConsumer(
	servers []string,
	username, password, appname string,
	l *slog.Logger, maxRetries int, interval time.Duration) (*Consumer, error) {

	client := &Consumer{
		l: l,
		latency: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "kafka_consumer_duration_seconds",
				Help:    "Duration of kafka consumers",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "error"},
		),
		connectionState: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "kafka_connection_state",
				Help: "Kafka connection state: 1 = connected, 0 = disconnected",
			},
			[]string{"appname"},
		),
		reconnectCount: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "kafka_reconnect_attempts_total",
				Help: "Total number of Kafka reconnection attempts",
			},
			[]string{"appname", "status"},
		),
		appname:    appname,
		handlers:   make(map[string]handler),
		servers:    servers,
		username:   username,
		password:   password,
		maxRetries: maxRetries,
		interval:   interval,
		connected:  false,
	}
	prometheus.MustRegister(client.latency, client.connectionState, client.reconnectCount)

	if err := client.connect(); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Consumer) connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Close existing connection if any
	if c.c != nil {
		c.c.Close()
	}

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":    strings.Join(c.servers, ","),
		"group.id":             c.appname,
		"enable.auto.commit":   false,
		"auto.offset.reset":    "latest", // or "earliest"
		"enable.partition.eof": false,
		"session.timeout.ms":   10000,
		"max.poll.interval.ms": 300000,
		"fetch.min.bytes":      1,
		"security.protocol":    "SASL_PLAINTEXT",
		"sasl.mechanism":       "PLAIN",
		"sasl.username":        c.username,
		"sasl.password":        c.password,
	})
	if err != nil {
		c.connected = false
		c.connectionState.WithLabelValues(c.appname).Set(0)
		return err
	}

	c.c = consumer
	c.connected = true
	c.connectionState.WithLabelValues(c.appname).Set(1)
	c.l.InfoContext(context.Background(), "kafka connected successfully", "appname", c.appname)

	// Re-subscribe to all topics
	for topic := range c.handlers {
		if err := c.c.SubscribeTopics([]string{topic}, nil); err != nil {
			c.l.ErrorContext(context.Background(), "failed to resubscribe to topic", "topic", topic, "error", err)
		}
	}

	return nil
}

func (c *Consumer) SubscribeTopics(topic, name string, h handlerFunc) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.handlers[topic]; ok {
		return fmt.Errorf("topic %s has already a registered handler", topic)
	}

	if err := c.c.SubscribeTopics([]string{topic}, nil); err != nil {
		return err
	}

	c.handlers[topic] = handler{name: name, executor: h}
	return nil
}

func (c *Consumer) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.c == nil {
		return nil
	}

	if err := c.c.Close(); err != nil {
		return err
	}
	for k := range c.handlers {
		delete(c.handlers, k)
	}
	c.connected = false
	c.connectionState.WithLabelValues(c.appname).Set(0)

	return nil
}

func (c *Consumer) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected
}

func (c *Consumer) Listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			c.l.InfoContext(ctx, "kafka consumer shutdown signal received, cleaning up...")
			if err := c.Close(); err != nil {
				c.l.ErrorContext(ctx, "an error occurred while closing the consumer", "error", err)
			}
			return

		default:
			if !c.IsConnected() {
				c.l.WarnContext(ctx, "kafka not connected, attempting to reconnect...")
				if err := c.reconnect(); err != nil {
					c.l.ErrorContext(ctx, "failed to reconnect", "error", err)
					time.Sleep(5 * time.Second)
					continue
				}
			}

			msg, err := c.c.ReadMessage(30 * time.Second)
			if err != nil {
				if c.shouldReconnectOnError(ctx, err) {
					c.l.WarnContext(ctx, "kafka connection lost, attempting to reconnect...")
					if err := c.reconnect(); err != nil {
						c.l.ErrorContext(ctx, "failed to reconnect", "error", err)
						time.Sleep(5 * time.Second)
					}
				}
				continue
			}

			if msg == nil {
				c.l.ErrorContext(ctx, "received an empty message")
				continue
			}

			c.handleMessage(ctx, msg)
		}
	}
}

func (c *Consumer) handleMessage(ctx context.Context, msg *kafka.Message) {
	if msg.TopicPartition.Topic == nil {
		c.l.ErrorContext(ctx, "received message with no topic", "message", msg)
		return
	}
	topic := *msg.TopicPartition.Topic

	h, ok := c.handlers[topic]
	if !ok {
		c.l.ErrorContext(ctx, "no handler registered for topic", "topic", topic)
		return
	}

	start := time.Now()
	rid := uuid.NewString()
	ctx = context.WithValue(ctx, "request_id", rid)
	ctx = context.WithValue(ctx, "method", h.name)
	var err error

	defer func() {
		c.latency.WithLabelValues(h.name, strconv.FormatBool(err != nil)).Observe(time.Since(start).Seconds())
	}()

	m := Message{
		Data:      msg.Value,
		Topic:     topic,
		Offset:    int64(msg.TopicPartition.Offset),
		Partition: msg.TopicPartition.Partition,
		Timestamp: msg.Timestamp,
		Id:        fmt.Sprintf("%d.%d", msg.TopicPartition.Partition, msg.TopicPartition.Offset),
	}
	if err = h.executor(ctx, m); err != nil {
		c.l.ErrorContext(ctx, "handler error", "topic", topic, "offset", msg.TopicPartition.Offset, "error", err)
		return
	}

	_, err = c.c.CommitMessage(msg)
	if err != nil {
		c.l.ErrorContext(ctx, "failed to commit message", "topic", topic, "offset", msg.TopicPartition.Offset, "error", err)
		return
	}

	c.l.InfoContext(ctx, "committed message", "topic", topic, "offset", msg.TopicPartition.Offset)
}

// shouldReconnectOnError determines if an error indicates a connection issue that requires reconnection
func (c *Consumer) shouldReconnectOnError(ctx context.Context, err error) bool {
	kerr, ok := err.(kafka.Error)
	if !ok {
		c.l.ErrorContext(ctx, "non-kafka error", "error", err)
		return true // Assume connection issue for non-kafka errors
	}

	// Check for timeout errors (these are normal and don't require reconnection)
	if kerr.Code() == kafka.ErrTimedOut {
		c.l.DebugContext(ctx, "poll timeout", "error", err, "hint", "no new messages were received, or the connection was interrupted")
		return false
	}

	// Check for specific error codes that indicate connection issues
	switch kerr.Code() {
	case kafka.ErrNetworkException:
		c.l.ErrorContext(ctx, "network exception", "error", err)
		return true
	case kafka.ErrBrokerNotAvailable:
		c.l.ErrorContext(ctx, "broker not available", "error", err)
		return true
	case kafka.ErrAllBrokersDown:
		c.l.ErrorContext(ctx, "all brokers down", "error", err)
		return true
	case kafka.ErrAuthentication:
		c.l.ErrorContext(ctx, "authentication failed", "error", err)
		return true
	case kafka.ErrInvalidSessionTimeout:
		c.l.ErrorContext(ctx, "invalid session timeout", "error", err)
		return true
	case kafka.ErrOffsetOutOfRange:
		c.l.ErrorContext(ctx, "offset out of range", "error", err)
		return true
	case kafka.ErrUnknownTopicOrPart:
		c.l.ErrorContext(ctx, "unknown topic or partition", "error", err)
		return true
	case kafka.ErrInvalidMsg:
		c.l.ErrorContext(ctx, "invalid message", "error", err)
		return true
	case kafka.ErrInvalidMsgSize:
		c.l.ErrorContext(ctx, "invalid message size", "error", err)
		return true
	case kafka.ErrInvalidPartitions:
		c.l.ErrorContext(ctx, "invalid partitions", "error", err)
		return true
	case kafka.ErrInvalidReplicationFactor:
		c.l.ErrorContext(ctx, "invalid replication factor", "error", err)
		return true
	case kafka.ErrInvalidReplicaAssignment:
		c.l.ErrorContext(ctx, "invalid replica assignment", "error", err)
		return true
	case kafka.ErrInvalidConfig:
		c.l.ErrorContext(ctx, "invalid config", "error", err)
		return true
	case kafka.ErrNotController:
		c.l.ErrorContext(ctx, "not controller", "error", err)
		return true
	case kafka.ErrInvalidRequiredAcks:
		c.l.ErrorContext(ctx, "invalid required acks", "error", err)
		return true
	case kafka.ErrIllegalGeneration:
		c.l.ErrorContext(ctx, "illegal generation", "error", err)
		return true
	case kafka.ErrInconsistentGroupProtocol:
		c.l.ErrorContext(ctx, "inconsistent group protocol", "error", err)
		return true
	case kafka.ErrInvalidGroupID:
		c.l.ErrorContext(ctx, "invalid group id", "error", err)
		return true
	case kafka.ErrUnknownMemberID:
		c.l.ErrorContext(ctx, "unknown member id", "error", err)
		return true
	case kafka.ErrRebalanceInProgress:
		c.l.WarnContext(ctx, "rebalance in progress", "error", err)
		return false // This is temporary, don't reconnect
	case kafka.ErrInvalidCommitOffsetSize:
		c.l.ErrorContext(ctx, "invalid commit offset size", "error", err)
		return true
	case kafka.ErrTopicAuthorizationFailed:
		c.l.ErrorContext(ctx, "topic authorization failed", "error", err)
		return true
	case kafka.ErrGroupAuthorizationFailed:
		c.l.ErrorContext(ctx, "group authorization failed", "error", err)
		return true
	case kafka.ErrClusterAuthorizationFailed:
		c.l.ErrorContext(ctx, "cluster authorization failed", "error", err)
		return true
	case kafka.ErrInvalidTimestamp:
		c.l.ErrorContext(ctx, "invalid timestamp", "error", err)
		return true
	case kafka.ErrUnsupportedSaslMechanism:
		c.l.ErrorContext(ctx, "unsupported sasl mechanism", "error", err)
		return true
	case kafka.ErrIllegalSaslState:
		c.l.ErrorContext(ctx, "illegal sasl state", "error", err)
		return true
	case kafka.ErrUnsupportedVersion:
		c.l.ErrorContext(ctx, "unsupported version", "error", err)
		return true
	case kafka.ErrTopicException:
		c.l.ErrorContext(ctx, "topic exception", "error", err)
		return true
	case kafka.ErrRecordListTooLarge:
		c.l.ErrorContext(ctx, "record list too large", "error", err)
		return true
	case kafka.ErrNotEnoughReplicas:
		c.l.ErrorContext(ctx, "not enough replicas", "error", err)
		return true
	case kafka.ErrNotEnoughReplicasAfterAppend:
		c.l.ErrorContext(ctx, "not enough replicas after append", "error", err)
		return true
	default:
		c.l.ErrorContext(ctx, "unknown error", "error", err, "code", kerr.Code())
		return true // Assume connection issue for unknown errors
	}
}

func (c *Consumer) reconnect() error {
	ctx := context.Background()
	attempt := 0
	for {
		attempt++

		// Check if we've exceeded max retries
		if c.maxRetries > 0 && attempt > c.maxRetries {
			c.l.ErrorContext(ctx, "max reconnection attempts exceeded", "attempts", attempt-1)
			c.reconnectCount.WithLabelValues(c.appname, "failed").Inc()
			return fmt.Errorf("max reconnection attempts (%d) exceeded", c.maxRetries)
		}

		c.l.InfoContext(ctx, "attempting to reconnect", "attempt", attempt, "interval", c.interval)
		c.reconnectCount.WithLabelValues(c.appname, "attempt").Inc()

		if err := c.connect(); err != nil {
			c.l.ErrorContext(ctx, "reconnection failed", "attempt", attempt, "error", err)
			time.Sleep(c.interval)
			continue
		}

		c.l.InfoContext(ctx, "reconnected successfully", "attempt", attempt)
		c.reconnectCount.WithLabelValues(c.appname, "success").Inc()
		return nil
	}
}
