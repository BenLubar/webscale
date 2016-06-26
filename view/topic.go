package view // import "github.com/BenLubar/webscale/view"

var Topic = parse("topic", `<section id="posts">
{{template "partial/posts-list" .Posts -}}
</section>
`)
