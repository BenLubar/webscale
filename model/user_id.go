// This file is generated. DO NOT EDIT.

package model // import "github.com/BenLubar/webscale/model"

import "database/sql/driver"

// UserID is the ID of a User.
type UserID ID

func (id *UserID) Scan(value interface{}) error {
	return (*ID)(id).Scan(value)
}
func (id UserID) Value() (driver.Value, error) {
	return ID(id).Value()
}

type UserIDs []UserID

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
func (ids UserIDs) Value() (driver.Value, error) {
	generic := make(IDs, len(ids))
	for i, id := range ids {
		generic[i] = ID(id)
	}

	return generic.Value()
}
