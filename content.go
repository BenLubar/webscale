package main // import "github.com/BenLubar/webscale"

import (
	_ "net/http/pprof"

	_ "github.com/BenLubar/webscale/controller"
	_ "github.com/BenLubar/webscale/static"
)
