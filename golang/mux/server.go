package mux

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func NewServer(subsystem string) *Server {
	l := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem: subsystem,
		Name:      "http_endpoints_latency",
		Help:      "Total duration of api ops in seconds",
		Buckets:   prometheus.DefBuckets,
	}, []string{"status_code", "method", "path", "error"})

	prometheus.Unregister(l)
	prometheus.MustRegister(l)

	return &Server{
		mux:     http.NewServeMux(),
		latency: l,
	}
}

type Server struct {
	mux     *http.ServeMux
	latency *prometheus.HistogramVec
}

func (s *Server) ListenAndServe(ctx context.Context, addr string) error {
	server := http.Server{
		Addr:    addr,
		Handler: s.mux,
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		slog.Info("Stopping HTTP server gracefully...")

		if err := server.Shutdown(shutdownCtx); err != nil {
			slog.Error("Failed to shutdown HTTP server", "error", err)
		}
	}()

	fmt.Println("Listening on", addr)
	return server.ListenAndServe()
}

func (s *Server) POST(pattern string, h http.HandlerFunc) {
	s.handle(http.MethodPost, pattern, h)
}

func (s *Server) GET(pattern string, h http.HandlerFunc) {
	s.handle(http.MethodGet, pattern, h)
}

func (s *Server) PUT(pattern string, h http.HandlerFunc) {
	s.handle(http.MethodPut, pattern, h)
}

func (s *Server) DELETE(pattern string, h http.HandlerFunc) {
	s.handle(http.MethodDelete, pattern, h)
}

func (s *Server) PATCH(pattern string, h http.HandlerFunc) {
	s.handle(http.MethodPatch, pattern, h)
}

func (s *Server) OPTIONS(pattern string, h http.HandlerFunc) {
	s.handle(http.MethodOptions, pattern, h)
}

func (s *Server) HEAD(pattern string, h http.HandlerFunc) {
	s.handle(http.MethodHead, pattern, h)
}

func (s *Server) handle(method, pattern string, h http.HandlerFunc) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				if rvr == http.ErrAbortHandler {
					// we don't recover http.ErrAbortHandler so the response
					// to the client is aborted, this should not be logged
					return
				}

				slog.ErrorContext(r.Context(), "panic", rvr)

				// This block checks if the request was not an upgrade (like a WebSocket), and if so,
				// it returns a 500 Internal Server Error response to the client.
				// Note: 'err' is not defined here. Most likely, the developer intended to send an error message.
				// For a correct implementation, you'd want to pass a meaningful error.
				if r.Header.Get("Connection") != "Upgrade" {
					InternalServerError(w)
				}
			}
		}()

		h.ServeHTTP(w, r)
	}
	s.mux.Handle(p(method, pattern), http.HandlerFunc(fn))
}

func p(method, pattern string) string {
	return fmt.Sprintf("%s %s", method, strings.ReplaceAll(pattern, " ", ""))
}
