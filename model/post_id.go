// This file is generated. DO NOT EDIT.

package model // import "github.com/BenLubar/webscale/model"

import (
	"database/sql"
	"database/sql/driver"

	"github.com/pkg/errors"
)

// PostID is the ID of a Post.
type PostID ID

// Scan implements sql.Scanner.
func (id *PostID) Scan(value interface{}) error {
	return (*ID)(id).Scan(value)
}

// Value implements driver.Valuer.
func (id PostID) Value() (driver.Value, error) {
	return ID(id).Value()
}

// PostIDs a slice of Post IDs.
type PostIDs []PostID

// Scan implements sql.Scanner.
func (ids *PostIDs) Scan(value interface{}) error {
	var generic IDs
	if err := generic.Scan(value); err != nil {
		return err
	}

	*ids = make(PostIDs, len(generic))
	for i, id := range generic {
		(*ids)[i] = PostID(id)
	}

	return nil
}

// Value implements driver.Valuer.
func (ids PostIDs) Value() (driver.Value, error) {
	generic := make(IDs, len(ids))
	for i, id := range ids {
		generic[i] = ID(id)
	}

	return generic.Value()
}

func scanPostRows(rows *sql.Rows, err error) ([]*Post, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var values []*Post

	for rows.Next() {
		v, err := scanPost(rows)
		if err != nil {
			return values, err
		}
		values = append(values, v)
	}

	return values, rows.Close()
}

// Get retrieves the Post from the database.
func (id PostID) Get(ctx *Context) (*Post, error) {
	v, err := scanPost(ctx.Tx.QueryRow(idGetPost, ctx.CurrentUser, ctx.Sudo, id))
	return v, errors.Wrap(err, "get Post by ID")
}

// Get retrieves each Post from the database. The returned slice has the same
// order as the ids. If any Post is missing from the database, a non-nil error
// will be returned.
func (ids PostIDs) Get(ctx *Context) ([]*Post, error) {
	values, err := scanPostRows(ctx.Tx.Query(idsGetPost, ctx.CurrentUser, ctx.Sudo, ids))

	sorted := make([]*Post, len(ids))
search:
	for i, id := range ids {
		if id == 0 {
			continue
		}

		for _, v := range values {
			if v.ID == id {
				sorted[i] = v
				continue search
			}
		}

		if err != nil {
			err = sql.ErrNoRows
		}
	}

	return sorted, errors.Wrap(err, "get Post by ID")
}
