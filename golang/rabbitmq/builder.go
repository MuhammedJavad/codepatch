package rabbitmq

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	amqp "github.com/rabbitmq/amqp091-go"
)

// New creates a new Manager instance
func New(
	host, port, vh, username, password, appname, hostname string,
	publisherChannelPoolSize int, logger *slog.Logger) *Client {

	appName = appname
	hostName = hostname
	cs := fmt.Sprintf("amqp://%s:%s@%s:%s/%s", username, password, host, port, vh)

	c := &Client{
		connectionString: cs,
		exchanges:        make(map[string]*exchange),
		logger:           logger,
		channels:         make(chan *amqp.Channel, publisherChannelPoolSize),
		channelLock:      &sync.Mutex{},
		latency: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "rabbitmq_consumer_duration_seconds",
				Help:    "Duration of rabbitmq handlers",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "error"},
		),
		connectionState: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "rabbitmq_connection_state",
				Help: "RabbitMQ connection state: 1 = connected, 0 = disconnected",
			},
			[]string{"appname"},
		),
		reconnectCount: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "rabbitmq_reconnect_attempts_total",
				Help: "Total number of RabbitMQ reconnection attempts",
			},
			[]string{"appname", "status"},
		),
	}
	prometheus.MustRegister(c.latency, c.connectionState, c.reconnectCount)
	return c
}

func (c *Client) DefineExchange(name, kind string, durable, autoDelete bool) *exchange {
	if e, ok := c.exchanges[name]; ok {
		return e
	}
	e := &exchange{
		name:       name,
		kind:       kind,
		durable:    durable,
		autoDelete: autoDelete,
		queues:     make(map[string]*queue),
	}
	c.exchanges[name] = e
	return e
}

func (e *exchange) DefineQueue(name, routingKey string, durable, autoDelete bool) *queue {
	if e.queues == nil {
		e.queues = make(map[string]*queue)
	}
	if q, ok := e.queues[name]; ok {
		return q
	}
	q := &queue{
		name:         name,
		routingKey:   routingKey,
		durable:      durable,
		autoDelete:   autoDelete,
		exchangeName: e.name,
		consumers:    make(map[string]*consumer),
		table:        make(map[string]interface{}),
	}
	e.queues[name] = q
	return q
}

func (e *exchange) DefineExclusiveQueue(name, routingKey string) *queue {
	queueName := fmt.Sprintf("%s-%s", name, hostName)
	q := e.DefineQueue(queueName, routingKey, false, false)
	q.exclusive = true
	return q
}

func (q *queue) WithDeadLetter(exchange, routingKey string) *queue {
	q.table["x-dead-letter-exchange"] = exchange
	q.table["x-dead-letter-routing-key"] = routingKey
	return q
}

func (q *queue) WithMessageTTL(ttl uint32) *queue {
	q.table["x-message-ttl"] = time.Duration(ttl) * time.Second
	return q
}

func (q *queue) WithRetryDeadLetter(c *Client) *queue {
	dlx := "dlx.retry.exchange"
	dlqRoutingKey := fmt.Sprintf("%s-%s", q.exchangeName, q.name)
	c.
		DefineExchange(dlx, "direct", false, false).
		DefineQueue(fmt.Sprintf("dlq.retry.%s", q.name), dlqRoutingKey, false, false).
		WithDeadLetter(q.exchangeName, q.routingKey).
		WithMessageTTL(10) // delays for 10 seconds

	q.WithDeadLetter(dlx, dlqRoutingKey)

	return q
}

func (q *queue) Consume(name string, autoAck, exclusive, requeue bool, prefetchCount int, h handler) *queue {
	name = fmt.Sprintf("%s.%s.%s", appName, name, hostName)
	if _, ok := q.consumers[name]; ok {
		return q
	}
	c := &consumer{
		name:          name,
		autoAck:       autoAck,
		exclusive:     exclusive,
		requeue:       requeue,
		prefetchCount: prefetchCount,
		h:             h,
	}
	q.consumers[name] = c
	return q
}
