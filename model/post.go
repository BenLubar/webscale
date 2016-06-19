//go:generate go run mkid.go Post
//go:generate go run mkid.go PostRevision

package model // import "github.com/BenLubar/webscale/model"

import (
	"time"

	"github.com/BenLubar/webscale/db"
)

type Post struct {
	ID        PostID
	Topic     TopicID
	Author    UserID
	Parent    PostID
	CreatedAt time.Time
	Content   string
	Tags      Strings
}

const postFields = `p.id, p.topic_id, case when can_post($1::bigint, 'post-view-author', $2::boolean, p.id) then p.user_id else null end as user_id, p.parent_post_id, p.created_at, p.content, case when can_post($1::bigint, 'post-view-tags', $2::boolean, p.id) then p.tags else array[]::citext[] end as tags`

func scanPost(s scanner) (*Post, error) {
	var p Post
	if err := s.Scan(&p.ID, &p.Topic, &p.Author, &p.Parent, &p.CreatedAt, &p.Content, &p.Tags); err != nil {
		return nil, err
	}
	return &p, nil
}

var idGetPost = db.Prepare(`select ` + postFields + ` from posts as p where can_post($1::bigint, 'post-meta', $2::boolean, p.id) and p.id = $3::bigint order by p.id asc;`)
var idsGetPost = db.Prepare(`select ` + postFields + ` from posts as p where can_post($1::bigint, 'post-meta', $2::boolean, p.id) and array[p.id] <@ $3::bigint[] order by p.id asc;`)

type PostRevision struct {
	ID        PostRevisionID
	Post      PostID
	Author    UserID
	CreatedAt time.Time
	Content   string
	Tags      Strings
}

const postRevisionFields = `pr.id, pr.post_id, case when can_post($1::bigint, 'post-view-author', $2::boolean, pr.post_id) then pr.user_id else null end as user_id, pr.created_at, pr.content, case when can_post($1::bigint, 'post-view-tags', $2::boolean, pr.post_id) then pr.tags else array[]::citext[] end as tags`

func scanPostRevision(s scanner) (*PostRevision, error) {
	var pr PostRevision
	if err := s.Scan(&pr.ID, &pr.Post, &pr.Author, &pr.CreatedAt, &pr.Content, &pr.Tags); err != nil {
		return nil, err
	}
	return &pr, nil
}

var idGetPostRevision = db.Prepare(`select ` + postRevisionFields + ` from post_revisions as pr where can_post($1::bigint, 'post-view-history', $2::boolean, pr.post_id) and pr.id = $3::bigint order by pr.id asc;`)
var idsGetPostRevision = db.Prepare(`select ` + postRevisionFields + ` from post_revisions as pr where can_post($1::bigint, 'post-view-history', $2::boolean, pr.post_id) and array[pr.id] <@ $3::bigint[] order by pr.id asc;`)
