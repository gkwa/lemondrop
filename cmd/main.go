package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/taylormonacelli/goldbug"
	"github.com/taylormonacelli/lemondrop"
)

var (
	verbose     bool
	veryVerbose bool
)

func main() {
	flag.BoolVar(&verbose, "v", false, "Show user friendly region name (shorthand)")
	flag.BoolVar(&verbose, "verbose", false, "Show user friendly region name")

	flag.BoolVar(&veryVerbose, "very-verbose", false, "Show debug")

	flag.Parse()

	if veryVerbose {
		goldbug.SetDefaultLoggerText(slog.LevelDebug)
	}

	lemondrop.WriteRegions(os.Stdout, verbose)
}
