package view // import "github.com/BenLubar/webscale/view"

var Index = parse("index", `{{with .Categories -}}
<section id="categories">
<h2>Categories</h2>
<ul>
{{range . -}}
<li><a href="/category/{{.ID}}/{{.Slug}}">{{.Name}}</a></li>
{{end -}}
</ul>
</section>
{{end -}}

{{with .Latest -}}
<section id="latest">
<h2>Latest</h2>
<ul>
{{range . -}}
<li><a href="/topic/{{.ID}}/{{.Slug}}">{{.Name}}</a></li>
{{end -}}
</ul>
</section>
{{end -}}
`)
