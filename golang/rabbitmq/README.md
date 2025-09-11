# RabbitMQ Client (slog + Prometheus)

This RabbitMQ client provides:
- Automatic reconnect with configurable delay
- Structured logging via Go's `log/slog`
- Prometheus metrics: connection state, reconnect attempts, handler latency
- Simple builder API to declare exchanges/queues and start consumers

## Install

Module path: `github.com/MuhammedJavad/codepatch/rabbitmq`

```bash
cd golang/rabbitmq
go mod tidy
```

## Usage

```go
import (
    "context"
    "log/slog"
    rmq "github.com/MuhammedJavad/codepatch/rabbitmq"
)

func main() {
    logger := slog.Default()

    client := rmq.New(
        "localhost", "5672", "/", "guest", "guest",
        "my-app", "host-1",
        8 /* channel publish pool size */, logger,
    )

    client.
        DefineExchange("ex.events", "direct", true, false).
        DefineQueue("q.events", "rk.events", true, false).
        Consume("handlerA", false, false, false, 50, func(ctx context.Context, m rmq.Message) error {
            logger.Info("got", "len", len(m.Body))
            return nil
        })

    if err := client.Do(context.Background()); err != nil {
        logger.Error("connect", "err", err)
        return
    }
    // block ...
}
```

## Metrics
- `rabbitmq_connection_state{appname}`: 1 when connected, 0 otherwise
- `rabbitmq_reconnect_attempts_total{appname,status}`: attempt|success|failed
- `rabbitmq_consumer_duration_seconds{method,error}`: handler latency

## Notes
- Each consumer uses its own dedicated channel.
    - Closing a publisher channel does not affect consumer channels.
    - High-throughput consumers are isolated; their channels are not shared.
- Logging APIs use `InfoContext/ErrorContext/WarnContext` with structured attrs.
- A channel pool is used for publishing; broker returns and channel closes are logged.
- All external dependencies are managed via `go.mod`; no dasht dependencies.
