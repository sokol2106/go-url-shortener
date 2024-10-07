package main

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

var NoOsExitAnalyzer = &analysis.Analyzer{
	Name: "noosexit",
	Doc:  "disallow direct use of os.Exit",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			callExpr, ok := node.(*ast.CallExpr)
			if !ok {
				return false
			}

			selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
			if !ok {
				return false
			}

			if ident, ok := selExpr.X.(*ast.Ident); ok && ident.Name == "os" && selExpr.Sel.Name == "Exit" {
				pass.Reportf(callExpr.Pos(), "direct use of os.Exit is not allowed")
			}
			return true
		})
	}

	return nil, nil
}
