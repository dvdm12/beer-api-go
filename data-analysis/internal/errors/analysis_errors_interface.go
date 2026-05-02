// Package errors provides application-level error types with codes and HTTP status mapping.
package errors

import "fmt"

// Error codes used across the application.
const (
	ErrNotFound      = "ERR_NOT_FOUND"
	ErrUnauthorized  = "ERR_UNAUTHORIZED"
	ErrForbidden     = "ERR_FORBIDDEN"
	ErrValidation    = "ERR_VALIDATION"
	ErrConflict      = "ERR_CONFLICT"
	ErrInternal      = "ERR_INTERNAL"
	ErrBadRequest    = "ERR_BAD_REQUEST"
	ErrUnprocessable = "ERR_UNPROCESSABLE"
)

// AppError defines a structured application error.
type AppError interface {
	error
	Code() string    // Returns machine-readable error code.
	StatusCode() int // Returns HTTP status code.
	Unwrap() error   // Returns underlying cause.
}

// businessError is the default AppError implementation.
type businessError struct {
	message    string
	code       string
	statusCode int
	cause      error
}

// Error returns the error message.
func (e *businessError) Error() string { return e.message }

// Code returns the error code.
func (e *businessError) Code() string { return e.code }

// StatusCode returns the HTTP status code.
func (e *businessError) StatusCode() int { return e.statusCode }

// Unwrap returns the underlying error.
func (e *businessError) Unwrap() error { return e.cause }

// New creates an AppError without a cause.
func New(message, code string, statusCode int) AppError {
	return &businessError{
		message:    message,
		code:       code,
		statusCode: statusCode,
	}
}

// Wrap creates an AppError wrapping another error.
func Wrap(cause error, message, code string, statusCode int) AppError {
	return &businessError{
		message:    fmt.Sprintf("%s: %v", message, cause),
		code:       code,
		statusCode: statusCode,
		cause:      cause,
	}
}

// NotFound returns a 404 error for a missing resource.
func NotFound(resource string) AppError {
	return New(fmt.Sprintf("%s not found", resource), ErrNotFound, 404)
}

// Unauthorized returns a 401 error.
func Unauthorized(message string) AppError {
	return New(message, ErrUnauthorized, 401)
}

// Forbidden returns a 403 error.
func Forbidden(message string) AppError {
	return New(message, ErrForbidden, 403)
}

// BadRequest returns a 400 error.
func BadRequest(message string) AppError {
	return New(message, ErrBadRequest, 400)
}

// ValidationError returns a 422 validation error.
func ValidationError(message string) AppError {
	return New(message, ErrValidation, 422)
}

// Conflict returns a 409 conflict error.
func Conflict(message string) AppError {
	return New(message, ErrConflict, 409)
}

// Internal returns a 500 error wrapping the given cause.
func Internal(cause error) AppError {
	return Wrap(cause, "internal server error", ErrInternal, 500)
}

// FromError converts a generic error into AppError.
// Returns the original if already AppError, otherwise wraps it as internal.
func FromError(err error) (AppError, bool) {
	if err == nil {
		return nil, false
	}
	if appErr, ok := err.(AppError); ok {
		return appErr, true
	}
	return Internal(err), false
}
