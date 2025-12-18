package mux

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (s *Server) StaticToken(token string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const (
			authPrefix = "Basic "
			header = "Authorization"
		)

		authValue := r.Header.Get(header)
		if !strings.HasPrefix(authValue, authPrefix) {
			Forbidden(w, errors.New("missing Authorization header"))
			return
		}

		reqToken := strings.TrimPrefix(authValue, authPrefix)
		if reqToken != token {
			Forbidden(w, errors.New("unauthorized"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Observability wraps an http.Handler to track metrics
func (s *Server) Observability(name string, next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const SCOPE_KEY = "scope"

		startTime := time.Now()

		// Get or create request ID
		rid := r.Header.Get(REQUEST_ID_HEADER)
		if rid == "" {
			rid = uuid.New().String()
			r.Header.Set(REQUEST_ID_HEADER, rid)
			w.Header().Set(REQUEST_ID_HEADER, rid)
		}

		// Add request ID to context
		ctx := context.WithValue(r.Context(), REQUEST_ID_HEADER, rid)
		ctx = context.WithValue(ctx, SCOPE_KEY, name)
		r = r.WithContext(ctx)

		// Wrap response writer to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the next handler
		next.ServeHTTP(wrapped, r)

		// Record metrics
		lvs := []string{
			strconv.Itoa(wrapped.statusCode),
			name,
			r.URL.Path,
			strconv.FormatBool(wrapped.statusCode >= 500 && wrapped.statusCode < 600),
		}
		s.latency.WithLabelValues(lvs...).Observe(time.Since(startTime).Seconds())
	})
}

// Logger wraps an http.Handler to log incoming requests and responses
func (s *Server) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Capture request info
		path := r.URL.Path
		params := r.URL.Query()
		bodyString := GetRequestBodyString(r)
		rid := r.Header.Get(REQUEST_ID_HEADER)

		slog.Info("request received",
			"request_id", rid,
			"path", path,
			"params", params,
			"body", bodyString,
			"method", r.Method,
		)

		// Wrap response writer to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the next handler
		next.ServeHTTP(wrapped, r)

		// Log response
		slog.Info("request completed",
			"request_id", rid,
			"path", path,
			"status_code", wrapped.statusCode,
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.written = true
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.written {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}
