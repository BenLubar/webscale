//go:generate -command asset go run asset.go -var=_
//go:generate npm install --silent --no-progress --production less less-plugin-clean-css
//go:generate node_modules/.bin/lessc --clean-css=advanced --strict-math=on --strict-units=on style.less style.css
//go:generate asset style.css

package static // import "github.com/BenLubar/webscale/static"

import (
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"strings"
)

func css(a asset) asset { return doAsset(a) }

var tags = make(map[string]string)

func doAsset(a asset) asset {
	http.Handle("/static/"+a.Name, a)
	if a.etag != "" {
		tagBytes, err := base64.StdEncoding.DecodeString(strings.Trim(a.etag, `"`))
		if err != nil {
			tagBytes = []byte(a.etag)
		}
		tag := hex.EncodeToString(tagBytes)
		http.Handle("/static/"+tag+"/"+a.Name, &cacheForever{a})
		tags[a.Name] = tag
	}
	return a
}

// Path returns the path the named resource.
func Path(name string) string {
	if tag, ok := tags[name]; ok {
		return "/static/" + tag + "/" + name
	}
	return "/static/" + name
}

type cacheForever struct {
	h http.Handler
}

func (c *cacheForever) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// from nginx's "expires max" directive
	w.Header().Set("Expires", "Thu, 31 Dec 2037 23:55:55 GMT")
	w.Header().Set("Cache-Control", "max-age=315360000")
	c.h.ServeHTTP(w, r)
}
