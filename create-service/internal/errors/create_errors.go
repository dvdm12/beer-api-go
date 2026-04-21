// Package errors defines domain-specific errors for the create service.
package errors

import "net/http"

const (
	// CodeDuplicateBeer indicates a duplicate resource.
	CodeDuplicateBeer = "DUPLICATE_BEER"

	// CodeInvalidInput indicates invalid input.
	CodeInvalidInput = "INVALID_INPUT"
)

// NewDuplicateBeerError returns an error for duplicate beer creation.
func NewDuplicateBeerError(message string) AppError {
	return &businessError{
		message:    message,
		code:       CodeDuplicateBeer,
		statusCode: http.StatusConflict,
	}
}

// NewValidationError returns an error for invalid input.
func NewValidationError(message string) AppError {
	return &businessError{
		message:    message,
		code:       CodeInvalidInput,
		statusCode: http.StatusUnprocessableEntity,
	}
}