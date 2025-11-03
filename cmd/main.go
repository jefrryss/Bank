package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/dig"

	"task2/internal/analitic"
	"task2/internal/app"
	"task2/internal/bankmanager"
	"task2/internal/export"
	"task2/internal/handmaker"
	"task2/internal/importer"
	"task2/internal/logstorage"
)

func main() {
	container := dig.New()

	_ = container.Provide(logstorage.NewLogStorage)
	_ = container.Provide(bankmanager.NewBankManager)
	_ = container.Provide(importer.NewImportFacade)

	_ = container.Provide(export.NewExportFacade)

	// Аналитика: берём актуальные operation/catagory
	_ = container.Provide(func(bm *bankmanager.BankManager, ls *logstorage.LogStorage) *analitic.AnalyticsFacade {
		ops, _ := bm.GetAllOperations()
		cats, _ := bm.GetAllCategories()
		return analitic.NewAnalyticsFacade(ops, cats, ls)
	})

	_ = container.Provide(handmaker.NewController)

	_ = container.Provide(app.NewApp)

	err := container.Invoke(func(app tea.Model) error {
		p := tea.NewProgram(app)
		if err := p.Start(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		os.Exit(1)
	}
}
