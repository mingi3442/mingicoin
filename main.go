package main

import (
	"github.com/mingi3442/mingicoin/cli"
	"github.com/mingi3442/mingicoin/db"
)

func main() {
	defer db.Close()
	cli.Start()
}
