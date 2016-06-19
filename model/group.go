//go:generate go run mkid.go Group

package model // import "github.com/BenLubar/webscale/model"

type Group struct {
	ID   GroupID
	Name string
	Slug string
}
