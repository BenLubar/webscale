package model_test

import (
	"testing"

	"github.com/BenLubar/webscale/db"
	"github.com/BenLubar/webscale/internal/testutils"
	"github.com/BenLubar/webscale/model"
)

func Context(t testing.TB) *model.Context {
	testutils.InitDB(t)

	tx, err := db.Begin()
	if err != nil {
		t.Skipf("database error: %+v", err)
	}

	return &model.Context{
		Tx: tx,
	}
}
