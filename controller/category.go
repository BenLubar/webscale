package controller // import "github.com/BenLubar/webscale/controller"

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/BenLubar/webscale/model"
	"github.com/BenLubar/webscale/view"
	"github.com/pkg/errors"
)

func category(w http.ResponseWriter, ctx *model.Context) error {
	path := strings.Split(strings.TrimPrefix(ctx.Request.URL.Path, "/category/"), "/")
	if len(path) < 1 || len(path) > 2 {
		handleError(w, ctx, nil, http.StatusNotFound)
		return nil
	}

	var category *model.Category

	if n, err := strconv.ParseInt(path[0], 10, 64); err != nil || n == 0 {
		handleError(w, ctx, nil, http.StatusNotFound)
		return nil
	} else if category, err = model.CategoryID(n).Get(ctx); err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			handleError(w, ctx, nil, http.StatusNotFound)
			return nil
		}
		return err
	}

	if len(path) < 2 || path[1] != category.Slug {
		http.Redirect(w, ctx.Request, fmt.Sprintf("/category/%d/%s", category.ID, category.Slug), http.StatusMovedPermanently)
		return nil
	}

	children, err := category.Children(ctx)
	if err != nil {
		return err
	}

	ctx.Header.Title = category.Name
	ctx.Header.Breadcrumb = []model.Breadcrumb{
		{
			Name: "#webscale",
			Path: "/",
		},
	}

	for _, id := range category.Path {
		c, err := id.Get(ctx)
		if err != nil {
			return err
		}

		ctx.Header.Breadcrumb = append(ctx.Header.Breadcrumb, model.Breadcrumb{
			Name: c.Name,
			Path: fmt.Sprintf("/category/%d/%s", c.ID, c.Slug),
		})
	}

	if err = ctx.Tx.Commit(); err != nil {
		return err
	}

	return view.Category.Execute(w, ctx, http.StatusOK, struct {
		Category *model.Category
		Children []*model.Category
	}{
		Category: category,
		Children: children,
	})
}

func init() {
	http.Handle("/category/", handler(category))
}
