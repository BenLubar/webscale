package view // import "github.com/BenLubar/webscale/view"

var _ = parse("partial/category-link", `<a href="/category/{{.ID}}/{{.Slug}}" data-category="{{.ID}}">{{.Name}}</a>`)
var _ = parse("partial/topic-link", `<a href="/topic/{{.ID}}/{{.Slug}}" data-topic="{{.ID}}">{{.Name}}</a>`)
var _ = parse("partial/user-link", `<a href="/user/{{.ID}}/{{.Slug}}" data-user="{{.ID}}">{{.Name}}</a>`)
