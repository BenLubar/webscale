// This file is generated. DO NOT EDIT.

package model // import "github.com/BenLubar/webscale/model"

import (
	"database/sql"
	"database/sql/driver"

	"github.com/pkg/errors"
)

// PermissionID is the ID of a Permission.
type PermissionID ID

// Scan implements sql.Scanner.
func (id *PermissionID) Scan(value interface{}) error {
	return (*ID)(id).Scan(value)
}

// Value implements driver.Valuer.
func (id PermissionID) Value() (driver.Value, error) {
	return ID(id).Value()
}

// PermissionIDs a slice of Permission IDs.
type PermissionIDs []PermissionID

// Scan implements sql.Scanner.
func (ids *PermissionIDs) Scan(value interface{}) error {
	var generic IDs
	if err := generic.Scan(value); err != nil {
		return err
	}

	*ids = make(PermissionIDs, len(generic))
	for i, id := range generic {
		(*ids)[i] = PermissionID(id)
	}

	return nil
}

// Value implements driver.Valuer.
func (ids PermissionIDs) Value() (driver.Value, error) {
	generic := make(IDs, len(ids))
	for i, id := range ids {
		generic[i] = ID(id)
	}

	return generic.Value()
}

func scanPermissionRows(rows *sql.Rows, err error) ([]*Permission, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var values []*Permission

	for rows.Next() {
		v, err := scanPermission(rows)
		if err != nil {
			return values, err
		}
		values = append(values, v)
	}

	return values, rows.Close()
}

// Get retrieves the Permission from the database.
func (id PermissionID) Get(ctx *Context) (*Permission, error) {
	v, err := scanPermission(ctx.Tx.QueryRow(idGetPermission, ctx.CurrentUser, ctx.Sudo, id))
	return v, errors.Wrap(err, "get Permission by ID")
}

// Get retrieves each Permission from the database. The returned slice has the same
// order as the ids. If any Permission is missing from the database, a non-nil error
// will be returned.
func (ids PermissionIDs) Get(ctx *Context) ([]*Permission, error) {
	values, err := scanPermissionRows(ctx.Tx.Query(idsGetPermission, ctx.CurrentUser, ctx.Sudo, ids))

	sorted := make([]*Permission, len(ids))
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

	return sorted, errors.Wrap(err, "get Permission by ID")
}
