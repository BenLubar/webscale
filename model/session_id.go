package model // import "github.com/BenLubar/webscale/model"

import (
	"database/sql"
	"database/sql/driver"

	"github.com/pkg/errors"
)

// SessionID is the ID of a Session.
type SessionID UUID

// Scan implements sql.Scanner.
func (id *SessionID) Scan(value interface{}) error {
	return (*UUID)(id).Scan(value)
}

// Value implements driver.Valuer.
func (id SessionID) Value() (driver.Value, error) {
	return UUID(id).Value()
}

func scanSessionRows(rows *sql.Rows, err error) ([]*Session, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var values []*Session

	for rows.Next() {
		v, err := scanSession(rows)
		if err != nil {
			return values, err
		}
		values = append(values, v)
	}

	return values, rows.Close()
}

// Get retrieves the Session from the database.
func (id SessionID) Get(ctx *Context) (*Session, error) {
	v, err := scanSession(ctx.Tx.QueryRow(idGetSession, ctx.CurrentUser, ctx.Sudo, id))
	return v, errors.Wrap(err, "get Session by ID")
}
