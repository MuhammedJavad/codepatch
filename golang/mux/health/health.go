package health

import (
	"commissioner/pkg/mux"
	"commissioner/pkg/rabbitmqx"
	"commissioner/pkg/redisx"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

type (
	option   func(context.Context) error
)

func Check(s *mux.Server, path string, opts ...option) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		for _, opt := range opts {
			if err := opt(r.Context()); err != nil {
				slog.ErrorContext(r.Context(), "health check error!!", "error", err.Error())
				mux.BadRequest(w, err)
				return
			}
		}
		mux.OK(w, map[string]string{"message": "All services are healthy."})
	}
	s.GET(path, fn)
}

func CheckMySql(name string, db *sql.DB) option {
	return func(ctx context.Context) error {
		if err := db.PingContext(ctx); err != nil {
			return fmt.Errorf("MySql connection failed. [err: %w]", err)
		}

		return nil
	}
}

func CheckRedis(name string, db *redisx.DB) option {
	return func(ctx context.Context) error {
		if err := db.Ping(ctx); err != nil {
			return fmt.Errorf("Redis connection failed. [err: %w]", err)
		}

		return nil
	}
}

func CheckRabbitMQ(name string, conn *rabbitmqx.Connection) option {
	return func(ctx context.Context) error {
		if !conn.IsHealthy() {
			return errors.New("RabbitMQ connection is not healthy")
		}

		return nil
	}
}