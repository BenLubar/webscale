package db // import "github.com/BenLubar/webscale/db"

import (
	"database/sql"
	"log"
	"runtime/debug"
	"sync"

	"github.com/lib/pq"
)

var (
	// Used by Init to wait for prepared queries to finish before returning.
	prepareGroup sync.WaitGroup
)

// Prepare creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the
// returned statement.
func Prepare(query string) *Stmt {
	prepareGroup.Add(1)

	stmt := &Stmt{
		query: query,
		from:  string(debug.Stack()),
	}

	go stmt.once.Do(stmt.init)

	return stmt
}

// Stmt is a prepared statement.
type Stmt struct {
	query string    // query source SQL
	from  string    // discarded after preparing succeeds
	once  sync.Once // used to synchronize init
	impl  *sql.Stmt
}

func (stmt *Stmt) init() {
	<-waitForInit

	if stmtImpl, err := theDB.Prepare(stmt.query); err != nil {
		code := "?"
		if pe, ok := err.(*pq.Error); ok {
			code = string(pe.Code)
		}
		log.Fatalf("preparing statement failed: %s: %v\n\nsource text:\n%s\n\nstack trace:\n%s", code, err, stmt.query, stmt.from)
	} else {
		stmt.from = ""
		stmt.impl = stmtImpl
		prepareGroup.Done()
	}
}
