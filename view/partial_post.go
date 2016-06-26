package view // import "github.com/BenLubar/webscale/view"

var _ = parse("partial/post", `{{with .PostAuthor -}}
<aside class="post-author" data-user="{{.ID}}">
<strong class="user-name">{{template "partial/user-link" .}}</strong>
{{if .Avatar -}}
<br>
<img src="{{.Avatar}}" alt="" title="{{.Name}}'s avatar">
{{end -}}
{{if .Location -}}
<br>
<span class="user-location">{{.Location}}</span>
{{end -}}
</aside>
{{end -}}
<article class="post-content">
{{.Post.Content}}
{{with .PostAuthor}}{{if .Signature -}}
<aside class="user-signature">
{{.Signature}}
</aside>
{{end}}{{end -}}
</article>
`)
