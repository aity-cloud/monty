package main

import (
	"github.com/aity-cloud/monty/internal/linter"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(linter.AnalyzerPlugin.GetAnalyzers()...)
}
