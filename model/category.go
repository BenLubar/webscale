//go:generate go run mkid.go Category

package model // import "github.com/BenLubar/webscale/model"

import (
	"github.com/BenLubar/webscale/db"
	"github.com/pkg/errors"
)

type Category struct {
	ID     CategoryID
	Name   string
	Slug   string
	Parent CategoryID
	Path   CategoryIDs
}

const categoryFields = `c.id, c.name, c.slug, c.parent_category_id, cp.path`

func scanCategory(s scanner) (*Category, error) {
	var c Category
	if err := s.Scan(&c.ID, &c.Name, &c.Slug, &c.Parent, &c.Path); err != nil {
		return nil, err
	}
	return &c, nil
}

var idGetCategory = db.Prepare(`select ` + categoryFields + ` from categories as c inner join categories_path as cp on cp.category_id = c.id where can_category($1::bigint, 'category-meta', $2::boolean, c.id) and c.id = $3::bigint order by c.id asc;`)
var idsGetCategory = db.Prepare(`select ` + categoryFields + ` from categories as c inner join categories_path as cp on cp.category_id = c.id where can_category($1::bigint, 'category-meta', $2::boolean, c.id) and array[c.id] <@ $3::bigint[] order by c.id asc;`)

var categoryTopLevelCategories = db.Prepare(`select ` + categoryFields + ` from categories as c inner join categories_path as cp on cp.category_id = c.id where can_category($1::bigint, 'category-meta', $2::boolean, c.id) and c.parent_category_id is null order by c.id asc;`)

func TopLevelCategories(ctx *Context) ([]*Category, error) {
	categories, err := scanCategoryRows(ctx.Tx.Query(categoryTopLevelCategories, ctx.CurrentUser, ctx.Sudo))
	return categories, errors.Wrap(err, "list top level categories")
}

var categoryChildren = db.Prepare(`select ` + categoryFields + ` from categories as c inner join categories_path as cp on cp.category_id = c.id where can_category($1::bigint, 'category-meta', $2::boolean, c.id) and c.parent_category_id = $3::bigint order by c.id asc;`)

func (c *Category) Children(ctx *Context) ([]*Category, error) {
	categories, err := scanCategoryRows(ctx.Tx.Query(categoryChildren, ctx.CurrentUser, ctx.Sudo, c.ID))
	return categories, errors.Wrapf(err, "list children of category %d", c.ID)
}

var categoryTopics = db.Prepare(`select ` + topicFields + ` from topics as t where can_category($1::bigint, 'category-list-topics', $2::boolean, $3::bigint) and can_topic($1::bigint, 'topic-meta', $2::boolean, t.id) and t.category_id = $3::bigint order by t.bumped_at desc, t.id asc limit $5::bigint offset $4::bigint * $5::bigint;`)

func (c *Category) Topics(ctx *Context, page int64) ([]*Topic, error) {
	topics, err := scanTopicRows(ctx.Tx.Query(categoryTopics, ctx.CurrentUser, ctx.Sudo, c.ID, page, perPage))
	return topics, errors.Wrapf(err, "list topics in category %d (page %d)", c.ID, page)
}
