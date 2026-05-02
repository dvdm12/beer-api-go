// Package errors defines application-level error types and helpers.
package errors

import "fmt"

const (
	// ErrNotFound indicates a missing resource.
	ErrNotFound = "ERR_NOT_FOUND"

	// ErrUnauthorized indicates missing or invalid authentication.
	ErrUnauthorized = "ERR_UNAUTHORIZED"

	// ErrForbidden indicates insufficient permissions.
	ErrForbidden = "ERR_FORBIDDEN"

	// ErrValidation indicates invalid input data.
	ErrValidation = "ERR_VALIDATION"

	// ErrConflict indicates a resource conflict.
	ErrConflict = "ERR_CONFLICT"

	// ErrInternal indicates an internal server error.
	ErrInternal = "ERR_INTERNAL"

	// ErrBadRequest indicates a malformed request.
	ErrBadRequest = "ERR_BAD_REQUEST"

	// ErrUnprocessable indicates semantically invalid input.
	ErrUnprocessable = "ERR_UNPROCESSABLE"
)

// AppError represents an application error with metadata.
type AppError interface {
	error
	Code() string
	StatusCode() int
	Unwrap() error
}

type businessError struct {
	message    string
	code       string
	statusCode int
	cause      error
}

func (e *businessError) Error() string   { return e.message }
func (e *businessError) Code() string    { return e.code }
func (e *businessError) StatusCode() int { return e.statusCode }
func (e *businessError) Unwrap() error   { return e.cause }

// New creates a new AppError.
func New(message, code string, statusCode int) AppError {
	return &businessError{
		message:    message,
		code:       code,
		statusCode: statusCode,
	}
}

// Wrap creates an AppError wrapping an underlying error.
func Wrap(cause error, message, code string, statusCode int) AppError {
	return &businessError{
		message:    fmt.Sprintf("%s: %v", message, cause),
		code:       code,
		statusCode: statusCode,
		cause:      cause,
	}
}

// NotFound returns a not found error.
func NotFound(resource string) AppError {
	return New(
		fmt.Sprintf("%s not found", resource),
		ErrNotFound,
		404,
	)
}

// Unauthorized returns an unauthorized error.
func Unauthorized(message string) AppError {
	return New(message, ErrUnauthorized, 401)
}

// Forbidden returns a forbidden error.
func Forbidden(message string) AppError {
	return New(message, ErrForbidden, 403)
}

// BadRequest returns a bad request error.
func BadRequest(message string) AppError {
	return New(message, ErrBadRequest, 400)
}

// ValidationError returns a validation error.
func ValidationError(message string) AppError {
	return New(message, ErrValidation, 422)
}

// Conflict returns a conflict error.
func Conflict(message string) AppError {
	return New(message, ErrConflict, 409)
}

// Internal returns an internal server error.
func Internal(cause error) AppError {
	return Wrap(cause, "internal server error", ErrInternal, 500)
}

// FromError attempts to cast an error to AppError.
func FromError(err error) (AppError, bool) {
	if err == nil {
		return nil, false
	}
	if appErr, ok := err.(AppError); ok {
		return appErr, true
	}
	return Internal(err), false
}
