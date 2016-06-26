package controller // import "github.com/BenLubar/webscale/controller"

import (
	"net/http"
	"strconv"

	"github.com/BenLubar/webscale/db"
	"github.com/BenLubar/webscale/model"
)

// If a non-nil error is returned, the ResponseWriter must not have been used.
type handler func(w http.ResponseWriter, ctx *model.Context) error

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var ctx model.Context
	ctx.Request = r

	if page := ctx.Request.FormValue("page"); page != "" {
		if page == "0" {
			u := *r.URL
			q := u.Query()
			q.Del("page")
			u.RawQuery = q.Encode()
			http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
			return
		}

		var err error
		ctx.Page, err = strconv.ParseInt(page, 10, 64)
		if err != nil || ctx.Page <= 0 || page[0] == '0' {
			handleError(w, &ctx, nil, http.StatusNotFound)
			return
		}
	}

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
