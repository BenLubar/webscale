// This file is generated. DO NOT EDIT.

package model // import "github.com/BenLubar/webscale/model"

import "database/sql/driver"

// PermissionID is the ID of a Permission.
type PermissionID ID

func (id *PermissionID) Scan(value interface{}) error {
	return (*ID)(id).Scan(value)
}
func (id PermissionID) Value() (driver.Value, error) {
	return ID(id).Value()
}

type PermissionIDs []PermissionID

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
func (ids PermissionIDs) Value() (driver.Value, error) {
	generic := make(IDs, len(ids))
	for i, id := range ids {
		generic[i] = ID(id)
	}

	return generic.Value()
}
