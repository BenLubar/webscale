package db // import "github.com/BenLubar/webscale/db"

import (
	"database/sql"

	"github.com/pkg/errors"
)

// Begin starts a transaction.
func Begin() (*Tx, error) {
	<-waitForInit

	tx, err := theDB.Begin()
	return &Tx{impl: tx}, errors.Wrap(err, "begin transaction")
}

// Tx is an in-progress database transaction.
//
// A transaction must end with a call to Commit or Rollback.
//
// Calling methods after Commit or Rollback have been called has no effect.
type Tx struct {
	impl *sql.Tx
}

// Rollback aborts the transaction.
func (tx *Tx) Rollback() error {
	return errors.Wrap(tx.impl.Rollback(), "transaction rollback")
}

// Commit commits the transaction.
func (tx *Tx) Commit() error {
	return errors.Wrap(tx.impl.Commit(), "transaction commit")
}

// Exec executes a query that doesn't return rows.
// For example: an INSERT and UPDATE.
func (tx *Tx) Exec(query *Stmt, args ...interface{}) (sql.Result, error) {
	query.once.Do(query.init)

	result, err := tx.impl.Stmt(query.impl).Exec(args...)
	return result, errors.Wrap(err, "transaction exec")
}

// Query executes a query that returns rows, typically a SELECT.
func (tx *Tx) Query(query *Stmt, args ...interface{}) (*sql.Rows, error) {
	query.once.Do(query.init)

	rows, err := tx.impl.Stmt(query.impl).Query(args...)
	return rows, errors.Wrap(err, "transaction query")
}

// QueryRow executes a query that is expected to return at most one row.
// QueryRow always returns a non-nil value. Errors are deferred until
// Row's Scan method is called.
func (tx *Tx) QueryRow(query *Stmt, args ...interface{}) *sql.Row {
	query.once.Do(query.init)

	return tx.impl.Stmt(query.impl).QueryRow(args...)
}
