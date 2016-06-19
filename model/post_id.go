// This file is generated. DO NOT EDIT.

package model // import "github.com/BenLubar/webscale/model"

import "database/sql/driver"

// PostID is the ID of a Post.
type PostID ID

func (id *PostID) Scan(value interface{}) error {
	return (*ID)(id).Scan(value)
}
func (id PostID) Value() (driver.Value, error) {
	return ID(id).Value()
}

type PostIDs []PostID

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
func (ids PostIDs) Value() (driver.Value, error) {
	generic := make(IDs, len(ids))
	for i, id := range ids {
		generic[i] = ID(id)
	}

	return generic.Value()
}
