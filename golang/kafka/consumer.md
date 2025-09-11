# Kafka Consumer Library with Automatic Reconnection

This enhanced Kafka consumer library provides robust connection management with automatic reconnection capabilities, configurable retry intervals, and comprehensive error handling for consuming messages from Kafka topics. It uses Go's structured logger `log/slog`.

## Features

- **Automatic Reconnection**: Automatically detects connection failures and attempts to reconnect
- **Configurable Retry Intervals**: Customizable retry intervals to prevent overwhelming the broker
- **Connection State Tracking**: Real-time monitoring of connection status
- **Comprehensive Error Handling**: Handles different types of Kafka errors appropriately
- **Metrics Integration**: Prometheus metrics for monitoring connection state and reconnection events
- **Thread-Safe**: Safe for concurrent use with proper locking mechanisms

## Usage

### Basic Usage

```go
import (
    "context"
    "log/slog"
    "time"
    "github.com/MuhammedJavad/codepatch/kafka"
)

// Create a logger
logger := slog.Default()

// Create consumer with reconnection settings
consumer, err := kafka.NewConsumer(
    []string{"localhost:9092"},
    "username",
    "password",
    "my-app",
    logger,
    0,                    // maxRetries: 0 = infinite retries
    5*time.Second,        // interval: 5 second delay between retries
)
if err != nil {
    logger.Error("Failed to create Kafka consumer", "error", err)
    return
}
defer consumer.Close()

// Subscribe to topics
err = consumer.SubscribeTopics("my-topic", "my-handler", func(ctx context.Context, msg kafka.Message) error {
    logger.Info("Received message", "topic", msg.Topic, "offset", msg.Offset)
    // Process your message here
    return nil
})
if err != nil {
    logger.Error("Failed to subscribe to topic", "error", err)
    return
}

// Start listening (handles reconnections automatically)
ctx := context.Background()
go consumer.Listen(ctx)
```

### Advanced Usage with Custom Reconnection Settings

```go
import "time"

// Create consumer with custom reconnection settings
consumer, err := kafka.NewConsumer(
    []string{"localhost:9092"},
    "username",
    "password",
    "my-app",
    logger,
    10,                   // maxRetries: Maximum 10 reconnection attempts
    2*time.Second,        // interval: 2 second delay between retries
)
if err != nil {
    logger.Error("Failed to create Kafka consumer", "error", err)
    return
}
defer consumer.Close()
```

## Configuration Options

### NewConsumer Parameters

| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| `servers` | `[]string` | List of Kafka broker addresses | `[]string{"localhost:9092"}` |
| `username` | `string` | SASL username for authentication | `"myuser"` |
| `password` | `string` | SASL password for authentication | `"mypass"` |
| `appname` | `string` | Consumer group ID | `"my-app"` |
| `logger` | `*slog.Logger` | Logger instance | `slog.Default()` |
| `maxRetries` | `int` | Maximum reconnection attempts (0 = infinite) | `0` or `10` |
| `interval` | `time.Duration` | Delay between reconnection attempts | `5*time.Second` |

## Error Handling

The consumer automatically handles various types of Kafka errors:

- **Connection Errors**: Network failures, broker unavailability, authentication failures
- **Temporary Errors**: Rebalance in progress (no reconnection needed)
- **Timeout Errors**: Normal polling timeouts (no reconnection needed)
- **Unknown Errors**: Assumes connection issue and attempts reconnection

The consumer uses a comprehensive error classification system that determines whether an error requires reconnection or can be handled gracefully.

## Metrics

The consumer exposes the following Prometheus metrics:

- `kafka_connection_state`: Connection state (1 = connected, 0 = disconnected)
- `kafka_reconnect_attempts_total`: Total reconnection attempts by status (attempt, success, failed)
- `kafka_consumer_duration_seconds`: Message processing duration

## Connection State Management

The consumer maintains connection state and provides methods to check status:

```go
// Check if connected
if consumer.IsConnected() {
    logger.Info("Kafka consumer is connected")
} else {
    logger.Warn("Kafka consumer is not connected")
}
```

## Automatic Reconnection Process

1. **Detection**: Consumer detects connection failures through error handling
2. **Retry Logic**: Waits with configurable interval before attempting reconnection
3. **Reconnection**: Creates new Kafka consumer with same configuration
4. **Re-subscription**: Automatically re-subscribes to all previously subscribed topics
5. **State Update**: Updates connection state and metrics
6. **Logging**: Logs all reconnection attempts and results

## Thread Safety

The consumer is thread-safe and can be used concurrently. All internal operations are protected with appropriate locking mechanisms.

## Best Practices

1. **Always check connection state** before critical operations
2. **Use appropriate log levels** to monitor reconnection events
3. **Configure reasonable retry limits** to prevent infinite loops
4. **Monitor metrics** to track connection health
5. **Handle context cancellation** properly in your message handlers

## Example Integration

This consumer library is designed specifically for consuming messages from Kafka topics with automatic reconnection capabilities. It provides a robust foundation for building reliable message processing applications.

### Key Benefits

- **Reliability**: Automatic reconnection ensures continuous message processing
- **Observability**: Comprehensive metrics and logging for monitoring
- **Flexibility**: Configurable retry behavior for different use cases
- **Simplicity**: Easy-to-use API with minimal configuration required

### Use Cases

- **Event Processing**: Consume and process business events from Kafka
- **Data Pipeline**: Build reliable data processing pipelines
- **Microservices**: Implement event-driven microservice architectures
- **Real-time Analytics**: Process streaming data for real-time insights
