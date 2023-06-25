package staticchecks

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestGetSARules(t *testing.T){

	rules := GetSARules()
	for _, v := range rules {
		analysistest.Run(t, analysistest.TestData(), v, "./...")
	}
}
