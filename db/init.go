// Package db handles database access for #webscale.
package db // import "github.com/BenLubar/webscale/db"

import (
	"database/sql"

	"github.com/BenLubar/webscale/db/internal/schema"
	_ "github.com/lib/pq" // postgres database driver
	"github.com/pkg/errors"
)

var (
	// The database. Safe for use by multiple goroutines, but not available
	// until waitForInit has been closed by Init.
	theDB *sql.DB

	// Closed when the database is ready.
	waitForInit = make(chan struct{})
)

// Init initializes the database. If a non-nil error is returned, the program
// should exit and no further calls to functions in this package should be made.
//
// Initialization has 4 steps:
//
// 1. The *sql.DB is opened. This will always succeed as the only possible
//    error is the driver name being wrong.
// 2. The database is pinged. This connects to the database using the DSN
//    provided and ensures that the database connection information is valid.
// 3. The schema is checked. The schema_changes table is created and locked for
//    synchronization and recording of updates. This step can fail if a schema
//    change script fails or the database is from a newer version of #webscale.
// 4. Prepared queries are handled.
func Init(dataSourceName string) error {
	if theDB != nil {
		panic("Init must only be called once")
	}

	var err error
	if theDB, err = sql.Open("postgres", dataSourceName); err != nil {
		return errors.Wrap(err, "database init")
	}

	if err = theDB.Ping(); err != nil {
		return errors.Wrap(err, "database down during init")
	}

	if err = schema.Upgrade(theDB); err != nil {
		return errors.Wrap(err, "check database schema")
	}

	// allow prepared statements and transactions to be made
	close(waitForInit)

	// wait for prepared statements to be checked for errors
	prepareGroup.Wait()

	return nil
}
