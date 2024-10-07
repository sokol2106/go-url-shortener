package main

import (
	"fmt"
	"github.com/kisielk/errcheck/errcheck"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	mychecks := []*analysis.Analyzer{
		NoOsExitAnalyzer,
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		errcheck.Analyzer,
		shift.Analyzer,
	}

	for _, v := range staticcheck.Analyzers {
		if v.Analyzer != nil {
			fmt.Println(v.Analyzer.Name)
		}
	}

	for _, v := range staticcheck.Analyzers {
		if v.Analyzer != nil && v.Analyzer.Name != "" {
			if v.Analyzer.Name[:2] == "SA" {
				mychecks = append(mychecks, v.Analyzer)
			}

			if v.Analyzer.Name == "S1005" {
				mychecks = append(mychecks, v.Analyzer)
			}

			if v.Analyzer.Name == "ST1005" {
				mychecks = append(mychecks, v.Analyzer)
			}

			if v.Analyzer.Name == "QF1003" {
				mychecks = append(mychecks, v.Analyzer)
			}
		}
	}

	multichecker.Main(
		mychecks...,
	)
}
