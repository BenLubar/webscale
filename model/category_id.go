// This file is generated. DO NOT EDIT.

package model // import "github.com/BenLubar/webscale/model"

import (
	"database/sql"
	"database/sql/driver"

	"github.com/pkg/errors"
)

// CategoryID is the ID of a Category.
type CategoryID ID

// Scan implements sql.Scanner.
func (id *CategoryID) Scan(value interface{}) error {
	return (*ID)(id).Scan(value)
}

// Value implements driver.Valuer.
func (id CategoryID) Value() (driver.Value, error) {
	return ID(id).Value()
}

// CategoryIDs a slice of Category IDs.
type CategoryIDs []CategoryID

// Scan implements sql.Scanner.
func (ids *CategoryIDs) Scan(value interface{}) error {
	var generic IDs
	if err := generic.Scan(value); err != nil {
		return err
	}

	*ids = make(CategoryIDs, len(generic))
	for i, id := range generic {
		(*ids)[i] = CategoryID(id)
	}

	return nil
}

// Value implements driver.Valuer.
func (ids CategoryIDs) Value() (driver.Value, error) {
	generic := make(IDs, len(ids))
	for i, id := range ids {
		generic[i] = ID(id)
	}

	return generic.Value()
}

func scanCategoryRows(rows *sql.Rows, err error) ([]*Category, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var values []*Category

	for rows.Next() {
		v, err := scanCategory(rows)
		if err != nil {
			return values, err
		}
		values = append(values, v)
	}

	return values, rows.Close()
}

// Get retrieves the Category from the database.
func (id CategoryID) Get(ctx *Context) (*Category, error) {
	v, err := scanCategory(ctx.Tx.QueryRow(idGetCategory, ctx.CurrentUser, ctx.Sudo, id))
	return v, errors.Wrap(err, "get Category by ID")
}

// Get retrieves each Category from the database. The returned slice has the same
// order as the ids. If any Category is missing from the database, a non-nil error
// will be returned.
func (ids CategoryIDs) Get(ctx *Context) ([]*Category, error) {
	values, err := scanCategoryRows(ctx.Tx.Query(idsGetCategory, ctx.CurrentUser, ctx.Sudo, ids))

	sorted := make([]*Category, len(ids))
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

	return sorted, errors.Wrap(err, "get Category by ID")
}
