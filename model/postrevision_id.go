// This file is generated. DO NOT EDIT.

package model // import "github.com/BenLubar/webscale/model"

import (
	"database/sql"
	"database/sql/driver"

	"github.com/pkg/errors"
)

// PostRevisionID is the ID of a PostRevision.
type PostRevisionID ID

// Scan implements sql.Scanner.
func (id *PostRevisionID) Scan(value interface{}) error {
	return (*ID)(id).Scan(value)
}

// Value implements driver.Valuer.
func (id PostRevisionID) Value() (driver.Value, error) {
	return ID(id).Value()
}

// PostRevisionIDs a slice of PostRevision IDs.
type PostRevisionIDs []PostRevisionID

// Scan implements sql.Scanner.
func (ids *PostRevisionIDs) Scan(value interface{}) error {
	var generic IDs
	if err := generic.Scan(value); err != nil {
		return err
	}

	*ids = make(PostRevisionIDs, len(generic))
	for i, id := range generic {
		(*ids)[i] = PostRevisionID(id)
	}

	return nil
}

// Value implements driver.Valuer.
func (ids PostRevisionIDs) Value() (driver.Value, error) {
	generic := make(IDs, len(ids))
	for i, id := range ids {
		generic[i] = ID(id)
	}

	return generic.Value()
}

func scanPostRevisionRows(rows *sql.Rows, err error) ([]*PostRevision, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var values []*PostRevision

	for rows.Next() {
		v, err := scanPostRevision(rows)
		if err != nil {
			return values, err
		}
		values = append(values, v)
	}

	return values, rows.Close()
}

// Get retrieves the PostRevision from the database.
func (id PostRevisionID) Get(ctx *Context) (*PostRevision, error) {
	v, err := scanPostRevision(ctx.Tx.QueryRow(idGetPostRevision, ctx.CurrentUser, ctx.Sudo, id))
	return v, errors.Wrap(err, "get PostRevision by ID")
}

// Get retrieves each PostRevision from the database. The returned slice has the same
// order as the ids. If any PostRevision is missing from the database, a non-nil error
// will be returned.
func (ids PostRevisionIDs) Get(ctx *Context) ([]*PostRevision, error) {
	values, err := scanPostRevisionRows(ctx.Tx.Query(idsGetPostRevision, ctx.CurrentUser, ctx.Sudo, ids))

	sorted := make([]*PostRevision, len(ids))
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

	return sorted, errors.Wrap(err, "get PostRevision by ID")
}
