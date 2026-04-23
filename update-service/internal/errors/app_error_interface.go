// Package errors provides domain and HTTP error definitions.
package errors

import "fmt"

// Standard HTTP error codes.
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

// AppError defines the standard error contract.
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

// New creates a base AppError.
func New(message, code string, statusCode int) AppError {
	return &businessError{message: message, code: code, statusCode: statusCode}
}

// Wrap embeds a root error into an AppError.
func Wrap(cause error, message, code string, statusCode int) AppError {
	return &businessError{
		message:    fmt.Sprintf("%s: %v", message, cause),
		code:       code,
		statusCode: statusCode,
		cause:      cause,
	}
}

// getCode overrides the default code if a custom one is provided.
func getCode(defaultCode string, customCode []string) string {
	if len(customCode) > 0 && customCode[0] != "" {
		return customCode[0]
	}
	return defaultCode
}

// --- Error Constructors ---

func NotFound(resource string, customCode ...string) AppError {
	return New(fmt.Sprintf("%s not found", resource), getCode(ErrNotFound, customCode), 404)
}

func BadRequest(message string, customCode ...string) AppError {
	return New(message, getCode(ErrBadRequest, customCode), 400)
}

func ValidationError(message string, customCode ...string) AppError {
	return New(message, getCode(ErrValidation, customCode), 422)
}

func Conflict(message string, customCode ...string) AppError {
	return New(message, getCode(ErrConflict, customCode), 409)
}

func Unauthorized(message string) AppError { return New(message, ErrUnauthorized, 401) }
func Forbidden(message string) AppError    { return New(message, ErrForbidden, 403) }
func Internal(cause error) AppError        { return Wrap(cause, "internal server error", ErrInternal, 500) }

// FromError safely casts a standard error to an AppError.
func FromError(err error) (AppError, bool) {
	if err == nil {
		return nil, false
	}
	if appErr, ok := err.(AppError); ok {
		return appErr, true
	}
	return Internal(err), false
}
