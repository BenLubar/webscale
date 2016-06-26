package view // import "github.com/BenLubar/webscale/view"

var Category = parse("category", `{{with .Children -}}
<section id="subcategories">
<h2>Subcategories</h2>
{{template "partial/categories-list" . -}}
</section>
{{end -}}

{{with .Topics -}}
<section id="topics">
<h2>Topics</h2>
{{template "partial/topics-list" . -}}
</section>
{{end -}}
`)
