package db

import (
	"database/sql"
	"sync"

	"github.com/BenLubar/webscale/internal"
)

var stmtWait sync.WaitGroup

// Stmt is a prepared SQL statement.
type Stmt struct {
	stmt  *sql.Stmt
	query string
	ready <-chan struct{}
	stack internal.StackTrace
}

// Prepare creates a prepared SQL statement. The database interaction happens
// asynchronously, and the program terminates if the query is malformed.
func Prepare(query string) *Stmt {
	done := make(chan struct{})
	stmt := &Stmt{
		query: query,
		ready: done,
		stack: internal.Callers(2),
	}
	stmtWait.Add(1)
	go stmt.prepare(done)

	return stmt
}

func (stmt *Stmt) prepare(done chan<- struct{}) {
	<-ready

	var err error
	stmt.stmt, err = conn.Prepare(stmt.query)
	if err == nil {
		close(done)
		stmtWait.Done()
		return
	}

	buf := []byte(err.Error())
	buf = stmt.appendErrorDetails(buf)
	buf = appendDriverErrorDetails(buf, err, stmt, nil)

	prepareFatal(string(buf), done, &stmtWait)
}

// variable for testing
var prepareFatal = func(message string, done chan<- struct{}, wg *sync.WaitGroup) {
	panic(message)
}

func (stmt *Stmt) appendErrorDetails(buf []byte) []byte {
	buf = append(buf, "\n\nQuery:\n\n"...)
	buf = append(buf, stmt.query...)
	buf = append(buf, "\n\nPrepared at:"...)
	buf = stmt.stack.AppendTo(buf)

	return buf
}

// WaitAll waits for all prepared statements to be ready.
func WaitAll() {
	stmtWait.Wait()
}

// Wait waits for the prepared statement to be ready.
func (stmt *Stmt) Wait() {
	<-stmt.ready
}
