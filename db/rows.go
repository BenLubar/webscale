package db

import (
	"database/sql"

	"github.com/BenLubar/webscale/internal"
)

// Rows is the result of a query. Its cursor starts before the first row
// of the result set. Use Next to advance through the rows:
//
//     rows, err := tx.Query(stmt, ...)
//     ...
//     defer rows.Close()
//     for rows.Next() {
//         var id int
//         var name string
//         err = rows.Scan(&id, &name)
//         ...
//     }
//     err = rows.Err() // get any error encountered during iteration
//     ...
type Rows struct {
	rows  *sql.Rows
	stmt  *Stmt
	tx    *Tx
	stack internal.StackTrace
}

func wrapRows(rows *sql.Rows, stmt *Stmt, tx *Tx) *Rows {
	if rows == nil {
		return nil
	}

	return &Rows{
		rows:  rows,
		stmt:  stmt,
		tx:    tx,
		stack: internal.Callers(3),
	}
}

// Close closes the Rows, preventing further enumeration. Unlike sql.Rows,
// Close does not return an error. Call Err to get the error.
func (rows *Rows) Close() {
	_ = rows.rows.Close()
}

// Err returns the error, if any, that was encountered during iteration.
// Err may be called after an explicit or implicit Close.
func (rows *Rows) Err() error {
	return wrapError(rows.rows.Err(), rows.stmt, rows.tx, rows.stack)
}

// Next prepares the next result row for reading with the Scan method. It
// returns true on success, or false if there is no next result row or an error
// happened while preparing it. Err should be consulted to distinguish between
// the two cases.
//
// Every call to Scan, even the first one, must be preceded by a call to Next.
func (rows *Rows) Next() bool {
	return rows.rows.Next()
}

// Scan copies the columns in the current row into the values pointed
// at by dest. The number of values in dest must be the same as the
// number of columns in Rows.
//
// Scan converts columns read from the database into the following
// common Go types and special types provided by the sql package:
//
//    *string
//    *[]byte
//    *int, *int8, *int16, *int32, *int64
//    *uint, *uint8, *uint16, *uint32, *uint64
//    *bool
//    *float32, *float64
//    *interface{}
//    *RawBytes
//    any type implementing Scanner (see Scanner docs)
//
// In the most simple case, if the type of the value from the source
// column is an integer, bool or string type T and dest is of type *T,
// Scan simply assigns the value through the pointer.
//
// Scan also converts between string and numeric types, as long as no
// information would be lost. While Scan stringifies all numbers
// scanned from numeric database columns into *string, scans into
// numeric types are checked for overflow. For example, a float64 with
// value 300 or a string with value "300" can scan into a uint16, but
// not into a uint8, though float64(255) or "255" can scan into a
// uint8. One exception is that scans of some float64 numbers to
// strings may lose information when stringifying. In general, scan
// floating point columns into *float64.
//
// If a dest argument has type *[]byte, Scan saves in that argument a
// copy of the corresponding data. The copy is owned by the caller and
// can be modified and held indefinitely. The copy can be avoided by
// using an argument of type *RawBytes instead; see the documentation
// for RawBytes for restrictions on its use.
//
// If an argument has type *interface{}, Scan copies the value
// provided by the underlying driver without conversion. When scanning
// from a source value of type []byte to *interface{}, a copy of the
// slice is made and the caller owns the result.
//
// Source values of type time.Time may be scanned into values of type
// *time.Time, *interface{}, *string, or *[]byte. When converting to
// the latter two, time.Format3339Nano is used.
//
// Source values of type bool may be scanned into types *bool,
// *interface{}, *string, *[]byte, or *RawBytes.
//
// For scanning into *bool, the source may be true, false, 1, 0, or
// string inputs parseable by strconv.ParseBool.
func (rows *Rows) Scan(dest ...interface{}) error {
	return wrapError(rows.rows.Scan(dest...), rows.stmt, rows.tx, rows.stack)
}

// Row is the result of calling QueryRow to select a single row.
type Row struct {
	row   *sql.Row
	stmt  *Stmt
	tx    *Tx
	stack internal.StackTrace
}

func wrapRow(row *sql.Row, stmt *Stmt, tx *Tx) *Row {
	return &Row{
		row:   row,
		stmt:  stmt,
		tx:    tx,
		stack: internal.Callers(3),
	}
}

// Scan copies the columns from the matched row into the values
// pointed at by dest. See the documentation on Rows.Scan for details.
// If more than one row matches the query,
// Scan uses the first row and discards the rest. If no row matches
// the query, Scan returns ErrNoRows.
func (row *Row) Scan(dest ...interface{}) error {
	return wrapError(row.row.Scan(dest...), row.stmt, row.tx, row.stack)
}
