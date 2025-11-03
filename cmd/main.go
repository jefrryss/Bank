package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jefrryss/Bank/internal/analitic"
	"github.com/jefrryss/Bank/internal/app"
	"github.com/jefrryss/Bank/internal/bankmanager"
	"github.com/jefrryss/Bank/internal/export"
	"github.com/jefrryss/Bank/internal/handmaker"
	"github.com/jefrryss/Bank/internal/importer"
	"github.com/jefrryss/Bank/internal/logstorage"
	"go.uber.org/dig"
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
