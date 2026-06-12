// Package errs defines the application-level error type shared by all layers.
// Wrap low-level failures into an *Error so callers can branch on Kind
// instead of inspecting error messages.
package errs

import (
	"errors"
	"fmt"
)

// Kind classifies an error into a failure category.
type Kind string

const (
	// KindInternal is an unexpected system failure (default category).
	KindInternal Kind = "INTERNAL"
	// KindUnavailable means an external dependency (DB, Telegram API, report service) is unreachable.
	KindUnavailable Kind = "UNAVAILABLE"
	// KindNotFound means a requested entity does not exist.
	KindNotFound Kind = "NOT_FOUND"
	// KindConflict means the operation contradicts current state (e.g. closing a closed incident).
	KindConflict Kind = "CONFLICT"
	// KindValidation means the input is malformed or violates business rules.
	KindValidation Kind = "VALIDATION"
)

// Error is the common application error.
type Error struct {
	Kind    Kind   // failure category
	Op      string // operation where the error occurred, e.g. "config.Load"
	Message string // human-readable description
	Err     error  // wrapped cause, may be nil
}

func (e *Error) Error() string {
	switch {
	case e.Err != nil && e.Message != "":
		return fmt.Sprintf("%s: %s: %v", e.Op, e.Message, e.Err)
	case e.Err != nil:
		return fmt.Sprintf("%s: %v", e.Op, e.Err)
	default:
		return fmt.Sprintf("%s: %s", e.Op, e.Message)
	}
}

func (e *Error) Unwrap() error { return e.Err }

// New creates an *Error without an underlying cause.
func New(kind Kind, op, message string) *Error {
	return &Error{Kind: kind, Op: op, Message: message}
}

// Wrap attaches a category and operation to an underlying error.
// Returns nil if err is nil.
func Wrap(kind Kind, op string, err error) *Error {
	if err == nil {
		return nil
	}
	return &Error{Kind: kind, Op: op, Err: err}
}

// Wrapf is Wrap with an additional human-readable message.
func Wrapf(kind Kind, op string, err error, format string, args ...any) *Error {
	if err == nil {
		return nil
	}
	return &Error{Kind: kind, Op: op, Message: fmt.Sprintf(format, args...), Err: err}
}

// KindOf returns the Kind of the first *Error in the chain,
// or KindInternal for unclassified errors.
func KindOf(err error) Kind {
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr.Kind
	}
	return KindInternal
}

// Is reports whether the error chain contains an *Error of the given Kind.
func Is(err error, kind Kind) bool {
	return KindOf(err) == kind
}
