package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"task2/internal/analitic"
	"task2/internal/bankmanager"
	"task2/internal/client"
	"task2/internal/domain/export"
	"task2/internal/domain/importer"
)

func main() {

	// ===== Создаём BankManager =====
	manager := bankmanager.NewBankManager()
	accounts, err := manager.GetAllAccounts()
	if err != nil {
		fmt.Println("Ошибка получения аккаунтов:", err)
		return
	}
	categories, err := manager.GetAllCategories()
	if err != nil {
		fmt.Println("Ошибка получения категорий:", err)
		return
	}
	operations, err := manager.GetAllOperations()
	if err != nil {
		fmt.Println("Ошибка получения операций:", err)
		return
	}
	accountsPtr := &accounts
	categoriesPtr := &categories
	operationsPtr := &operations

	importJSON := importer.NewJSONParser()
	importerFacade := importer.NewImportFacade()

	exportJSON := export.NewExportJSON()
	exporterFacade := export.NewExportFacade()
	exporterFacade.Init(categoriesPtr, operationsPtr, accountsPtr)

	analyticsFacade := analitic.NewAnalyticsFacade(operations, categories)

	// ===== Создаём модель для Bubble Tea =====
	m := client.NewModel(manager)

	// ===== Переопределяем performAction для интеграции фасадов =====
	m.PerformAction = func(activeID string) string {
		switch activeID {
		case "import":
			if err := importerFacade.Init(importJSON); err != nil {
				return fmt.Sprintf("Ошибка импорта: %v", err)
			}
			return "Импорт данных завершён"
		case "export":
			exporterFacade.StartExport(exportJSON)
			return "Экспорт завершён"
		case "analytics_balance":
			analyticsFacade.CalcBalance()
			return "Расчёт баланса выполнен"
		case "analytics_category":
			analyticsFacade.GroupByCategory()
			return "Группировка по категориям выполнена"
		case "view_errors":
			errors := manager.GetErrors()
			if len(errors) == 0 {
				return "Ошибок нет"
			}
			var lines []string
			for _, e := range errors {
				lines = append(lines, fmt.Sprintf("%v : %v", e.Line, e.Err))
			}
			return strings.Join(lines, "\n")
		default:
			return ""
		}
	}

	// ===== Запуск TUI =====
	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Printf("Ошибка запуска интерфейса: %v\n", err)
		os.Exit(1)
	}
}
