package controller // import "github.com/BenLubar/webscale/controller"

import (
	"log"
	"net/http"

	"github.com/BenLubar/webscale/model"
	"github.com/BenLubar/webscale/proxy"
	"github.com/BenLubar/webscale/view"
	"github.com/pkg/errors"
)

func handleError(w http.ResponseWriter, ctx *model.Context, err error, status int) {
	if status != http.StatusNotFound {
		log.Printf("%v %v [ERROR] %+v", proxy.RequestIP(ctx.Request), ctx.Request.URL, err)
	}

	if err == nil {
		err = errors.New(http.StatusText(status))
	}

	ctx.Header.Title = "Error"
	ctx.Header.Breadcrumb = nil

	_ = view.Error.Execute(w, ctx, status, struct {
		Error error
	}{
		Error: err,
	})
}
