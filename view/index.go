package view // import "github.com/BenLubar/webscale/view"

var Index = parse("index", `<section id="categories">
<h2>Categories</h2>
{{with .Categories -}}
<ul>
{{range . -}}
<li><a href="/category/{{.ID}}/{{.Slug}}">{{.Name}}</a></li>
{{end -}}
</ul>
{{else -}}
<p>No categories found</p>
{{end -}}
</section>
`)
