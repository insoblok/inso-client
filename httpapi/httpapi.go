package httpapi

import (
	"encoding/json"
	"net/http"
)

// APIError represents a structured error in API responses
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// APIResponse is a generic wrapper for any API response payload
type APIResponse[T any] struct {
	Status int       `json:"status"`          // HTTP status code
	Data   *T        `json:"data,omitempty"`  // Successful payload
	Error  *APIError `json:"error,omitempty"` // Optional error
}

// WriteJSON writes an APIResponse to the http.ResponseWriter
func WriteJSON[T any](w http.ResponseWriter, status int, data *T, err *APIError) {
	resp := APIResponse[T]{
		Status: status,
		Data:   data,
		Error:  err,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(resp) // You may choose to log this if it fails
}

// Helper for sending success responses
func WriteOK[T any](w http.ResponseWriter, data *T) {
	WriteJSON(w, http.StatusOK, data, nil)
}

// Helper for sending errors
func WriteError(w http.ResponseWriter, status int, code, message string) {
	err := &APIError{Code: code, Message: message}
	WriteJSON[any](w, status, nil, err)
}
