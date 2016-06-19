// This file is generated. DO NOT EDIT.

package model // import "github.com/BenLubar/webscale/model"

import "database/sql/driver"

// CategoryID is the ID of a Category.
type CategoryID ID

func (id *CategoryID) Scan(value interface{}) error {
	return (*ID)(id).Scan(value)
}
func (id CategoryID) Value() (driver.Value, error) {
	return ID(id).Value()
}

type CategoryIDs []CategoryID

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
func (ids CategoryIDs) Value() (driver.Value, error) {
	generic := make(IDs, len(ids))
	for i, id := range ids {
		generic[i] = ID(id)
	}

	return generic.Value()
}
