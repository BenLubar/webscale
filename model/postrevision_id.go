// This file is generated. DO NOT EDIT.

package model // import "github.com/BenLubar/webscale/model"

import "database/sql/driver"

// PostRevisionID is the ID of a PostRevision.
type PostRevisionID ID

func (id *PostRevisionID) Scan(value interface{}) error {
	return (*ID)(id).Scan(value)
}
func (id PostRevisionID) Value() (driver.Value, error) {
	return ID(id).Value()
}

type PostRevisionIDs []PostRevisionID

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
func (ids PostRevisionIDs) Value() (driver.Value, error) {
	generic := make(IDs, len(ids))
	for i, id := range ids {
		generic[i] = ID(id)
	}

	return generic.Value()
}
