package view // import "github.com/BenLubar/webscale/view"

import (
	"bytes"
	"html/template"
	"io"
	"net/http"

	"github.com/BenLubar/webscale/model"
	"github.com/pkg/errors"
)

type Template struct {
	Name string
}

var tmpl = template.New("")

func parse(name, content string) *Template {
	template.Must(tmpl.New(name).Parse(content))

	return &Template{Name: name}
}

func (t *Template) Execute(w http.ResponseWriter, ctx *model.Context, status int, data interface{}) error {
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "header", ctx.Header); err != nil {
		return errors.Wrap(err, "rendering header")
	}
	if err := tmpl.ExecuteTemplate(&buf, t.Name, data); err != nil {
		return errors.Wrapf(err, "rendering template %q", t.Name)
	}
	if err := tmpl.ExecuteTemplate(&buf, "footer", ctx.Footer); err != nil {
		return errors.Wrap(err, "rendering footer")
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	// Ignore errors generated while doing the actual sending of the
	// response. We've touched the ResponseWriter, so returning an error
	// would violate our contract.
	_, _ = io.Copy(w, &buf)

	return nil
}
