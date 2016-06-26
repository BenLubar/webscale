//go:generate go run mkid.go Topic

package model // import "github.com/BenLubar/webscale/model"

import (
	"database/sql"
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

var topicPosts = db.Prepare(`select ` + postFields + ` from posts as p where can_topic($1::bigint, 'topic-meta', $2::boolean, $3::bigint) and can_post($1::bigint, 'post-meta', $2::boolean, p.id) and p.topic_id = $3::bigint order by p.id asc limit $5::bigint offset $4::bigint * $5::bigint;`)

func (t *Topic) Posts(ctx *Context, page int64) ([]*Post, error) {
	posts, err := scanPostRows(ctx.Tx.Query(topicPosts, ctx.CurrentUser, ctx.Sudo, t.ID, page, perPage))
	return posts, errors.Wrapf(err, "list posts in topic %d (page %d)", t.ID, page)
}

var topicPostsPageCount = db.Prepare(`select ` + pageCountField + ` from posts as p where can_topic($1::bigint, 'topic-meta', $2::boolean, $3::bigint) and can_post($1::bigint, 'post-meta', $2::boolean, p.id) and p.topic_id = $3::bigint;`)

func (t *Topic) PostsPageCount(ctx *Context) (int64, error) {
	var count int64
	err := ctx.Tx.QueryRow(topicPostsPageCount, ctx.CurrentUser, ctx.Sudo, t.ID, perPage).Scan(&count)
	return count, errors.Wrapf(err, "count pages in topic %d", t.ID)
}

var topicsLastPosts = db.Prepare(`select ` + postFields + ` from posts as p where p.id = (select p2.id from posts as p2 where can_post($1::bigint, 'post-meta', $2::boolean, p2.id) and p2.topic_id = p.topic_id order by p2.created_at desc, p2.id asc limit 1) and can_topic($1::bigint, 'topic-meta', $2::boolean, p.topic_id) and array[p.topic_id] <@ $3::bigint[] order by p.topic_id asc;`)

func (ids TopicIDs) LastPosts(ctx *Context) ([]*Post, error) {
	values, err := scanPostRows(ctx.Tx.Query(topicsLastPosts, ctx.CurrentUser, ctx.Sudo, ids))

	sorted := make([]*Post, len(ids))
search:
	for i, id := range ids {
		if id == 0 {
			continue
		}

		for _, v := range values {
			if v.Topic == id {
				sorted[i] = v
				continue search
			}
		}

		if err != nil {
			err = sql.ErrNoRows
		}
	}

	return sorted, errors.Wrap(err, "get last post in topics")
}
