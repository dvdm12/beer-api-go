// Package errors defines domain-specific errors for the analysis service.
package errors

import (
	stderrors "errors"
	"net/http"

	"dataanalysis/internal/repository"
)

// Domain-specific error codes.
const (
	CodeBeerNotFound     = "BEER_NOT_FOUND"
	CodeEmptyCollection  = "EMPTY_COLLECTION"
	CodeInvalidID        = "INVALID_ID"
	CodeInvalidFilter    = "INVALID_FILTER"
	CodeInvalidSortField = "INVALID_SORT_FIELD"
	CodeInvalidAlcohol   = "INVALID_ALCOHOL"
	CodeInvalidYear      = "INVALID_YEAR"
	CodeEmptyBrand       = "EMPTY_BRAND"
	CodeEmptyName        = "EMPTY_NAME"
)

// InvalidID returns a 400 error for a malformed beer ID.
func InvalidID(id string) AppError {
	return &businessError{
		message:    "'" + id + "' is not a valid beer ID",
		code:       CodeInvalidID,
		statusCode: http.StatusBadRequest,
	}
}

// InvalidFilter returns a 422 error for invalid filter input.
func InvalidFilter(reason string) AppError {
	return &businessError{
		message:    "invalid filter: " + reason,
		code:       CodeInvalidFilter,
		statusCode: http.StatusUnprocessableEntity,
	}
}

// InvalidSortField returns a 400 error for an unknown sort field.
func InvalidSortField(field string) AppError {
	return &businessError{
		message:    "'" + field + "' is not a valid sort field",
		code:       CodeInvalidSortField,
		statusCode: http.StatusBadRequest,
	}
}

// InvalidAlcoholRange returns a 422 error for invalid alcohol bounds.
func InvalidAlcoholRange(min, max float64) AppError {
	return &businessError{
		message:    "invalid alcohol range",
		code:       CodeInvalidAlcohol,
		statusCode: http.StatusUnprocessableEntity,
	}
}

// InvalidYearRange returns a 422 error for invalid year bounds.
func InvalidYearRange(from, to int) AppError {
	return &businessError{
		message:    "invalid year range",
		code:       CodeInvalidYear,
		statusCode: http.StatusUnprocessableEntity,
	}
}

// EmptyBrand returns a 400 error when brand is empty.
func EmptyBrand() AppError {
	return &businessError{
		message:    "brand cannot be empty",
		code:       CodeEmptyBrand,
		statusCode: http.StatusBadRequest,
	}
}

// EmptyName returns a 400 error when name is empty.
func EmptyName() AppError {
	return &businessError{
		message:    "name cannot be empty",
		code:       CodeEmptyName,
		statusCode: http.StatusBadRequest,
	}
}

// BeerNotFound returns a 404 error for a missing beer.
func BeerNotFound(id string) AppError {
	return &businessError{
		message:    "beer with id '" + id + "' not found",
		code:       CodeBeerNotFound,
		statusCode: http.StatusNotFound,
	}
}

// EmptyCollection returns a 404 error when no beers exist.
func EmptyCollection() AppError {
	return &businessError{
		message:    "no beers found in database",
		code:       CodeEmptyCollection,
		statusCode: http.StatusNotFound,
	}
}

// NoResultsForFilter returns a 404 error when a query yields no results.
func NoResultsForFilter() AppError {
	return &businessError{
		message:    "no beers match the given filter",
		code:       CodeEmptyCollection,
		statusCode: http.StatusNotFound,
	}
}

// FromRepositoryError maps repository errors into AppError.
// Domain errors are translated; infrastructure errors become internal.
func FromRepositoryError(err error) AppError {
	if err == nil {
		return nil
	}

	// Already an AppError (e.g. validation layer)
	if appErr, ok := err.(AppError); ok {
		return appErr
	}

	// Repository domain errors
	switch err {
	case repository.ErrNoBeerFound:
		return BeerNotFound("unknown")
	case repository.ErrEmptyCollection:
		return EmptyCollection()
	case repository.ErrInvalidID:
		return InvalidID("unknown")
	}

	// Repository infrastructure errors → internal
	var repoErr *repository.RepoError
	if stderrors.As(err, &repoErr) {
		return Internal(err)
	}

	// Fallback
	return Internal(err)
}
