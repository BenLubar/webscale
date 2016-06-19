package controller // import "github.com/BenLubar/webscale/controller"

import (
	"net/http"

	"github.com/BenLubar/webscale/model"
	"github.com/BenLubar/webscale/view"
)

func index(w http.ResponseWriter, ctx *model.Context) error {
	if ctx.Request.URL.Path != "/" {
		handleError(w, ctx, nil, http.StatusNotFound)
		return nil
	}

	categories, err := model.TopLevelCategories(ctx)
	if err != nil {
		return err
	}

	ctx.Header.Title = ""
	ctx.Header.Breadcrumb = []model.Breadcrumb{
		{
			Name: "#webscale",
			Path: "/",
		},
	}

	err = ctx.Tx.Commit()
	if err != nil {
		return err
	}

	return view.Index.Execute(w, ctx, http.StatusOK, struct {
		Categories []*model.Category
	}{
		Categories: categories,
	})
}

func init() {
	http.Handle("/", handler(index))
}
