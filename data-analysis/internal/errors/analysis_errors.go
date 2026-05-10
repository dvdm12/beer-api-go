// Package errors defines domain-specific errors for the analysis service.
package errors

import (
	stderrors "errors"
	"net/http"

	"dataanalysis/internal/repository"
)

// Domain-specific error codes.
const (
	CodeEmptyCollection = "EMPTY_COLLECTION"
)

// EmptyCollection returns a 404 error when no beers exist.
func EmptyCollection() AppError {
	return &businessError{
		message:    "no beers found in database",
		code:       CodeEmptyCollection,
		statusCode: http.StatusNotFound,
	}
}

// FromRepositoryError maps repository errors into AppError.
// CategoryNotFound → EmptyCollection.
// Infrastructure errors → Internal.
func FromRepositoryError(err error) AppError {
	if err == nil {
		return nil
	}

	if appErr, ok := err.(AppError); ok {
		return appErr
	}

	var repoErr *repository.RepoError
	if stderrors.As(err, &repoErr) {
		if repoErr.Category == repository.CategoryNotFound {
			return EmptyCollection()
		}
		return Internal(err)
	}

	return Internal(err)
}
