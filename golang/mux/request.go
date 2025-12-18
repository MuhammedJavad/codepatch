package mux

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// GetQueryParam extracts a query parameter from an HTTP request and converts it to type T.
// Supported types: string, int, int64, float64, bool.
// Returns the zero value of T if the param is missing or cannot be converted.
func GetQueryParam[T any](r *http.Request, name string) (T, error) {
	value := r.URL.Query().Get(name)
	return cast[T](value)
}

func GetPathValue[T any](r *http.Request, name string) (T, error) {
	value := r.PathValue(name)
	return cast[T](value)
}

func GetRequestID(r *http.Request) string {
	return r.Header.Get("X-Request-ID")
}

func JsonBody[T interface{}](r *http.Request) (T, error) {
	var b T
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		return b, err
	}

	return b, nil
}

func GetRequestBodyString(req *http.Request) string {
	if req.Body == nil {
		return ""
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return ""
	}
	// Restore the body to its original state
	req.Body = io.NopCloser(bytes.NewBuffer(body))
	return string(body)
}

// cast converts a string to a specified generic type T.
// Supported types: string, int, int64, uint32, uint64, float64, bool.
func cast[T any](actual string) (T, error) {
	var zero T
	if actual == "" {
		return zero, errors.New("value is required")
	}

	switch any(zero).(type) {
	case string:
		return any(actual).(T), nil

	case int:
		v, err := strconv.Atoi(actual)
		if err != nil {
			return zero, fmt.Errorf("invalid int value: %w", err)
		}
		return any(v).(T), nil

	case int64:
		v, err := strconv.ParseInt(actual, 10, 64)
		if err != nil {
			return zero, fmt.Errorf("invalid int64 value: %w", err)
		}
		return any(v).(T), nil

	case uint32:
		v, err := strconv.ParseUint(actual, 10, 32)
		if err != nil {
			return zero, fmt.Errorf("invalid uint32 value: %w", err)
		}
		return any(uint32(v)).(T), nil

	case uint64:
		v, err := strconv.ParseUint(actual, 10, 64)
		if err != nil {
			return zero, fmt.Errorf("invalid uint64 value: %w", err)
		}
		return any(v).(T), nil

	case float64:
		v, err := strconv.ParseFloat(actual, 64)
		if err != nil {
			return zero, fmt.Errorf("invalid float64 value: %w", err)
		}
		return any(v).(T), nil

	case bool:
		v, err := strconv.ParseBool(actual)
		if err != nil {
			return zero, fmt.Errorf("invalid bool value: %w", err)
		}
		return any(v).(T), nil

	default:
		return zero, fmt.Errorf("unsupported type: %T", zero)
	}
}