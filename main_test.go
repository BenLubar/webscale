package main

import (
	"flag"
	"testing"

	"github.com/BenLubar/webscale/db"
)

var flagDB = flag.String("db", "host=webscale-postgres user=postgres sslmode=disable", "PostgreSQL connection string")

func TestInit(t *testing.T) {
	if err := db.Init(*flagDB); err != nil {
		t.Error(err)
	}
}
