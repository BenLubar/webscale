package db

import (
	"context"
	"database/sql"

	"github.com/BenLubar/webscale/internal"
)

// Tx is a database transaction.
type Tx struct {
	tx    *sql.Tx
	ctx   context.Context
	stack internal.StackTrace
}

func wrapTx(tx *sql.Tx, ctx context.Context) *Tx {
	if tx == nil {
		return nil
	}

	return &Tx{
		tx:    tx,
		ctx:   ctx,
		stack: internal.Callers(3),
	}
}

// Begin starts a transaction. If a non-default isolation level is used that
// the driver doesn't support an error will be returned. Different drivers
// may have slightly different meanings for the same isolation level.
func Begin(ctx context.Context) (*Tx, error) {
	<-ready
	tx, err := conn.BeginContext(ctx)
	wrappedTx := wrapTx(tx, ctx)
	return wrappedTx, wrapError(err, nil, wrappedTx)
}

func (tx *Tx) appendErrorDetails(buf []byte) []byte {
	buf = append(buf, "\n\nTransaction started here:"...)
	buf = tx.stack.AppendTo(buf)

	return buf
}

// Cancel rolls back the transaction without returning an error. It is assumed
// that by the time Cancel is called, Commit has already been called or another
// (more important) error has occurred.
func (tx *Tx) Cancel() {
	_ = tx.tx.Rollback()
}

// Rollback aborts the transaction.
func (tx *Tx) Rollback() error {
	return wrapError(tx.tx.Rollback(), nil, tx)
}

// Commit commits the transaction.
func (tx *Tx) Commit() error {
	return wrapError(tx.tx.Commit(), nil, tx)
}

var closedChanStruct = func() <-chan struct{} {
	ch := make(chan struct{})
	close(ch)
	return ch
}()

// Prepare prepares a statement local to this transaction.
func (tx *Tx) Prepare(query string) (*Stmt, error) {
	stmt, err := tx.tx.PrepareContext(tx.ctx, query)
	if stmt == nil {
		return nil, wrapError(err, nil, tx)
	}

	return &Stmt{
		stmt:  stmt,
		query: query,
		ready: closedChanStruct,
		stack: internal.Callers(2),
	}, wrapError(err, nil, tx)
}

// Exec executes a prepared statement with the given arguments and returns
// the number of rows affected for INSERT and UPDATE statements.
func (tx *Tx) Exec(stmt *Stmt, args ...interface{}) (int64, error) {
	<-stmt.ready
	res, err := tx.tx.StmtContext(tx.ctx, stmt.stmt).ExecContext(tx.ctx, args...)
	var affected int64
	if res != nil {
		affected, _ = res.RowsAffected()
	}
	return affected, wrapError(err, stmt, tx)
}

// Query executes a prepared query statement with the given arguments and
// returns the query results as a *Rows.
func (tx *Tx) Query(stmt *Stmt, args ...interface{}) (*Rows, error) {
	<-stmt.ready
	rows, err := tx.tx.StmtContext(tx.ctx, stmt.stmt).QueryContext(tx.ctx, args...)
	return wrapRows(rows, stmt, tx), wrapError(err, stmt, tx)
}

// QueryRow executes a prepared query statement with the given arguments.
// If an error occurs during the execution of the statement, that error will
// be returned by a call to Scan on the returned *Row, which is always non-nil.
// If the query selects no rows, the *Row's Scan will return ErrNoRows.
// Otherwise, the *Row's Scan scans the first selected row and discards
// the rest.
//
// Example usage:
//
//  var name string
//  err := tx.QueryRow(nameByUseridStmt, id).Scan(&name)
func (tx *Tx) QueryRow(stmt *Stmt, args ...interface{}) *Row {
	<-stmt.ready
	row := tx.tx.StmtContext(tx.ctx, stmt.stmt).QueryRowContext(tx.ctx, args...)
	return wrapRow(row, stmt, tx)
}
