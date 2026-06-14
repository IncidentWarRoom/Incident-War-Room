package errs

import (
	"errors"
	"fmt"
)

type Kind string

const (
	KindInternal Kind = "INTERNAL"

	KindUnavailable Kind = "UNAVAILABLE"

	KindNotFound Kind = "NOT_FOUND"

	KindConflict Kind = "CONFLICT"

	KindValidation Kind = "VALIDATION"
)

type Error struct {
	Kind    Kind
	Op      string
	Message string
	Err     error
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

func New(kind Kind, op, message string) *Error {
	return &Error{Kind: kind, Op: op, Message: message}
}

func Wrap(kind Kind, op string, err error) *Error {
	if err == nil {
		return nil
	}
	return &Error{Kind: kind, Op: op, Err: err}
}

func Wrapf(kind Kind, op string, err error, format string, args ...any) *Error {
	if err == nil {
		return nil
	}
	return &Error{Kind: kind, Op: op, Message: fmt.Sprintf(format, args...), Err: err}
}

func KindOf(err error) Kind {
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr.Kind
	}
	return KindInternal
}

func Is(err error, kind Kind) bool {
	return KindOf(err) == kind
}
