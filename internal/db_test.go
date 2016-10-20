package internal_test

import "flag"

// internal does not use package db, so we define a dummy flag here to make
// testing easier.
var _ = flag.String("db", "", "data source name (ignored)")
