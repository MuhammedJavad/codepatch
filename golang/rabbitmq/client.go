package rabbitmq

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	// appName is the name of the application.
	appName string

	// hostName is a unique identifier for each running instance (used for naming purposes).
	hostName string
)

const reconnectDelay = time.Second * 5

type (
	Client struct {
		connectionString       string
		connection             atomic.Pointer[amqp.Connection]
		channels               chan *amqp.Channel
		channelPoolSizeTracker int32
		channelLock            *sync.Mutex
		logger                 *slog.Logger
		latency                *prometheus.HistogramVec // Tracks how long it takes to process a message.
		connectionState        *prometheus.GaugeVec     // Reports connection state: 1 = connected, 0 = disconnected
		reconnectCount         *prometheus.CounterVec   // Counts reconnection attempts
		exchanges              map[string]*exchange
	}

	exchange struct {
		name       string
		kind       string
		durable    bool
		autoDelete bool
		queues     map[string]*queue
	}

	queue struct {
		name         string
		exchangeName string
		routingKey   string
		durable      bool
		autoDelete   bool
		exclusive    bool
		consumers    map[string]*consumer
		table        map[string]interface{}
	}
)

func (c *Client) Close() error {
	conn := c.connection.Load()
	if conn == nil {
		return nil
	}

	if err := conn.Close(); err != nil {
		return fmt.Errorf("failed to close RabbitMQ connection: %w", err)
	}
	c.connectionState.WithLabelValues(appName).Set(0)
	return nil
}

func (c *Client) Do(ctx context.Context) error {
	conn := c.connection.Load()
	if conn != nil && !conn.IsClosed() {
		return nil
	}

	conn, err := amqp.Dial(c.connectionString)
	if err != nil {
		fmt.Println(c.connectionString)
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	c.connection.Store(conn)

	if err := c.buildObjects(ctx); err != nil {
		if er := conn.Close(); er != nil {
			return fmt.Errorf("failed to close the channel: %w.\nfailed to build objects: %w", er, err)
		}
		return fmt.Errorf("failed to build objects: %w", err)
	}

	c.logger.InfoContext(ctx, "connected to RabbitMQ and successfully built objects")
	c.connectionState.WithLabelValues(appName).Set(1)

	go c.monitorConnection(ctx)

	return nil
}

func (c *Client) buildObjects(ctx context.Context) error {
	conn := c.connection.Load()
	if conn == nil {
		return fmt.Errorf("no register connection")
	}

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel for consumer: error: %w", err)
	}
	defer func() { ch.Close() }()

	for _, e := range c.exchanges {
		if err := ch.ExchangeDeclare(e.name, e.kind, e.durable, e.autoDelete, false, false, nil); err != nil {
			return fmt.Errorf("failed to declare exchange %q: %w", e.name, err)
		}

		c.logger.InfoContext(ctx, "exchange declared successfully", "exchange", e.name)

		for _, q := range e.queues {
			_, err := ch.QueueDeclare(q.name, q.durable, q.autoDelete, q.exclusive, false, q.table)
			if err != nil {
				return fmt.Errorf("failed to declare queue %q: %w", q.name, err)
			}

			c.logger.InfoContext(ctx, "queue declared successfully", "queue", q.name)

			if err := ch.QueueBind(q.name, q.routingKey, e.name, false, nil); err != nil {
				return fmt.Errorf("failed to bind queue %q to exchange %q: %w", q.name, e.name, err)
			}

			c.logger.InfoContext(ctx, "bound queue to exchange successfully", "queue", q.name, "exchange", e.name)

			for _, consumer := range q.consumers {
				go c.startConsumer(ctx, q.name, *consumer)
			}
		}
	}
	return nil
}

func (c *Client) monitorConnection(ctx context.Context) {
	conn := c.connection.Load()
	if conn == nil {
		return
	}

	errCh := conn.NotifyClose(make(chan *amqp.Error, 1))
	blockCh := conn.NotifyBlocked(make(chan amqp.Blocking, 1))

	select {
	case <-ctx.Done():
		c.logger.InfoContext(ctx, "rabbitmq connection closed")
		if err := c.Close(); err != nil {
			c.logger.ErrorContext(ctx, "error while closing rabbitmq connection", "error", err)
		}
		return // app is shutting down
	case err := <-errCh:
		if err != nil {
			c.logger.ErrorContext(ctx, "rabbitmq connection closed", "error", err)
			c.reconnectLoop(ctx)
		}
	case b := <-blockCh:
		c.logger.ErrorContext(ctx, "rabbitmq connection blocked", "active", b.Active, "reason", b.Reason)
	}
}

func (c *Client) reconnectLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			if err := c.Close(); err != nil {
				c.logger.ErrorContext(ctx, "error while closing rabbitmq connection", "error", err)
			}
			return
		default:
			c.logger.InfoContext(context.Background(), "attempting to reconnect...")
			err := c.Do(ctx)
			if err != nil {
				c.logger.ErrorContext(context.Background(), "reconnect failed", "error", err)
				time.Sleep(reconnectDelay)
				continue
			}
			c.logger.InfoContext(context.Background(), "reconnected successfully")
			go c.monitorConnection(ctx) // start monitoring new connection
			return
		}
	}
}
