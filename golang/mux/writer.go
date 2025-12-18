package mux

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func OK(w http.ResponseWriter, data interface{}) {
	writeJSON(w, http.StatusOK, data)
}

func Created(w http.ResponseWriter, data interface{}) {
	writeJSON(w, http.StatusCreated, data)
}

func Accepted(w http.ResponseWriter) {
	w.WriteHeader(http.StatusAccepted)
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func SeeOther(w http.ResponseWriter, location string) {
	w.Header().Set(LOCATION_HEADER, location)
	w.WriteHeader(http.StatusSeeOther)
}

func NotModified(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotModified)
}

func UseProxy(w http.ResponseWriter, location string) {
	w.Header().Set(LOCATION_HEADER, location)
	w.WriteHeader(http.StatusUseProxy)
}

func TemporaryRedirect(w http.ResponseWriter, location string) {
	w.Header().Set(LOCATION_HEADER, location)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func PermanentRedirect(w http.ResponseWriter, location string) {
	w.Header().Set(LOCATION_HEADER, location)
	w.WriteHeader(http.StatusPermanentRedirect)
}

func InternalServerError(w http.ResponseWriter) {
	writeError(w, http.StatusInternalServerError, fmt.Errorf("internal server error"))
}

func BadRequest(w http.ResponseWriter, err error) {
	writeError(w, http.StatusBadRequest, err)
}

func NotFound(w http.ResponseWriter, err error) {
	writeError(w, http.StatusNotFound, err)
}

func NotAcceptable(w http.ResponseWriter, err error) {
	writeError(w, http.StatusNotAcceptable, err)
}

func RequestTimeout(w http.ResponseWriter, err error) {
	writeError(w, http.StatusRequestTimeout, err)
}

func Conflict(w http.ResponseWriter, err error) {
	writeError(w, http.StatusConflict, err)
}

func UnprocessableEntity(w http.ResponseWriter, err error) {
	writeError(w, http.StatusUnprocessableEntity, err)
}

func Forbidden(w http.ResponseWriter, err error) {
	writeError(w, http.StatusForbidden, err)
}

// Helper methods
func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	// Delete the Content-Length header, which might be for some other content.
	// Assuming the error string fits in the writer's buffer, we'll figure
	// out the correct Content-Length for it later.
	//
	// We don't delete Content-Encoding, because some middleware sets
	// Content-Encoding: gzip and wraps the ResponseWriter to compress on-the-fly.
	// See https://go.dev/issue/66343.
	w.Header().Del("Content-Length")

	// There might be content type already set, but we reset it to
	// text/plain for the error message.
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)

	response := happyResponse(data)
	json.NewEncoder(w).Encode(response)
}

func writeError(w http.ResponseWriter, statusCode int, err error) {
	// Delete the Content-Length header, which might be for some other content.
	// Assuming the error string fits in the writer's buffer, we'll figure
	// out the correct Content-Length for it later.
	//
	// We don't delete Content-Encoding, because some middleware sets
	// Content-Encoding: gzip and wraps the ResponseWriter to compress on-the-fly.
	// See https://go.dev/issue/66343.
	w.Header().Del(CONTENT_LENGTH_HEADER)

	// There might be content type already set, but we reset it to
	// text/plain for the error message.
	w.Header().Set(CONTENT_TYPE_HEADER, "application/json")
	w.Header().Set(X_CONTENT_TYPE_OPTIONS_HEADER, "nosniff")
	w.WriteHeader(statusCode)

	response := unhappyResponse(err.Error())
	json.NewEncoder(w).Encode(response)
}

type response struct {
	Ok    bool        `json:"ok"`
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

func happyResponse(data interface{}) response {
	if data == nil {
		data = struct{}{}
	}
	return response{
		Ok:    true,
		Data:  data,
		Error: "",
	}
}

func unhappyResponse(error string) response {
	return response{
		Ok:    false,
		Data:  struct{}{},
		Error: error,
	}
}
