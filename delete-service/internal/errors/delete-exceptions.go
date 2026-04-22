// Package errors defines domain-specific errors for the delete service.
package errors

import "net/http"

const (
	// CodeBeerNotFound indicates the beer does not exist.
	CodeBeerNotFound = "BEER_NOT_FOUND"

	// CodeInvalidID indicates the provided ID is not a valid ObjectID.
	CodeInvalidID = "INVALID_ID"
)

// NewBeerNotFoundError returns an error when the beer does not exist.
// HTTP 404 Not Found.
func NewBeerNotFoundError(id string) AppError {
	return &businessError{
		message:    "beer with id '" + id + "' not found",
		code:       CodeBeerNotFound,
		statusCode: http.StatusNotFound, // 404
	}
}

// NewInvalidIDError returns an error when the ID format is invalid.
// HTTP 400 Bad Request.
func NewInvalidIDError(id string) AppError {
	return &businessError{
		message:    "'" + id + "' is not a valid beer ID",
		code:       CodeInvalidID,
		statusCode: http.StatusBadRequest, // 400
	}
}