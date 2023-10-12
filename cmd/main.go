package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/taylormonacelli/goldbug"
	"github.com/taylormonacelli/lemondrop"
)

var verbose bool

func main() {
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output")
	flag.BoolVar(&verbose, "v", false, "Enable verbose output (shorthand)")
	flag.Parse()

	if verbose {
		goldbug.SetDefaultLoggerText(slog.LevelDebug)
	}

	lemondrop.WriteRegions(os.Stdout, verbose)
}
