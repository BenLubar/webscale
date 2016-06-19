//go:generate go run mkid.go Permission

package model // import "github.com/BenLubar/webscale/model"

type Permission struct {
	ID   ID
	Slug string
}
