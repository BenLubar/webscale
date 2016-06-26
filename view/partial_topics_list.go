package view // import "github.com/BenLubar/webscale/view"

var _ = parse("partial/topics-list", `<ul>
{{range . -}}
<li data-topic="{{.Topic.ID}}">
{{template "partial/topic-link" .Topic}}
{{if .Post -}}
(last post {{timestamp .Post.CreatedAt}}{{if .PostAuthor}} by {{template "partial/user-link" .PostAuthor}}{{end}})
{{end -}}
</li>
{{end -}}
</ul>
`)
