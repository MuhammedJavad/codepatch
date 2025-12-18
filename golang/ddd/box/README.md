# Box Package

The `box` package provides a generic message box implementation that can be used for both inbox and outbox patterns in Domain-Driven Design (DDD) architectures.

## Overview

The Box pattern is a common DDD pattern used for:
- **Inbox**: Receiving and processing incoming messages/events
- **Outbox**: Storing and publishing outgoing messages/events

This package provides a unified interface (`Box`) that can be used for both patterns, making it easy to implement reliable message processing.

## Components

### Box Interface

The `Box` interface provides two main operations:

```go
type Box interface {
    Save(ctx context.Context, message Message) error
    Iterator(ctx context.Context, source, event string, limit int) (Iterator, error)
}
```

- `Save`: Stores a new message in the box
- `Iterator`: Retrieves messages matching the given source and event filters, returning an iterator for processing

### Message

The `Message` struct represents a message with the following fields:

- `ID`: Unique identifier
- `Data`: Serialized message payload (JSON bytes)
- `ReceivedAt`: Timestamp when the message was received
- `ProcessedAt`: Timestamp when the message was processed (nil if not processed)
- `Source`: Source system/service identifier
- `Event`: Event type identifier
- `Error`: Error message if processing failed
- `RetryCount`: Number of retry attempts
- `RetryMax`: Maximum retry attempts (-1 for unlimited)

### Iterator

The `Iterator` interface provides a way to iterate through messages:

```go
type Iterator interface {
    Next() bool
    Current() Message
    Confirm(ctx context.Context, err error)
}
```

- `Next`: Advances to the next message, returns false if no more messages
- `Current`: Returns the current message
- `Confirm`: Confirms processing of the current message (with optional error)

## Usage Example

```go
package box

import (
	"context"
	"fmt"
)

func outboxAccessor() Accessor {}

func inboxAccessor() Accessor {}

func main() {
	// Save outbox message
	outbox := New(outboxAccessor())
	payload := map[string]string{"action": "create", "user": "john"}
	outboxMsg, err := NewMessage(payload, "user-service", "user.created", -1)
	if err != nil {
		panic(err)
	}
	if err := outbox.Save(context.Background(), outboxMsg); err != nil {
		panic(err)
	}
	
	// Save inbox message
	inbox := New(inboxAccessor())
	payload = map[string]string{"action": "update", "user": "jane"}
	inboxMsg, err := NewMessage(payload, "user-service", "user.updated", -1)
	if err != nil {
		panic(err)
	}
	if err := inbox.Save(context.Background(), inboxMsg); err != nil {
		panic(err)
	}
}

func inboxIterator() {
	ctx := context.Background()
	inbox := New(inboxAccessor())
	iterator, err := inbox.Iterator(ctx, "user-service", "user.created", 10)
	if err != nil {
		panic(err)
	}
	for iterator.Next() {
		message := iterator.Current()
		fmt.Printf("Processing inbox message ID: %d, Source: %s, Event: %s\n", message.ID, message.Source, message.Event)
		// process message
		iterator.Confirm(ctx, nil)
	}
}

func outboxIterator() {
	ctx := context.Background()
	outbox := New(outboxAccessor())
	iterator, err := outbox.Iterator(ctx, "user-service", "user.updated", 10)
	if err != nil {
		panic(err)
	}
	for iterator.Next() {
		message := iterator.Current()
		fmt.Printf("Processing inbox message ID: %d, Source: %s, Event: %s\n", message.ID, message.Source, message.Event)
		// process message
		iterator.Confirm(ctx, nil)
	}
}
```

1. **Mock Accessors**: Implementation of mock accessors for both inbox and outbox
2. **Box Instances**: Creating Box instances for inbox and outbox
3. **Message Processing**: Using iterators to process inbox messages
4. **Message Publishing**: Using iterators to publish outbox messages

### Creating a Box Instance

```go
// Create a mock accessor
accessor := newMockInboxAccessor()

// Create a Box instance
inbox := box.New(accessor)
```

### Saving Messages

```go
message, _ := box.NewMessage(
    map[string]string{"action": "create", "user": "john"},
    "user-service",
    "user.created",
    3, // retry max
)

inbox.Save(ctx, message)
```

### Processing Messages with Iterator

```go
iterator, err := inbox.Iterator(ctx, "user-service", "", 10)
if err != nil {
    return err
}

for iterator.Next() {
    message := iterator.Current()
    // Process the message
    err := processMessage(message)
    // Confirm processing (marks as processed or records error)
    iterator.Confirm(ctx, err)
}
```

## Accessor Interface

To use the Box, you need to implement the `accessor` interface:

```go
type accessor interface {
    Insert(ctx context.Context, message Message) error
    Update(ctx context.Context, message Message) error
    Retrieve(ctx context.Context, source, event string, limit int) ([]Message, error)
}
```

The accessor is responsible for:
- **Insert**: Storing new messages
- **Update**: Updating existing messages (used when confirming)
- **Retrieve**: Retrieving unprocessed messages matching filters

## Example Functions

### ProcessInboxMessages

Processes messages from the inbox using an iterator. This function:
- Creates an iterator for unprocessed messages
- Iterates through each message
- Processes the message
- Confirms processing (success or error)

### PublishOutboxMessages

Publishes messages from the outbox using an iterator. This function:
- Creates an iterator for unprocessed messages
- Iterates through each message
- Publishes the message to an external system
- Confirms publishing (success or error)

## Running the Example

To run the example code:

```bash
cd ddd/box/example
go run main.go
```

The example demonstrates:
- Creating mock accessors for both inbox and outbox
- Creating Box instances for inbox and outbox
- Saving messages to both boxes
- Processing inbox messages using `ProcessInboxMessages` function with Iterator
- Publishing outbox messages using `PublishOutboxMessages` function with Iterator

You can also see the example code in `docs.go` which contains the same implementation as part of the box package.

## Benefits

1. **Unified Interface**: Same interface for both inbox and outbox patterns
2. **Reliable Processing**: Iterator pattern ensures messages are processed and confirmed
3. **Error Handling**: Built-in retry mechanism and error tracking
4. **Flexible Filtering**: Filter messages by source and event type
5. **Testable**: Easy to mock accessors for testing

## Integration

This package can be integrated with:
- Database accessors (PostgreSQL, MySQL, etc.)
- Message brokers (Kafka, RabbitMQ, etc.)
- Event sourcing systems
- CQRS architectures

