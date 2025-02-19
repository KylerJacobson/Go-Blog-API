package httperr

import (
	"encoding/json"
	"errors" // Now we can use the standard httperr package clearly
	"net/http"
)

// Error represents a structured HTTP error response
type Error struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

// Error implements the error interface
func (e *Error) Error() string {
	return e.Message
}

// New creates a new Error
func New(status int, message string, detail string) *Error {
	return &Error{
		Status:  status,
		Message: message,
		Detail:  detail,
	}
}

// Common constructors
func BadRequest(message string, detail string) *Error {
	return New(http.StatusBadRequest, message, detail)
}

func NotFound(message string, detail string) *Error {
	return New(http.StatusNotFound, message, detail)
}

func Internal(message string, detail string) *Error {
	return New(http.StatusInternalServerError, message, detail)
}

// Write sends the error response to the http.ResponseWriter
func Write(w http.ResponseWriter, err error) {
	var httpErr *Error

	// Check if the error is already our type
	if !errors.As(err, &httpErr) {
		// If not, wrap it as an internal error
		httpErr = Internal(
			"An unexpected error occurred",
			err.Error(),
		)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpErr.Status)
	_ = json.NewEncoder(w).Encode(httpErr)
}
