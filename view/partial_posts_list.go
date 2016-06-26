package view // import "github.com/BenLubar/webscale/view"

var _ = parse("partial/posts-list", `<ul>
{{range . -}}
<li data-post="{{.Post.ID}}">
{{template "partial/post" . -}}
</li>
{{end -}}
</ul>
`)
