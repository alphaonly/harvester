//Package staticchecks describes which analyzers of staticcheck.io are to be involved in checks array
package staticchecks

import (
	"strings"

	"golang.org/x/tools/go/analysis"
	"honnef.co/go/tools/staticcheck"
)

// GetSARules - collect a slice of staticcheck analyzers that start of "SA" prefix
func GetSARules() []*analysis.Analyzer {
	var saRules []*analysis.Analyzer
	for _, analyzer := range staticcheck.Analyzers {
		if strings.HasPrefix(analyzer.Analyzer.Name, "SA") {
			saRules = append(saRules, analyzer.Analyzer)
		}
	}
	return saRules
}
