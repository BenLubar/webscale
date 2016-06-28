// This file is generated. DO NOT EDIT.

package model // import "github.com/BenLubar/webscale/model"

import (
	"database/sql"
	"database/sql/driver"

	"github.com/pkg/errors"
)

// TopicID is the ID of a Topic.
type TopicID ID

// Scan implements sql.Scanner.
func (id *TopicID) Scan(value interface{}) error {
	return (*ID)(id).Scan(value)
}

// Value implements driver.Valuer.
func (id TopicID) Value() (driver.Value, error) {
	return ID(id).Value()
}

// TopicIDs a slice of Topic IDs.
type TopicIDs []TopicID

// Scan implements sql.Scanner.
func (ids *TopicIDs) Scan(value interface{}) error {
	var generic IDs
	if err := generic.Scan(value); err != nil {
		return err
	}

	*ids = make(TopicIDs, len(generic))
	for i, id := range generic {
		(*ids)[i] = TopicID(id)
	}

	return nil
}

// Value implements driver.Valuer.
func (ids TopicIDs) Value() (driver.Value, error) {
	generic := make(IDs, len(ids))
	for i, id := range ids {
		generic[i] = ID(id)
	}

	return generic.Value()
}

func scanTopicRows(rows *sql.Rows, err error) ([]*Topic, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var values []*Topic

	for rows.Next() {
		v, err := scanTopic(rows)
		if err != nil {
			return values, err
		}
		values = append(values, v)
	}

	return values, rows.Close()
}

// Get retrieves the Topic from the database.
func (id TopicID) Get(ctx *Context) (*Topic, error) {
	v, err := scanTopic(ctx.Tx.QueryRow(idGetTopic, ctx.CurrentUser, ctx.Sudo, id))
	return v, errors.Wrap(err, "get Topic by ID")
}

// Get retrieves each Topic from the database. The returned slice has the same
// order as the ids. If any Topic is missing from the database, a non-nil error
// will be returned.
func (ids TopicIDs) Get(ctx *Context) ([]*Topic, error) {
	values, err := scanTopicRows(ctx.Tx.Query(idsGetTopic, ctx.CurrentUser, ctx.Sudo, ids))

	sorted := make([]*Topic, len(ids))
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

	return sorted, errors.Wrap(err, "get Topic by ID")
}
