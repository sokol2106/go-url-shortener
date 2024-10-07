package main

import (
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
	// Определяем список анализаторов, которые будем использовать
	mychecks := []*analysis.Analyzer{
		NoOsExitAnalyzer,   // Собственный анализатор, запрещающий os.Exit
		printf.Analyzer,    // Анализатор проверки правильности использования форматирования в printf
		shadow.Analyzer,    // Анализатор, предупреждающий о затенении переменных
		structtag.Analyzer, // Анализатор для проверки правильности тегов структур
		errcheck.Analyzer,  // Анализатор, проверяющий необработанные ошибки
		shift.Analyzer,     // Анализатор проверки правильности сдвигов битовых операций
	}

	// Проходим по всем анализаторам из staticcheck и добавляем нужные анализаторы
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

	// Запускаем мультианализатор с собранным списком проверок
	multichecker.Main(
		mychecks...,
	)
}
