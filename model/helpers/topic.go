package helpers // import "github.com/BenLubar/webscale/model/helpers"

import (
	"database/sql"

	"github.com/BenLubar/webscale/model"
	"github.com/pkg/errors"
)

type TopicWithLastPost struct {
	Topic       *model.Topic
	TopicAuthor *model.User
	PostWithAuthor
}

func TopicsLastPosts(ctx *model.Context, topics []*model.Topic) ([]TopicWithLastPost, error) {
	tlp := make([]TopicWithLastPost, len(topics))

	tids := make(model.TopicIDs, len(topics))
	uids := make(model.UserIDs, len(topics))

	for i, t := range topics {
		if t != nil {
			tlp[i].Topic = t
			tids[i] = t.ID
			uids[i] = t.Author
		}
	}

	users, err := uids.Get(ctx)
	if errors.Cause(err) == sql.ErrNoRows {
		err = nil
	} else if err != nil {
		return tlp, errors.Wrap(err, "get topics authors")
	}

	for i, u := range users {
		tlp[i].TopicAuthor = u
	}

	posts, err := tids.LastPosts(ctx)
	if errors.Cause(err) == sql.ErrNoRows {
		err = nil
	} else if err != nil {
		return tlp, errors.Wrap(err, "get topics last posts")
	}

	pa, err := PostsAuthors(ctx, posts)
	if err != nil {
		return tlp, errors.Wrap(err, "get topics last posts authors")
	}

	for i, p := range pa {
		tlp[i].PostWithAuthor = p
	}

	return tlp, nil
}
