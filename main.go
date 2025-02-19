package main

import (
	"flag"
	"log"
)

func main() {
	// read configuration
	pathFlagPointer := flag.String("path", "", "Path to openvpn-status.log")

	flag.Parse()

	if *pathFlagPointer == "" {
		log.Fatal("No path provided")
	}
	err := runTelegraf(*pathFlagPointer)
	if err != nil {
		log.Fatal(err)
	}
}
