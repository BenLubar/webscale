package main

import (
	"flag"

	"github.com/BenLubar/webscale/db"
)

func main() {
	flagDataSource := flag.String("db", "", "data source name")

	flag.Parse()

	db.Init(*flagDataSource)

	// PLACEHOLDER
}
