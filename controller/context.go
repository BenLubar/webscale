package controller // import "github.com/BenLubar/webscale/controller"

import (
	"net/http"

	"github.com/BenLubar/webscale/db"
	"github.com/BenLubar/webscale/model"
)

// If a non-nil error is returned, the ResponseWriter must not have been used.
type handler func(w http.ResponseWriter, ctx *model.Context) error

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var ctx model.Context
	ctx.Request = r

	var err error
	ctx.Tx, err = db.Begin()
	if err != nil {
		handleError(w, &ctx, err, http.StatusServiceUnavailable)
		return
	}
	defer ctx.Tx.Rollback()

	err = h(w, &ctx)
	if err != nil {
		handleError(w, &ctx, err, http.StatusInternalServerError)
	}
}
