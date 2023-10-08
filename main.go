package main

import (
	"bifrost/src"
	"flag"
)

func main() {
	initFlag := flag.Bool("init", false, "Initialize Bifrost (configure target folders)")
	flag.Parse()

	if *initFlag {
		src.Init()
		return
	} else {
		src.StartServer()
	}
}
