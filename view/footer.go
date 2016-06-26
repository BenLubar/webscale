package view // import "github.com/BenLubar/webscale/view"

var _ = parse("footer", `{{pages .CurrentPage .PageCount -}}
</main>
<footer>
</footer>
</body>
</html>
`)
