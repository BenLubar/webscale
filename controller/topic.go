package controller // import "github.com/BenLubar/webscale/controller"

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/BenLubar/webscale/model"
	"github.com/BenLubar/webscale/model/helpers"
	"github.com/BenLubar/webscale/view"
	"github.com/pkg/errors"
)

func topic(w http.ResponseWriter, ctx *model.Context) error {
	if ctx.Request.URL.Path == "/topic/" {
		http.Redirect(w, ctx.Request, "/", http.StatusMovedPermanently)
		return nil
	}

	path := strings.Split(strings.TrimPrefix(ctx.Request.URL.Path, "/topic/"), "/")
	if len(path) < 1 || len(path) > 2 {
		handleError(w, ctx, nil, http.StatusNotFound)
		return nil
	}

	var err error
	var data struct {
		Topic    *model.Topic
		Author   *model.User
		Category *model.Category
		Posts    []helpers.PostWithAuthor
	}

	if n, err := strconv.ParseInt(path[0], 10, 64); err != nil || n == 0 {
		handleError(w, ctx, nil, http.StatusNotFound)
		return nil
	} else if data.Topic, err = model.TopicID(n).Get(ctx); err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			handleError(w, ctx, nil, http.StatusNotFound)
			return nil
		}
		return err
	}

	if data.Author, err = data.Topic.Author.Get(ctx); err != nil && errors.Cause(err) != sql.ErrNoRows {
		return err
	}

	if len(path) < 2 || path[1] != data.Topic.Slug {
		url := fmt.Sprintf("/topic/%d/%s", data.Topic.ID, data.Topic.Slug)
		if ctx.Page != 0 {
			url = fmt.Sprintf("%s?page=%d", url, ctx.Page)
		}

		http.Redirect(w, ctx.Request, url, http.StatusMovedPermanently)
		return nil
	}

	var posts []*model.Post
	if posts, err = data.Topic.Posts(ctx, ctx.Page); err != nil {
		return err
	}
	if data.Posts, err = helpers.PostsAuthors(ctx, posts); err != nil {
		return err
	}

	if ctx.Page != 0 && len(data.Posts) == 0 {
		handleError(w, ctx, nil, http.StatusNotFound)
		return nil
	}

	ctx.Header.Title = data.Topic.Name
	ctx.Header.Breadcrumb = []model.Breadcrumb{
		{
			Name: "#webscale",
			Path: "/",
		},
	}

	if data.Category, err = data.Topic.Category.Get(ctx); err != nil {
		return err
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
	ctx.Header.Breadcrumb = append(ctx.Header.Breadcrumb, model.Breadcrumb{
		Name: data.Topic.Name,
		Path: fmt.Sprintf("/topic/%d/%s", data.Topic.ID, data.Topic.Slug),
	})

	if err = ctx.Tx.Commit(); err != nil {
		return err
	}

	return view.Topic.Execute(w, ctx, http.StatusOK, data)
}

func init() {
	http.Handle("/topic/", handler(topic))
}
