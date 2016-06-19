// +build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	if err := ioutil.WriteFile(strings.ToLower(os.Args[1])+"_id.go", []byte(fmt.Sprintf(`// This file is generated. DO NOT EDIT.

package model // import "github.com/BenLubar/webscale/model"

import (
	"database/sql"
	"database/sql/driver"

	"github.com/pkg/errors"
)

// %[1]sID is the ID of a %[1]s.
type %[1]sID ID

// Scan implements sql.Scanner.
func (id *%[1]sID) Scan(value interface{}) error {
	return (*ID)(id).Scan(value)
}

// Value implements driver.Valuer.
func (id %[1]sID) Value() (driver.Value, error) {
	return ID(id).Value()
}

// %[1]sIDs a slice of %[1]s IDs.
type %[1]sIDs []%[1]sID

// Scan implements sql.Scanner.
func (ids *%[1]sIDs) Scan(value interface{}) error {
	var generic IDs
	if err := generic.Scan(value); err != nil {
		return err
	}

	*ids = make(%[1]sIDs, len(generic))
	for i, id := range generic {
		(*ids)[i] = %[1]sID(id)
	}

	return nil
}

// Value implements driver.Valuer.
func (ids %[1]sIDs) Value() (driver.Value, error) {
	generic := make(IDs, len(ids))
	for i, id := range ids {
		generic[i] = ID(id)
	}

	return generic.Value()
}

func scan%[1]sRows(rows *sql.Rows, err error) ([]*%[1]s, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var values []*%[1]s

	for rows.Next() {
		v, err := scan%[1]s(rows)
		if err != nil {
			return values, err
		}
		values = append(values, v)
	}

	return values, rows.Close()
}

// Get retrieves the %[1]s from the database.
func (id %[1]sID) Get(ctx *Context) (*%[1]s, error) {
	v, err := scan%[1]s(ctx.Tx.QueryRow(idGet%[1]s, ctx.CurrentUser, ctx.Sudo, id))
	return v, errors.Wrap(err, "get %[1]s by ID")
}

// Get retrieves each %[1]s from the database. The returned slice has the same
// order as the ids. If any %[1]s is missing from the database, a non-nil error
// will be returned.
func (ids %[1]sIDs) Get(ctx *Context) ([]*%[1]s, error) {
	values, err := scan%[1]sRows(ctx.Tx.Query(idsGet%[1]s, ctx.CurrentUser, ctx.Sudo, ids))

	sorted := make([]*%[1]s, len(ids))
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

	return sorted, errors.Wrap(err, "get %[1]s by ID")
}
`, os.Args[1])), 0664); err != nil {
		panic(err)
	}
}
