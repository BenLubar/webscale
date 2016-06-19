package view // import "github.com/BenLubar/webscale/view"

var Category = parse("category", `{{with .Children -}}
<section id="subcategories">
<h2>Subcategories</h2>
<ul>
{{range . -}}
<li><a href="/category/{{.ID}}/{{.Slug}}">{{.Name}}</a></li>
{{end -}}
</ul>
</section>
{{end -}}
`)
