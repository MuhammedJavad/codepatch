package rabbitmq

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type (
	handler func(ctx context.Context, data Message) error

	consumer struct {
		name          string
		autoAck       bool
		exclusive     bool
		requeue       bool
		prefetchCount int
		h             handler
	}
)

func (c *Client) startConsumer(ctx context.Context, queueName string, consumer consumer) {
	conn := c.connection.Load()
	if conn == nil {
		return
	}

	ch, err := conn.Channel()
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to open channel for consumer", "error", err)
		return
	}
	defer func() { ch.Close() }()

	if err := ch.Qos(consumer.prefetchCount, 0, true); err != nil {
		c.logger.ErrorContext(ctx, "failed to set QoS", "error", err)
		return
	}

	deliveries, err := ch.Consume(queueName, consumer.name, consumer.autoAck, consumer.exclusive, false, false, nil)
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to consume", "error", err)
		return
	}

	closeChan := ch.NotifyClose(make(chan *amqp.Error, 1))
	for {
		select {
		case msg := <-deliveries:
			go c.handleMessage(msg, consumer)
		case err := <-closeChan:
			c.logger.ErrorContext(ctx, "AMQP channel closed", "error", err)
			return
		case <-ctx.Done():
			c.logger.InfoContext(ctx, "consumer stopped")
			return
		}
	}
}

func (c *Client) handleMessage(msg amqp.Delivery, consumer consumer) {
	rid := uuid.New().String()
	ctx := context.WithValue(context.Background(), "request_id", rid)
	ctx = context.WithValue(ctx, "parent_method", consumer.name)

	var err error
	start := time.Now()
	defer func() {
		c.latency.WithLabelValues(consumer.name, fmt.Sprint(err != nil)).Observe(time.Since(start).Seconds())
	}()

	content := Message{
		Body:            msg.Body,
		ContentType:     msg.ContentType,
		ContentEncoding: msg.ContentEncoding,
		Timestamp:       msg.Timestamp,
	}
	if err = consumer.h(ctx, content); err != nil {
		c.logger.ErrorContext(ctx, "handler error", "error", err)

		if isRejectError(err) {
			c.logger.InfoContext(ctx, "message rejected")
			msg.Reject(false)
			return
		}

		c.logger.InfoContext(ctx, "message nacked")
		msg.Nack(false, consumer.requeue)
		return
	}

	if err = msg.Ack(false); err != nil {
		c.logger.ErrorContext(ctx, "ack error", "error", err)
	}
}
