// This file is generated. DO NOT EDIT.

package model // import "github.com/BenLubar/webscale/model"

import "database/sql/driver"

// TopicID is the ID of a Topic.
type TopicID ID

func (id *TopicID) Scan(value interface{}) error {
	return (*ID)(id).Scan(value)
}
func (id TopicID) Value() (driver.Value, error) {
	return ID(id).Value()
}

type TopicIDs []TopicID

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
func (ids TopicIDs) Value() (driver.Value, error) {
	generic := make(IDs, len(ids))
	for i, id := range ids {
		generic[i] = ID(id)
	}

	return generic.Value()
}
