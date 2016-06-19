//go:generate go run mkid.go Category

package model // import "github.com/BenLubar/webscale/model"

import (
	"database/sql"

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

func categoryScan(s scanner) (*Category, error) {
	var c Category
	if err := s.Scan(&c.ID, &c.Name, &c.Slug, &c.Parent, &c.Path); err != nil {
		return nil, errors.Wrap(err, "scan category")
	}
	return &c, nil
}

func categoriesScan(rows *sql.Rows, err error) ([]*Category, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*Category

	for rows.Next() {
		c, err := categoryScan(rows)
		if err != nil {
			return categories, err
		}
		categories = append(categories, c)
	}

	return categories, rows.Close()
}

var categoryIDGet = db.Prepare(`select ` + categoryFields + ` from categories as c inner join categories_path as cp on cp.category_id = c.id where can_category($1::bigint, 'category-meta', $2::boolean, c.id) and c.id = $3::bigint order by c.id asc;`)

func (id CategoryID) Get(ctx *Context) (*Category, error) {
	c, err := categoryScan(ctx.Tx.QueryRow(categoryIDGet, ctx.CurrentUser, ctx.Sudo, id))
	return c, errors.Wrap(err, "get category by ID")
}

var categoryTopLevelCategories = db.Prepare(`select ` + categoryFields + ` from categories as c inner join categories_path as cp on cp.category_id = c.id where can_category($1::bigint, 'category-meta', $2::boolean, c.id) and c.parent_category_id is null order by c.id asc;`)

func TopLevelCategories(ctx *Context) ([]*Category, error) {
	categories, err := categoriesScan(ctx.Tx.Query(categoryTopLevelCategories, ctx.CurrentUser, ctx.Sudo))
	return categories, errors.Wrap(err, "list top level categories")
}

var categoryChildren = db.Prepare(`select ` + categoryFields + ` from categories as c inner join categories_path as cp on cp.category_id = c.id where can_category($1::bigint, 'category-meta', $2::boolean, c.id) and c.parent_category_id = $3::bigint order by c.id asc;`)

func (c *Category) Children(ctx *Context) ([]*Category, error) {
	categories, err := categoriesScan(ctx.Tx.Query(categoryChildren, ctx.CurrentUser, ctx.Sudo, c.ID))
	return categories, errors.Wrapf(err, "list children of category %d", c.ID)
}
