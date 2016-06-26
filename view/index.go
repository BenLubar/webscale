package view // import "github.com/BenLubar/webscale/view"

var Index = parse("index", `{{with .Categories -}}
<section id="categories">
<h2>Categories</h2>
{{template "partial/categories-list" . -}}
</section>
{{end -}}

{{with .Latest -}}
<section id="latest">
<h2>Latest</h2>
{{template "partial/topics-list" . -}}
</section>
{{end -}}
`)
