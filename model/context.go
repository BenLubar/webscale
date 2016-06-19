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
	Header      struct {
		Title      string
		Breadcrumb []Breadcrumb
	}
	Footer struct {
	}
}

type Breadcrumb struct {
	Name string
	Path string
}
