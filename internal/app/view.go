package app

import (
	"fmt"
	"strings"
)

func (m *app) View() string {
	switch m.state {
	case stateMain:
		return styleBox.Render(m.viewMenu(mainMenuItems, m.cursor, "*** Главное меню ***"))
	case stateImportSub:
		return styleBox.Render(m.viewMenu(importSubItems, m.subCursor, "--- Подменю: Импорт ---"))
	case stateExportSub:
		return styleBox.Render(m.viewMenu(exportSubItems, m.subCursor, "--- Подменю: Экспорт ---"))
	case stateBalanceSub:
		return styleBox.Render(m.viewMenu(balanceSubItems, m.subCursor, "--- Подменю: Пересчёт баланса ---"))
	case stateAnalyticsSub:
		return styleBox.Render(m.viewMenu(analyticsSubItems, m.subCursor, "--- Подменю: Аналитика ---"))

	case stateManualEntityMenu:
		return styleBox.Render(m.viewMenu(manualEntityItems, m.subCursor, "--- Ручной ввод: сущность ---"))
	case stateManualActionMenu:
		return styleBox.Render(m.viewMenu(manualActionItems, m.subCursor, "--- Ручной ввод: действие ---"))

	case stateManualFormInput:
		label := m.formFields[m.formStep].Label
		return styleBox.Render(
			stylePrompt.Render(label) + "\n\n" +
				styleInput.Render(m.textInput.View()),
		)

	case stateManualDeleteIDInput:
		return styleBox.Render(
			stylePrompt.Render("Введите ID для удаления:") + "\n\n" +
				styleInput.Render(m.textInput.View()),
		)

	case stateImportPathInput, stateExportNameInput:
		return styleBox.Render(
			stylePrompt.Render(m.renderPrompt()) + "\n\n" +
				styleInput.Render(m.textInput.View()),
		)

	case stateAnalyticsDateFromInput:
		return styleBox.Render(
			stylePrompt.Render("Введите дату начала (YYYY-MM-DD):") + "\n\n" +
				styleInput.Render(m.textInput.View()),
		)

	case stateAnalyticsDateToInput:
		return styleBox.Render(
			stylePrompt.Render("Введите дату конца (YYYY-MM-DD):") + "\n\n" +
				styleInput.Render(m.textInput.View()),
		)

	case stateManualAccountIDInput:
		return styleBox.Render(
			stylePrompt.Render("Введите ID счёта:") + "\n\n" +
				styleInput.Render(m.textInput.View()),
		)

	case stateManualNewBalanceInput:
		return styleBox.Render(
			stylePrompt.Render("Введите новый баланс (число):") + "\n\n" +
				styleInput.Render(m.textInput.View()),
		)

	case stateLogs:
		return styleBox.Render(m.viewLogs())

	case stateShowMessage:
		if m.isBusy {

			body := styleSpinner.Render(m.spin.View()+" "+m.busyMsg) + "\n\n" +
				styleHint.Render("(Операция выполняется… Esc — отмены нет, но можно выйти в меню)")
			return styleBox.Render(body)
		}
		return styleBox.Render(
			colorizeInfo(m.infoMessage) + "\n\n" +
				styleHint.Render("(Enter или Esc — назад)"),
		)

	case stateShowImportErrors:
		return styleBox.Render(m.viewImportErrors())
	default:
		return ""
	}
}

func (m *app) viewMenu(items []string, cursor int, title string) string {
	var b strings.Builder
	b.WriteString(styleTitle.Render(title))
	b.WriteString("\n")
	b.WriteString(styleHint.Render("(Навигация: ↑/↓, Enter — выбрать, Esc — назад)"))
	b.WriteString("\n\n")

	for i, it := range items {
		if i == cursor {
			b.WriteString(styleMenuCursor.Render("›"))
			b.WriteString(" ")
			b.WriteString(styleMenuItemSelected.Render(it))
			b.WriteString("\n")
		} else {
			b.WriteString("  ")
			b.WriteString(styleMenuItem.Render(it))
			b.WriteString("\n")
		}
	}
	return b.String()
}

func (m *app) viewLogs() string {
	logs := m.logStore.GetAll()
	var b strings.Builder
	b.WriteString(styleTitle.Render("--- Логи ---"))
	b.WriteString("\n\n")
	for _, l := range logs {
		b.WriteString(fmt.Sprintf("%v\n", l))
	}
	b.WriteString("\n")
	b.WriteString(styleHint.Render("(Enter или Esc — назад)"))
	return b.String()
}

func (m *app) viewImportErrors() string {
	if len(m.lastImportErrors) == 0 {
		return styleSuccess.Render("Ошибок импорта нет") + "\n\n" +
			styleHint.Render("(Enter или Esc — назад)")
	}
	var b strings.Builder
	b.WriteString(styleTitle.Render("--- Ошибки импорта ---"))
	b.WriteString("\n\n")
	for _, e := range m.lastImportErrors {
		b.WriteString(styleError.Render(e))
		b.WriteString("\n")
	}
	b.WriteString("\n")
	b.WriteString(styleHint.Render("(Enter или Esc — назад)"))
	return b.String()
}

func (m *app) renderPrompt() string {
	switch m.state {
	case stateImportPathInput:
		return "Введите путь к файлу для импорта:"
	case stateExportNameInput:
		return "Введите имя файла для экспорта:"
	}
	return ""
}

func colorizeInfo(s string) string {
	ls := strings.ToLower(s)
	switch {
	case strings.Contains(ls, "ошиб"):
		return styleError.Render(s)
	case strings.Contains(ls, "успеш") || strings.Contains(ls, "заверш"):
		return styleSuccess.Render(s)
	default:
		return stylePrompt.Render(s)
	}
}
