// This file is generated. DO NOT EDIT.

package model // import "github.com/BenLubar/webscale/model"

import (
	"database/sql"
	"database/sql/driver"

	"github.com/pkg/errors"
)

// UserID is the ID of a User.
type UserID ID

// Scan implements sql.Scanner.
func (id *UserID) Scan(value interface{}) error {
	return (*ID)(id).Scan(value)
}

// Value implements driver.Valuer.
func (id UserID) Value() (driver.Value, error) {
	return ID(id).Value()
}

// UserIDs a slice of User IDs.
type UserIDs []UserID

// Scan implements sql.Scanner.
func (ids *UserIDs) Scan(value interface{}) error {
	var generic IDs
	if err := generic.Scan(value); err != nil {
		return err
	}

	*ids = make(UserIDs, len(generic))
	for i, id := range generic {
		(*ids)[i] = UserID(id)
	}

	return nil
}

// Value implements driver.Valuer.
func (ids UserIDs) Value() (driver.Value, error) {
	generic := make(IDs, len(ids))
	for i, id := range ids {
		generic[i] = ID(id)
	}

	return generic.Value()
}

func scanUserRows(rows *sql.Rows, err error) ([]*User, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var values []*User

	for rows.Next() {
		v, err := scanUser(rows)
		if err != nil {
			return values, err
		}
		values = append(values, v)
	}

	return values, rows.Close()
}

// Get retrieves the User from the database.
func (id UserID) Get(ctx *Context) (*User, error) {
	v, err := scanUser(ctx.Tx.QueryRow(idGetUser, ctx.CurrentUser, ctx.Sudo, id))
	return v, errors.Wrap(err, "get User by ID")
}

// Get retrieves each User from the database. The returned slice has the same
// order as the ids. If any User is missing from the database, a non-nil error
// will be returned.
func (ids UserIDs) Get(ctx *Context) ([]*User, error) {
	values, err := scanUserRows(ctx.Tx.Query(idsGetUser, ctx.CurrentUser, ctx.Sudo, ids))

	sorted := make([]*User, len(ids))
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

	return sorted, errors.Wrap(err, "get User by ID")
}
