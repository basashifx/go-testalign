package testalign_test

import (
	"testing"

	testalign "github.com/basashifx/go-testalign"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()

	tests := []string{
		"basic",
		"noviolation",
		"subtests",
		"unexported",
		"mixed",
		"multifile",
		"pkgfuncs",
		"multi_receiver",
	}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			analysistest.Run(t, testdata, testalign.Analyzer, tt)
		})
	}
}

func TestAnalyzer_ExternalTestPackage(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, testalign.Analyzer, "externalapi")
}
