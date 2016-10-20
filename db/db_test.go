package db

import (
	"flag"
	"sync"
)

var flagDataSource = flag.String("db", "", "data source name")
var testInitOnce sync.Once

func InitForTesting() {
	testInitOnce.Do(func() {
		Init(*flagDataSource)
	})
}

var TestingApplyFilterAppendErrorDetails = &applyFilterAppendErrorDetails
var TestingPrepareFatal = &prepareFatal
