package helpers // import "github.com/BenLubar/webscale/model/helpers"

import (
	"database/sql"

	"github.com/BenLubar/webscale/model"
	"github.com/pkg/errors"
)

type PostWithAuthor struct {
	Post       *model.Post
	PostAuthor *model.User
}

func PostsAuthors(ctx *model.Context, posts []*model.Post) ([]PostWithAuthor, error) {
	pa := make([]PostWithAuthor, len(posts))

	uids := make(model.UserIDs, len(posts))

	for i, p := range posts {
		if p != nil {
			pa[i].Post = p
			uids[i] = p.Author
		}
	}

	users, err := uids.Get(ctx)
	if errors.Cause(err) == sql.ErrNoRows {
		err = nil
	} else if err != nil {
		return pa, errors.Wrap(err, "get post authors")
	}

	for i, u := range users {
		pa[i].PostAuthor = u
	}

	return pa, nil
}
