//go:generate go run mkid.go Post
//go:generate go run mkid.go PostRevision

package model // import "github.com/BenLubar/webscale/model"

import "time"

type Post struct {
	ID        PostID
	Topic     TopicID
	Author    UserID
	Parent    PostID
	CreatedAt time.Time
	Content   string
	Tags      []string
}

type PostRevision struct {
	ID        PostRevisionID
	Topic     PostID
	Author    UserID
	CreatedAt time.Time
	Content   string
	Tags      []string
}
