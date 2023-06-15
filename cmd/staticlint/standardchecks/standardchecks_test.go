package standardchecks

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestGetStandardRules(t *testing.T) {

	rules := GetStandardRules()
	for _, v := range rules {
		analysistest.Run(t, analysistest.TestData(), v, "./...")
	}
}
