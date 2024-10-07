package main

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

// Определяем новый анализатор NoOsExitAnalyzer
var NoOsExitAnalyzer = &analysis.Analyzer{
	Name: "noosexit",
	Doc:  "disallow direct use of os.Exit",
	Run:  run,
}

// Функция run будет вызываться для каждого файла, который анализируется
func run(pass *analysis.Pass) (interface{}, error) {
	// Проходим по всем файлам, которые анализируются
	for _, file := range pass.Files {
		// Функция Inspect позволяет обойти все узлы абстрактного синтаксического дерева (AST) файла
		ast.Inspect(file, func(node ast.Node) bool {
			// Проверяем, является ли текущий узел вызовом функции
			callExpr, ok := node.(*ast.CallExpr)
			if !ok {
				return false
			}

			// Проверяем, является ли вызов функции селектором (т.е. вызовом вида pkg.Func)
			selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
			if !ok {
				return false
			}

			// Проверяем, что функция вызывается из пакета os и называется Exit
			if ident, ok := selExpr.X.(*ast.Ident); ok && ident.Name == "os" && selExpr.Sel.Name == "Exit" {
				pass.Reportf(callExpr.Pos(), "direct use of os.Exit is not allowed")
			}
			return true
		})
	}

	return nil, nil
}
