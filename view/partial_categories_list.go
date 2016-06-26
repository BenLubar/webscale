package view // import "github.com/BenLubar/webscale/view"

var _ = parse("partial/categories-list", `<ul>
{{range . -}}
<li data-category="{{.Category.ID}}">
{{template "partial/category-link" .Category}}
{{if .Topic -}}
(last post {{timestamp .Topic.BumpedAt}}{{if .PostAuthor}} by {{template "partial/user-link" .PostAuthor}}{{end}} in {{template "partial/topic-link" .Topic}})
{{end -}}
</li>
{{end -}}
</ul>
`)
