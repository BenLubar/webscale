package view // import "github.com/BenLubar/webscale/view"

var _ = parse("header", `<!DOCTYPE html>
<html class="tmpl-{{.Template}}">
<head>
<meta charset="utf-8">
<title>{{with .Title}}{{.}} - {{end}}#webscale</title>
<meta name="viewport" content="width=device-width, initial-scale=1">
<link rel="stylesheet" href="/static/style.css">
</head>
<body>
<header>
<h1>{{with .Title}}{{.}}{{else}}#webscale{{end}}</h1>
{{with .Breadcrumb -}}
<ol class="breadcrumb" itemscope itemtype="http://schema.org/BreadcrumbList">
{{range . -}}
<li itemprop="itemListElement" itemscope itemtype="http://schema.org/ListItem"><a itemprop="item" itemscope itemtype="http://schema.org/Thing" href="{{.Path}}"><span itemprop="name">{{.Name}}</span></a></li>
{{end -}}
</ol>
{{end -}}
</header>
<main>
`)
