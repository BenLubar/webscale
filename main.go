// Command webscale implements the #webscale forum server.
package main // import "github.com/BenLubar/webscale"

import (
	"flag"
	"log"
	"net/http"

	"github.com/BenLubar/webscale/db"
)

func main() {
	flagDB := flag.String("db", "host=webscale-postgres user=postgres sslmode=disable", "PostgreSQL connection string")
	flagAddr := flag.String("addr", ":4567", "address to listen on for HTTP connections")

	flag.Parse()

	if err := db.Init(*flagDB); err != nil {
		log.Fatalf("fatal database error: %+v", err)
	}

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	log.Fatalf("fatal listen error: %+v", http.ListenAndServe(*flagAddr, nil))
}
