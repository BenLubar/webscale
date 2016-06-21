//go:generate go run mkid.go Topic

package model // import "github.com/BenLubar/webscale/model"

import (
	"time"

	"github.com/BenLubar/webscale/db"
	"github.com/pkg/errors"
)

type Topic struct {
	ID        TopicID
	Name      string
	Slug      string
	Author    UserID
	Category  CategoryID
	CreatedAt time.Time
	BumpedAt  time.Time
}

const topicFields = `t.id, t.name, t.slug, case when can_topic($1::bigint, 'topic-view-author', $2::boolean, t.id) then t.user_id else null end as user_id, t.category_id, t.created_at, t.bumped_at`

func scanTopic(s scanner) (*Topic, error) {
	var t Topic
	if err := s.Scan(&t.ID, &t.Name, &t.Slug, &t.Author, &t.Category, &t.CreatedAt, &t.BumpedAt); err != nil {
		return nil, err
	}
	return &t, nil
}

var idGetTopic = db.Prepare(`select ` + topicFields + ` from topics as t where can_topic($1::bigint, 'topic-meta', $2::boolean, t.id) and t.id = $3::bigint order by t.id asc;`)
var idsGetTopic = db.Prepare(`select ` + topicFields + ` from topics as t where can_topic($1::bigint, 'topic-meta', $2::boolean, t.id) and array[t.id] <@ $3::bigint[] order by t.id asc;`)

var topicLatestTopics = db.Prepare(`select ` + topicFields + ` from topics as t where can_topic($1::bigint, 'topic-meta', $2::boolean, t.id) order by t.bumped_at desc, t.id asc limit $4::bigint offset $3::bigint * $4::bigint;`)

func LatestTopics(ctx *Context, page int64) ([]*Topic, error) {
	topics, err := scanTopicRows(ctx.Tx.Query(topicLatestTopics, ctx.CurrentUser, ctx.Sudo, page, perPage))
	return topics, errors.Wrapf(err, "get latest topics (page %d)", page)
}
