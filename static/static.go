//go:generate -command asset go run asset.go
//go:generate npm install --silent --no-progress --production less less-plugin-clean-css
//go:generate node_modules/.bin/lessc --clean-css=advanced --strict-math=on --strict-units=on style.less style.css
//go:generate asset style.css

package static // import "github.com/BenLubar/webscale/static"

import "net/http"

func css(a asset) asset {
	http.Handle("/static/"+a.Name, a)
	return a
}
