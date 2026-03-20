package apierror

import (
	"errors"
	"fmt"
)

// APIError is a structured error returned by API handlers.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// New creates a generic APIError with the given code and message.
func New(code, message string) *APIError {
	return &APIError{Code: code, Message: message}
}

// NotFound returns a 404-style APIError.
func NotFound(message string) *APIError {
	return &APIError{Code: CodeNotFound, Message: message}
}

// Unauthorized returns a 401-style APIError.
func Unauthorized(message string) *APIError {
	return &APIError{Code: CodeUnauthorized, Message: message}
}

// Forbidden returns a 403-style APIError.
func Forbidden(message string) *APIError {
	return &APIError{Code: CodeForbidden, Message: message}
}

// BadRequest returns a 400-style APIError.
func BadRequest(code, message string) *APIError {
	return &APIError{Code: code, Message: message}
}

// Conflict returns a 409-style APIError.
func Conflict(message string) *APIError {
	return &APIError{Code: CodeConflict, Message: message}
}

// Internal returns a 500-style APIError.
func Internal(message string) *APIError {
	return &APIError{Code: CodeInternalError, Message: message}
}

// Is implements errors.Is support so callers can check with errors.As.
func (e *APIError) Is(target error) bool {
	var t *APIError
	if errors.As(target, &t) {
		return e.Code == t.Code
	}
	return false
}
