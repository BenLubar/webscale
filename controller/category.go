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
	if ctx.Request.URL.Path == "/category/" {
		http.Redirect(w, ctx.Request, "/", http.StatusMovedPermanently)
		return nil
	}

	path := strings.Split(strings.TrimPrefix(ctx.Request.URL.Path, "/category/"), "/")
	if len(path) < 1 || len(path) > 2 {
		handleError(w, ctx, nil, http.StatusNotFound)
		return nil
	}

	var err error
	var data struct {
		Category *model.Category
		Children []*model.Category
		Topics   []*model.Topic
	}

	if n, err := strconv.ParseInt(path[0], 10, 64); err != nil || n == 0 {
		handleError(w, ctx, nil, http.StatusNotFound)
		return nil
	} else if data.Category, err = model.CategoryID(n).Get(ctx); err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			handleError(w, ctx, nil, http.StatusNotFound)
			return nil
		}
		return err
	}

	if len(path) < 2 || path[1] != data.Category.Slug {
		url := fmt.Sprintf("/category/%d/%s", data.Category.ID, data.Category.Slug)
		if ctx.Page != 0 {
			url = fmt.Sprintf("%s?page=%d", url, ctx.Page)
		}

		http.Redirect(w, ctx.Request, url, http.StatusMovedPermanently)
		return nil
	}

	if data.Topics, err = data.Category.Topics(ctx, ctx.Page); err != nil {
		return err
	}

	if ctx.Page != 0 && len(data.Topics) == 0 {
		handleError(w, ctx, nil, http.StatusNotFound)
		return nil
	}

	if data.Children, err = data.Category.Children(ctx); err != nil {
		return err
	}

	ctx.Header.Title = data.Category.Name
	ctx.Header.Breadcrumb = []model.Breadcrumb{
		{
			Name: "#webscale",
			Path: "/",
		},
	}

	breadcrumb, err := data.Category.Path.Get(ctx)
	if err != nil {
		return err
	}

	for _, c := range breadcrumb {
		ctx.Header.Breadcrumb = append(ctx.Header.Breadcrumb, model.Breadcrumb{
			Name: c.Name,
			Path: fmt.Sprintf("/category/%d/%s", c.ID, c.Slug),
		})
	}

	if err = ctx.Tx.Commit(); err != nil {
		return err
	}

	return view.Category.Execute(w, ctx, http.StatusOK, data)
}

func init() {
	http.Handle("/category/", handler(category))
}
