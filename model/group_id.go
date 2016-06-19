// This file is generated. DO NOT EDIT.

package model // import "github.com/BenLubar/webscale/model"

import "database/sql/driver"

// GroupID is the ID of a Group.
type GroupID ID

func (id *GroupID) Scan(value interface{}) error {
	return (*ID)(id).Scan(value)
}
func (id GroupID) Value() (driver.Value, error) {
	return ID(id).Value()
}

type GroupIDs []GroupID

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
func (ids GroupIDs) Value() (driver.Value, error) {
	generic := make(IDs, len(ids))
	for i, id := range ids {
		generic[i] = ID(id)
	}

	return generic.Value()
}
