// Package db handles database connections, transactions, and prepared
// statements for #webscale.
package db

import (
	"database/sql"
	"sync/atomic"

	"github.com/BenLubar/webscale/internal"
)

var conn *sql.DB
var initStarted uint32
var ready = make(chan struct{})

// Init sets the data source name and allows prepared statements and
// transactions to run.
func Init(dataSource string) {
	if !atomic.CompareAndSwapUint32(&initStarted, 0, 1) {
		panic("db: Init must only be called once")
	}

	db, err := sql.Open(driverName, dataSource)
	// The only possible error is an invalid driver name.
	internal.ImpossibleError(err)

	conn = db

	close(ready)
}
