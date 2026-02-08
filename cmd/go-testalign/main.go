package main

import (
	testalign "github.com/basashifx/go-testalign"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(testalign.Analyzer)
}
