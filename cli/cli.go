package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/ParkChangSun/nomadcoin/explorer"
	"github.com/ParkChangSun/nomadcoin/rest"
)

func usage() {
	fmt.Println("-port:set port")
	fmt.Println("-mode:select mode")
	runtime.Goexit()
}

func Start() {

	if len(os.Args) == 1 {
		usage()
	}

	port := flag.Int("port", 4000, "set port")
	mode := flag.String("mode", "rest", "select mode")

	flag.Parse()

	switch *mode {
	case "rest":
		rest.Start(*port)
	case "html":
		explorer.Start(*port)
	default:
		usage()
	}
}
