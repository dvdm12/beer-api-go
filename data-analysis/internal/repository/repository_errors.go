package repository

import (
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

// Category classifies repository errors.
type Category string

const (
	CategoryNotFound Category = "NOT_FOUND"
	CategoryTimeout  Category = "TIMEOUT"
	CategoryNetwork  Category = "NETWORK"
	CategoryDecode   Category = "DECODE"
	CategoryCursor   Category = "CURSOR"
	CategoryUnknown  Category = "UNKNOWN"
)

// Operation identifies the repository action that failed.
type Operation string

const (
	OpFind    Operation = "find"
	OpFindTop Operation = "findTop"
)

// Collection defines known MongoDB collections.
type Collection string

const (
	CollectionBeers Collection = "beers"
)

// Message constants used in RepoError.
const (
	MsgNotFound        = "document not found"
	MsgNilCursor       = "nil cursor received"
	MsgClientDisconn   = "client disconnected"
	MsgTimeout         = "operation timed out"
	MsgServerInterrupt = "server interrupted"
	MsgSocketError     = "socket or key error"
	MsgNetworkFailure  = "network failure"
	MsgDecodeFailed    = "document decode failed"
	MsgUnmapped        = "unmapped mongo error"
)

// RepoError represents a structured repository error.
type RepoError struct {
	Category  Category
	Message   string
	Cause     error
	Operation Operation
}

// Error implements the error interface.
func (e *RepoError) Error() string {
	if e.Operation != "" {
		return fmt.Sprintf("[%s] %s (%s): %v", e.Category, e.Message, e.Operation, e.Cause)
	}
	return fmt.Sprintf("[%s] %s: %v", e.Category, e.Message, e.Cause)
}

// Unwrap returns the underlying error.
func (e *RepoError) Unwrap() error {
	return e.Cause
}

func newRepoError(cat Category, msg string, op Operation, cause error) *RepoError {
	return &RepoError{Category: cat, Message: msg, Operation: op, Cause: cause}
}

// Logger defines minimal structured logging.
type Logger interface {
	Error(msg string, args ...any)
}

type noopLogger struct{}

func (n *noopLogger) Error(_ string, _ ...any) {}

// defaultLogger avoids nil checks when no logger is provided.
var defaultLogger Logger = &noopLogger{}

// String patterns used as fallback matching.
var (
	timeoutPatterns = []string{"context deadline exceeded", "timed out", "timeout"}
	networkPatterns = []string{"connection refused", "no reachable servers", "server selection error", "i/o timeout", "eof"}
	decodePatterns  = []string{"cannot decode", "no decoder found"}
)

// MapMongoError maps MongoDB errors into RepoError.
func MapMongoError(err error, op Operation, collection Collection, logger Logger) error {
	if err == nil {
		return nil
	}
	if logger == nil {
		logger = defaultLogger
	}

	logger.Error("mongo operation failed",
		"operation", op,
		"collection", collection,
		"error", err,
	)

	// Typed driver errors
	if errors.Is(err, mongo.ErrNoDocuments) {
		return newRepoError(CategoryNotFound, MsgNotFound, op, err)
	}
	if errors.Is(err, mongo.ErrNilCursor) {
		return newRepoError(CategoryCursor, MsgNilCursor, op, err)
	}
	if errors.Is(err, mongo.ErrClientDisconnected) {
		return newRepoError(CategoryNetwork, MsgClientDisconn, op, err)
	}

	// Command errors
	var cmdErr mongo.CommandError
	if errors.As(err, &cmdErr) {
		return mapCommandError(cmdErr, op, err)
	}

	// Fallback
	return mapByMessage(err, op)
}

// mapCommandError maps MongoDB command errors.
func mapCommandError(cmd mongo.CommandError, op Operation, original error) *RepoError {
	switch cmd.Code {
	case 50:
		return newRepoError(CategoryTimeout, MsgTimeout, op, original)
	case 11600, 11601:
		return newRepoError(CategoryNetwork, MsgServerInterrupt, op, original)
	case 9001, 211:
		return newRepoError(CategoryNetwork, MsgSocketError, op, original)
	default:
		return newRepoError(
			CategoryUnknown,
			fmt.Sprintf("command error code=%d: %s", cmd.Code, cmd.Message),
			op,
			original,
		)
	}
}

// mapByMessage matches errors by message patterns.
// Used only as a last-resort fallback.
func mapByMessage(err error, op Operation) *RepoError {
	msg := strings.ToLower(err.Error())

	for _, p := range timeoutPatterns {
		if strings.Contains(msg, p) {
			return newRepoError(CategoryTimeout, MsgTimeout, op, err)
		}
	}
	for _, p := range networkPatterns {
		if strings.Contains(msg, p) {
			return newRepoError(CategoryNetwork, MsgNetworkFailure, op, err)
		}
	}
	for _, p := range decodePatterns {
		if strings.Contains(msg, p) {
			return newRepoError(CategoryDecode, MsgDecodeFailed, op, err)
		}
	}

	return newRepoError(CategoryUnknown, MsgUnmapped, op, err)
}
