// Package staticlint contains main logic and main function for the analyzer program.
// It includes all standard rules, every staticcheck.io analyzer with "SA" prefix
// and analyzer that checks the use of os.Exit in main func
package main

import (
	"github.com/alphaonly/harvester/cmd/staticlint/mainanalyzer"
	"github.com/alphaonly/harvester/cmd/staticlint/standardchecks"
	"github.com/alphaonly/harvester/cmd/staticlint/staticchecks"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {

	//Add SA staticcheck analyzers
	rules := staticchecks.GetSARules()
	//Add standard analyzers
	rules = append(rules, standardchecks.GetStandardRules()...)
	//Add "os.Exit call in main" check analyzer
	rules = append(rules, mainanalyzer.MainOsExit)

	//Add the slice to multichecker
	multichecker.Main(
		rules...,
	)
}
