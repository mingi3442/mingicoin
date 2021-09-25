package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/mingi3442/mingicoin/explorer"
	"github.com/mingi3442/mingicoin/rest"
)

func usage() {
	fmt.Printf("Welcome to MINGI Coin\n\n")
	fmt.Printf("Please use the following flags:\n\n")
	fmt.Printf("-port:		Set the PORT of the server\n")
	fmt.Printf("-mode:		Choose between 'html' and 'rest'\n\n")
	runtime.Goexit() //모든 함수를 제거하지만 그 전에 defer를 먼저 이행한다
}

func Start() {
	if len(os.Args) == 1 {
		usage()
	}

	port := flag.Int("port", 4000, "Set port of the server")
	mode := flag.String("mode", "rest", "Choose between 'html' and 'rest'")

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
