package httpapi

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func ParseAPIResponse[T any](resp *http.Response) (*T, *APIError, error) {
	defer resp.Body.Close()

	var parsed APIResponse[T]
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	if parsed.Error != nil {
		return nil, parsed.Error, nil
	}
	return parsed.Data, nil, nil
}

func PostWithAPIResponse[T any](url string, payload any) (*T, *APIError, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal request payload: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, nil, fmt.Errorf("HTTP POST failed: %w", err)
	}
	defer resp.Body.Close()

	return ParseAPIResponse[T](resp)
}

func PostWithAPIResponseNoPayload[T any](url string) (*T, *APIError, error) {
	resp, err := http.Post(url, "application/json", bytes.NewReader([]byte{}))
	if err != nil {
		return nil, nil, fmt.Errorf("HTTP POST failed: %w", err)
	}
	defer resp.Body.Close()

	return ParseAPIResponse[T](resp)
}
