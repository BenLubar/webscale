package helpers // import "github.com/BenLubar/webscale/model/helpers"

import (
	"database/sql"

	"github.com/BenLubar/webscale/model"
	"github.com/pkg/errors"
)

type CategoryWithLatestTopic struct {
	Category *model.Category
	TopicWithLastPost
}

func CategoriesLatestTopics(ctx *model.Context, categories []*model.Category) ([]CategoryWithLatestTopic, error) {
	clt := make([]CategoryWithLatestTopic, len(categories))

	cids := make(model.CategoryIDs, len(categories))
	for i, c := range categories {
		if c != nil {
			clt[i].Category = c
			cids[i] = c.ID
		}
	}

	topics, err := cids.LatestTopics(ctx)
	if errors.Cause(err) == sql.ErrNoRows {
		err = nil
	} else if err != nil {
		return clt, errors.Wrap(err, "get categories latest topics")
	}

	tlp, err := TopicsLastPosts(ctx, topics)
	if err != nil {
		return clt, errors.Wrap(err, "get categories latest topics last posts")
	}

	for i, tp := range tlp {
		clt[i].TopicWithLastPost = tp
	}

	return clt, nil
}
