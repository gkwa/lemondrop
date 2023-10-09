package main

import (
	"flag"
	"os"

	"github.com/taylormonacelli/lemondrop"
)

var verbose bool

func main() {
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output")
	flag.BoolVar(&verbose, "v", false, "Enable verbose output (shorthand)")
	flag.Parse()

	lemondrop.WriteRegions(os.Stdout, verbose)
}
