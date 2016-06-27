package testutils // import "github.com/BenLubar/webscale/internal/testutils"

import (
	"flag"
	"sync"
	"testing"

	"github.com/BenLubar/webscale/db"
)

// FlagDB holds the connection string for the test database.
var FlagDB = flag.String("db", "host=webscale-postgres user=postgres sslmode=disable", "PostgreSQL connection string")

var once sync.Once
var err error

// InitDB initializes the database and skips the test if there is an error in
// database initialization. The benchmark timer is stopped while waiting for
// the database.
func InitDB(t testing.TB) {
	startStop, ok := t.(interface {
		StartTimer()
		StopTimer()
	})
	if ok {
		startStop.StopTimer()
	}
	once.Do(func() {
		err = db.Init(*FlagDB)
	})
	if err != nil {
		t.Skipf("database init failed: %+v", err)
	}
	if ok {
		startStop.StartTimer()
	}
}
