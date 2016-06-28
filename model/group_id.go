// This file is generated. DO NOT EDIT.

package model // import "github.com/BenLubar/webscale/model"

import (
	"database/sql"
	"database/sql/driver"

	"github.com/pkg/errors"
)

// GroupID is the ID of a Group.
type GroupID ID

// Scan implements sql.Scanner.
func (id *GroupID) Scan(value interface{}) error {
	return (*ID)(id).Scan(value)
}

// Value implements driver.Valuer.
func (id GroupID) Value() (driver.Value, error) {
	return ID(id).Value()
}

// GroupIDs a slice of Group IDs.
type GroupIDs []GroupID

// Scan implements sql.Scanner.
func (ids *GroupIDs) Scan(value interface{}) error {
	var generic IDs
	if err := generic.Scan(value); err != nil {
		return err
	}

	*ids = make(GroupIDs, len(generic))
	for i, id := range generic {
		(*ids)[i] = GroupID(id)
	}

	return nil
}

// Value implements driver.Valuer.
func (ids GroupIDs) Value() (driver.Value, error) {
	generic := make(IDs, len(ids))
	for i, id := range ids {
		generic[i] = ID(id)
	}

	return generic.Value()
}

func scanGroupRows(rows *sql.Rows, err error) ([]*Group, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var values []*Group

	for rows.Next() {
		v, err := scanGroup(rows)
		if err != nil {
			return values, err
		}
		values = append(values, v)
	}

	return values, rows.Close()
}

// Get retrieves the Group from the database.
func (id GroupID) Get(ctx *Context) (*Group, error) {
	v, err := scanGroup(ctx.Tx.QueryRow(idGetGroup, ctx.CurrentUser, ctx.Sudo, id))
	return v, errors.Wrap(err, "get Group by ID")
}

// Get retrieves each Group from the database. The returned slice has the same
// order as the ids. If any Group is missing from the database, a non-nil error
// will be returned.
func (ids GroupIDs) Get(ctx *Context) ([]*Group, error) {
	values, err := scanGroupRows(ctx.Tx.Query(idsGetGroup, ctx.CurrentUser, ctx.Sudo, ids))

	sorted := make([]*Group, len(ids))
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

		if err == nil {
			err = sql.ErrNoRows
		}
	}

	return sorted, errors.Wrap(err, "get Group by ID")
}
