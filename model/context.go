package model // import "github.com/BenLubar/webscale/model"

import (
	"net/http"

	"github.com/BenLubar/webscale/db"
)

type Context struct {
	Tx          *db.Tx
	CurrentUser UserID
	Sudo        bool
	Request     *http.Request
	Page        int64
	Header      struct {
		Template   string
		Title      string
		Breadcrumb []Breadcrumb
	}
	Footer struct {
		CurrentPage int64
		PageCount   int64
	}
}

type Breadcrumb struct {
	Name string
	Path string
}
