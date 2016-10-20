package db_test

import (
	"testing"

	"github.com/BenLubar/webscale/db"
)

func TestInitTwice(t *testing.T) {
	db.InitForTesting()

	defer func() {
		r := recover()
		if r != "db: Init must only be called once" {
			t.Errorf("panic was not expected value: %#v", r)
		}
	}()

	db.Init("")
}
