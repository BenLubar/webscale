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

import "database/sql/driver"

// %[1]sID is the ID of a %[1]s.
type %[1]sID ID

func (id *%[1]sID) Scan(value interface{}) error {
	return (*ID)(id).Scan(value)
}
func (id %[1]sID) Value() (driver.Value, error) {
	return ID(id).Value()
}

type %[1]sIDs []%[1]sID

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
func (ids %[1]sIDs) Value() (driver.Value, error) {
	generic := make(IDs, len(ids))
	for i, id := range ids {
		generic[i] = ID(id)
	}

	return generic.Value()
}
`, os.Args[1])), 0664); err != nil {
		panic(err)
	}
}
