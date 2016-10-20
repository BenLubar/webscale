package main

import (
	"flag"

	"github.com/BenLubar/webscale/db"
)

func main() {
	flagDataSource := flag.String("db", "host=webscale-postgres user=postgres sslmode=disable", "data source name")

	flag.Parse()

	db.Init(*flagDataSource)

	// TODO
}
