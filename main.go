package main

import (
	"github.com/ponder2000/rdpms25-template/pkg"
)

var (
	version   string = "unknown"
	buildTime string = "unknown"
)

func main() {
	pkg.Start(version, buildTime)
}
