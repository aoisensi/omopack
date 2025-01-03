package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Parse()
	source := flag.Arg(0)
	if source == "" {
		source = "."
	}
	p := new(Packer)
	p.Source = source
	if err := p.openManifest(); err != nil {
		fmt.Println("Press Enter to exit")
		fmt.Scanln()
		os.Exit(1)
	}
	p.Pack()
	p.Dest()
	fmt.Println("Press Enter to exit")
	fmt.Scanln()
}
