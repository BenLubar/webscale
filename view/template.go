package view // import "github.com/BenLubar/webscale/view"

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/BenLubar/webscale/model"
	"github.com/BenLubar/webscale/static"
	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"
)

type Template struct {
	Name string
}

var tmpl = template.New("").Funcs(template.FuncMap{
	"static":    static.Path,
	"timestamp": timestamp,
	"pages":     pages,
})

func parse(name, content string) *Template {
	template.Must(tmpl.New(name).Parse(content))

	return &Template{Name: name}
}

func (t *Template) Execute(w http.ResponseWriter, ctx *model.Context, status int, data interface{}) error {
	ctx.Header.Template = t.Name
	ctx.Footer.CurrentPage = ctx.Page

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
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	w.WriteHeader(status)
	// Ignore errors generated while doing the actual sending of the
	// response. We've touched the ResponseWriter, so returning an error
	// would violate our contract.
	_, _ = io.Copy(w, &buf)

	return nil
}

func timestamp(t time.Time) template.HTML {
	human := humanize.Time(t)
	return template.HTML(`<time title="` + template.HTMLEscapeString(t.Format(time.RFC3339Nano)) + `" datetime="` + template.HTMLEscapeString(t.Format("2006-01-02T15:04:05.999Z07:00")) + `">` + template.HTMLEscapeString(human) + `</time>`)
}

func pages(current, count int64) template.HTML {
	if count <= 1 {
		return ""
	}

	// TODO: add quick links to the left and right of the current page and
	//       first/last links where applicable.
	b := []byte("<ul class=\"pages\">")
	b = append(b, "<li><form action=\"\" method=\"GET\"><input type=\"submit\" value=\"go to page\"><input type=\"number\" name=\"page\" min=\"0\" max=\""...)
	b = strconv.AppendInt(b, count-1, 10)
	b = append(b, "\" step=\"1\" value=\""...)
	b = strconv.AppendInt(b, current, 10)
	b = append(b, "\"></form></li>"...)
	b = append(b, "</ul>\n"...)

	return template.HTML(b)
}
