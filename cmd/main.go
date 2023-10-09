package main

import (
	"os"

	"github.com/taylormonacelli/lemondrop"
)

func main() {
	lemondrop.GetRegions(os.Stdout)
}
