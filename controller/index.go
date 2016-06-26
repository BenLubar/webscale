package controller // import "github.com/BenLubar/webscale/controller"

import (
	"net/http"

	"github.com/BenLubar/webscale/model"
	"github.com/BenLubar/webscale/model/helpers"
	"github.com/BenLubar/webscale/view"
)

func index(w http.ResponseWriter, ctx *model.Context) error {
	if ctx.Request.URL.Path != "/" {
		handleError(w, ctx, nil, http.StatusNotFound)
		return nil
	}

	var err error
	var data struct {
		Categories []helpers.CategoryWithLatestTopic
		Latest     []helpers.TopicWithLastPost
	}

	if ctx.Page == 0 {
		var categories []*model.Category
		if categories, err = model.TopLevelCategories(ctx); err != nil {
			return err
		}

		if data.Categories, err = helpers.CategoriesLatestTopics(ctx, categories); err != nil {
			return err
		}
	}

	var topics []*model.Topic
	if topics, err = model.LatestTopics(ctx, ctx.Page); err != nil {
		return err
	}
	if data.Latest, err = helpers.TopicsLastPosts(ctx, topics); err != nil {
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
