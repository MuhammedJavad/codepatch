package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Message struct {
	ContentType     string
	ContentEncoding string
	Priority        uint8
	CorrelationID   string
	ReplyTo         string
	Expiration      string
	Timestamp       time.Time
	Body            []byte
}

/*
UnsafePublish sends a message to the server, but wouldn't wait for Publish.Confirm.

Since publishing is asynchronous in RabbitMQ, the server may respond with:
- Basic.Return for undeliverable messages
- Basic.Ack/Basic.Nack for publish confirmations

Basic.Return indicates an undeliverable message when:
- The mandatory flag is true and no queue matches the routing key, or
- The immediate flag is true and no consumer is ready (deprecated since RabbitMQ 2.9)

In this package:
- Mandatory flag is always true (Basic.Return events are logged)
- Immediate flag is not supported (deprecated)

Publish confirmations:
- Basic.Ack indicates successful processing when:
  - Message was consumed by all target consumers, or
  - Message was enqueued and persisted (if requested)

- Basic.Nack indicates failure when:
  - Message cannot be routed (with mandatory=false), or
  - Internal broker error occurs

Important behaviors:
1. Non-existent exchange: Channel will be closed (returns error)
2. Durable message to non-durable exchange:
  - Basic.Ack is sent (message accepted but not persistent)
  - No Basic.Return (if exchange exists at publish time)

3. Performance considerations:
  - Blocking for confirmations is resource-intensive, so its not considered here
  - Only Basic.Return is logged by default
  - Implement async handlers for delivery confirmation if needed

This method returns an error if the channel, connection, or socket is closed.

NOTES:
- Context is not supported in github.com/rabbitmq/amqp091-go v1.10.0
- Extend this package for custom confirmation handling
*/
func (c *Client) UnsafePublish(ctx context.Context, exchange, key string, msg Message) error {
	var err error

	ch, err0 := c.getChannel(ctx)
	if err0 != nil {
		err = err0
		return err
	}
	defer func() {
		c.channels <- ch
	}()

	dao := amqp.Publishing{
		ContentType:     msg.ContentType,
		ContentEncoding: msg.ContentEncoding,
		DeliveryMode:    1,
		Priority:        msg.Priority,
		CorrelationId:   msg.CorrelationID,
		ReplyTo:         msg.ReplyTo,
		Expiration:      msg.Expiration,
		Timestamp:       msg.Timestamp,
		Body:            msg.Body,
		MessageId:       uuid.NewString(),
	}
	err = ch.PublishWithContext(ctx, exchange, key, true, false, dao)
	if err != nil {
		return fmt.Errorf("failed to publish: %w", err)
	}

	return nil
}

// getChannel retrieves an active AMQP channel for publishing.
// It first attempts to reuse a channel from the pool (`p.collection`),
// and returns it if it's not nil and still open.
//
// If no valid channel is available, it creates a new one from the
// current AMQP connection. The method also sets up a goroutine to
// listen for returned messages (undeliverable messages) via
// ch.NotifyReturn, logging each returned message.
//
// Returns the active channel or an error if channel creation fails.
//
// Notes: make sure to Put the channel back in the p.collection
func (c *Client) getChannel(ctx context.Context) (*amqp.Channel, error) {
	if c.channels == nil {
		return nil, errors.New("configure publisher")
	}

	c.channelLock.Lock()
	defer c.channelLock.Unlock()

	connectionFactory := func() (*amqp.Channel, error) {
		connection := c.connection.Load()
		if connection == nil {
			return nil, fmt.Errorf("no register connection")
		}

		ch, err := connection.Channel()
		if err != nil {
			return nil, fmt.Errorf("failed to open channel for publisher: %w", err)
		}

		returnChan := ch.NotifyReturn(make(chan amqp.Return, 1))
		closeChan := ch.NotifyClose(make(chan *amqp.Error, 1))
		go func() {
			for {
				select {
				case ret := <-returnChan:
					c.logger.InfoContext(context.Background(), "AMQP return received",
						"exchange", ret.Exchange,
						"routing_key", ret.RoutingKey,
						"content_type", ret.ContentType,
						"message_id", ret.MessageId,
						"correlation_id", ret.CorrelationId,
						"timestamp", ret.Timestamp,
						"delivery_mode", ret.DeliveryMode,
						"priority", ret.Priority,
						"expiration", ret.Expiration,
						"user_id", ret.UserId,
						"app_id", ret.AppId,
					)
				case err := <-closeChan:
					c.logger.InfoContext(context.Background(), "AMQP channel closed", "error", err)
					return
				}
			}
		}()

		c.channelPoolSizeTracker++
		return ch, nil
	}

	if c.channelPoolSizeTracker < int32(cap(c.channels)) {
		return connectionFactory()
	}

	ch := <-c.channels
	if ch.IsClosed() {
		c.channelPoolSizeTracker--
		return connectionFactory()
	}
	return ch, nil
}
