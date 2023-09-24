package main

import (
	"flag"
	"plex/src"
)

func main() {
	initFlag := flag.Bool("init", false, "Initialize Bifrost (configure media folders)")
	flag.Parse()

	if *initFlag {
		src.Init()
		return
	} else {
		src.StartServer()
	}
}
