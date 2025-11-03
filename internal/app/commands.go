package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"task2/domain/ports"
	"task2/internal/balancemanager"
	"task2/internal/export"
	"task2/internal/importer"
)

type balanceDoneMsg struct{ msg string }
type exportDoneMsg struct{ msg string }
type importDoneMsg struct {
	msg  string
	errs []string
}

func (m *app) cmdAutoBalance() tea.Cmd {
	return func() tea.Msg {
		auto := balancemanager.NewAutoBalanceManager(m.bankManager)
		if err := auto.Recalculate(); err != nil {
			return balanceDoneMsg{msg: fmt.Sprintf("Ошибка при пересчёте: %v", err)}
		}
		return balanceDoneMsg{msg: "Баланс автоматически пересчитан по операциям."}
	}
}

func (m *app) cmdManualBalance(accountID string, newBalance float64) tea.Cmd {
	return func() tea.Msg {
		man := balancemanager.NewManualBalanceManager(m.bankManager, accountID, newBalance)
		if err := man.Recalculate(); err != nil {
			return balanceDoneMsg{msg: fmt.Sprintf("Ошибка при обновлении баланса: %v", err)}
		}
		return balanceDoneMsg{msg: fmt.Sprintf("Баланс счёта %s установлен на %.2f", accountID, newBalance)}
	}
}

func (m *app) cmdExport(name string, isCSV bool) tea.Cmd {
	return func() tea.Msg {
		ops, _ := m.bankManager.GetAllOperations()
		cats, _ := m.bankManager.GetAllCategories()
		accs, _ := m.bankManager.GetAllAccounts()

		m.exportFacade.Init(cats, ops, accs, m.logStore)

		var visitor export.ExporterVisitor
		if isCSV {
			visitor = export.NewExportCSV()
		} else {
			visitor = export.NewExportJSON()
		}
		if err := visitor.SetFilePath(name); err != nil {
			return exportDoneMsg{msg: fmt.Sprintf("Ошибка с файлом экспорта: %v", err)}
		}
		if err := m.exportFacade.StartExport(visitor); err != nil {
			return exportDoneMsg{msg: fmt.Sprintf("Ошибка экспорта: %v", err)}
		}
		return exportDoneMsg{msg: fmt.Sprintf("Экспорт успешно завершён: %s", name)}
	}
}

func (m *app) cmdImport(path string, isCSV bool) tea.Cmd {
	return func() tea.Msg {
		var di ports.DataImporter
		if isCSV {
			di = importer.NewCSVParser()
		} else {
			di = importer.NewJSONParser()
		}

		if err := di.SetFilePath(path); err != nil {
			return importDoneMsg{msg: fmt.Sprintf("Ошибка с файлом: %v", err)}
		}

		fac := importer.NewImportFacade()
		if err := fac.Init(di, m.logStore); err != nil {
			return importDoneMsg{msg: fmt.Sprintf("Ошибка импорта: %v", err)}
		}

		m.bankManager.AddAccounts(convertToPointers(fac.GetAccounts()))
		m.bankManager.AddCategories(convertToPointers(fac.GetCategory()))
		m.bankManager.AddOperations(convertToPointers(fac.GetOperation()))

		errs := make([]string, 0, len(fac.GetErrors()))
		for _, e := range fac.GetErrors() {
			errs = append(errs, fmt.Sprintf("%v: %v", e.Line, e.Err))
		}

		msg := fmt.Sprintf("Импорт завершён: %d счетов, %d категорий, %d операций",
			len(fac.GetAccounts()), len(fac.GetCategory()), len(fac.GetOperation()))
		if len(errs) > 0 {
			msg += fmt.Sprintf(" (с ошибками: %d)", len(errs))
		}

		return importDoneMsg{msg: msg, errs: errs}
	}
}

func convertToPointers[T any](slice []T) []*T {
	ptrs := make([]*T, len(slice))
	for i := range slice {
		ptrs[i] = &slice[i]
	}
	return ptrs
}
