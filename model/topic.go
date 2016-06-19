//go:generate go run mkid.go Topic

package model // import "github.com/BenLubar/webscale/model"

import "time"

type Topic struct {
	ID        TopicID
	Name      string
	Slug      string
	Author    UserID
	Category  CategoryID
	CreatedAt time.Time
	BumpedAt  time.Time
}
