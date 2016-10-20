package db

import (
	"context"
	"database/sql"
	"database/sql/driver"

	"github.com/BenLubar/hook"
	"github.com/BenLubar/webscale/internal"
	"github.com/pkg/errors"
)

// FilterAppendErrorDetails registers a function to be called when an error is
// being described. Filter functions should either return the byte slice
// unchanged or append text starting with a newline '\n'.
var FilterAppendErrorDetails = hook.NewFilter(&applyFilterAppendErrorDetails).(func(func([]byte, context.Context) ([]byte, error), int))
var applyFilterAppendErrorDetails func([]byte, context.Context) ([]byte, error)

// Error is the type of most errors returned by package db.
type Error struct {
	cause  error
	stmt   *Stmt
	tx     *Tx
	stacks []internal.StackTrace
}

func wrapError(err error, stmt *Stmt, tx *Tx, stacks ...internal.StackTrace) error {
	if err == nil || err == sql.ErrTxDone || err == sql.ErrNoRows || err == driver.ErrBadConn {
		return err
	}

	return &Error{
		cause:  err,
		stmt:   stmt,
		tx:     tx,
		stacks: append(stacks, internal.Callers(3)),
	}
}

// Cause implements the cause interface from github.com/pkg/errors.
func (err *Error) Cause() error {
	return err.cause
}

// Error implements the error interface.
func (err *Error) Error() string {
	buf := []byte(err.cause.Error())

	for _, stack := range err.stacks {
		buf = append(buf, "\n\nStack trace:"...)
		buf = stack.AppendTo(buf)
	}

	if err.tx != nil {
		buf = err.tx.appendErrorDetails(buf)
	}
	if err.stmt != nil {
		buf = err.stmt.appendErrorDetails(buf)
	}
	if err.tx != nil {
		ctxBuf, ctxErr := applyFilterAppendErrorDetails(nil, err.tx.ctx)
		if ctxErr != nil {
			buf = append(buf, "\n\n!!! ERROR WHILE GETTING ERROR DETAILS !!!\n"...)
			buf = append(buf, ctxErr.Error()...)
		} else {
			buf = append(buf, ctxBuf...)
		}
	}
	buf = appendDriverErrorDetails(buf, err.cause, err.stmt, err.tx)

	return string(buf)
}

// Timeout returns true if the cause of this error was a timeout.
func (err *Error) Timeout() bool {
	if te, ok := err.cause.(interface {
		Timeout() bool
	}); ok {
		return te.Timeout()
	}

	return false
}

// Temporary returns true if the cause of this error was temporary.
func (err *Error) Temporary() bool {
	if te, ok := err.cause.(interface {
		Temporary() bool
	}); ok {
		return te.Temporary()
	}

	return false
}

// StackTrace implements the stackTrace interface from github.com/pkg/errors.
func (err *Error) StackTrace() errors.StackTrace {
	return err.stacks[len(err.stacks)-1].StackTrace()
}

// IsConstraint returns true if the given error is a violation of a constraint
// with the given name.
func IsConstraint(err error, name string) bool {
	if e, ok := err.(*Error); ok {
		if constraint, ok := constraintFromDriverError(e.cause); ok {
			return constraint == name
		}
	}
	return false
}
