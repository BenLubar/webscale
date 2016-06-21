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

	var err error
	var data struct {
		Categories []*model.Category
		Latest     []*model.Topic
	}

	if ctx.Page == 0 {
		if data.Categories, err = model.TopLevelCategories(ctx); err != nil {
			return err
		}
	}

	if data.Latest, err = model.LatestTopics(ctx, ctx.Page); err != nil {
		return err
	}

	if ctx.Page != 0 && len(data.Latest) == 0 {
		handleError(w, ctx, nil, http.StatusNotFound)
		return nil
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

	return view.Index.Execute(w, ctx, http.StatusOK, data)
}

func init() {
	http.Handle("/", handler(index))
}
