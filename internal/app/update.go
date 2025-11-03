package app

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jefrryss/Bank/internal/analitic"
)

func (m *app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spin, cmd = m.spin.Update(msg)
		if m.isBusy {
			return m, cmd
		}
		return m, nil

	case balanceDoneMsg:
		m.isBusy = false
		m.infoMessage = msg.msg
		return m, nil

	case exportDoneMsg:
		m.isBusy = false
		m.infoMessage = msg.msg
		return m, nil

	case importDoneMsg:
		m.isBusy = false
		m.lastImportErrors = msg.errs
		m.infoMessage = msg.msg
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {

		case stateMain:
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(mainMenuItems)-1 {
					m.cursor++
				}
			case "enter":
				switch mainMenuItems[m.cursor] {
				case "Импорт":
					m.state, m.subCursor = stateImportSub, 0
				case "Экспорт":
					m.state, m.subCursor = stateExportSub, 0
				case "Пересчёт баланса":
					m.state, m.subCursor = stateBalanceSub, 0
				case "Аналитика":
					m.state, m.subCursor = stateAnalyticsSub, 0
				case "Ручной ввод":
					m.state, m.subCursor = stateManualEntityMenu, 0
				case "Просмотр логов":
					m.state = stateLogs
				case "Выход":
					return m, tea.Quit
				}
			case "esc", "q":
				return m, tea.Quit
			}

		case stateImportSub:
			switch msg.String() {
			case "up", "k":
				if m.subCursor > 0 {
					m.subCursor--
				}
			case "down", "j":
				if m.subCursor < len(importSubItems)-1 {
					m.subCursor++
				}
			case "enter":
				switch importSubItems[m.subCursor] {
				case "Импорт CSV":
					m.state = stateImportPathInput
					m.textInput.Placeholder = "Введите путь к CSV файлу"
					m.textInput.SetValue("")
				case "Импорт JSON":
					m.state = stateImportPathInput
					m.textInput.Placeholder = "Введите путь к JSON файлу"
					m.textInput.SetValue("")
				case "Просмотр ошибок при импорте":
					m.state = stateShowImportErrors
				case "Назад":
					m.state = stateMain
				}
			case "esc":
				m.state = stateMain
			}

		case stateExportSub:
			switch msg.String() {
			case "up", "k":
				if m.subCursor > 0 {
					m.subCursor--
				}
			case "down", "j":
				if m.subCursor < len(exportSubItems)-1 {
					m.subCursor++
				}
			case "enter":
				switch exportSubItems[m.subCursor] {
				case "Экспорт CSV":
					m.state = stateExportNameInput
					m.textInput.Placeholder = "Введите имя CSV файла"
					m.textInput.SetValue("")
				case "Экспорт JSON":
					m.state = stateExportNameInput
					m.textInput.Placeholder = "Введите имя JSON файла"
					m.textInput.SetValue("")
				case "Назад":
					m.state = stateMain
				}
			case "esc":
				m.state = stateMain
			}

		case stateBalanceSub:
			switch msg.String() {
			case "up", "k":
				if m.subCursor > 0 {
					m.subCursor--
				}
			case "down", "j":
				if m.subCursor < len(balanceSubItems)-1 {
					m.subCursor++
				}
			case "enter":
				switch balanceSubItems[m.subCursor] {
				case "Автоматически":
					m.isBusy = true
					m.busyMsg = "Пересчитываю баланс…"
					m.state = stateShowMessage
					return m, tea.Batch(m.spin.Tick, m.cmdAutoBalance())

				case "Вручную":
					m.textInput.Placeholder = "Введите ID счёта"
					m.textInput.SetValue("")
					m.state = stateManualAccountIDInput

				case "Назад":
					m.state = stateMain
				}
			case "esc":
				m.state = stateMain
			}

		case stateManualEntityMenu:
			switch msg.String() {
			case "up", "k":
				if m.subCursor > 0 {
					m.subCursor--
				}
			case "down", "j":
				if m.subCursor < len(manualEntityItems)-1 {
					m.subCursor++
				}
			case "enter":
				switch manualEntityItems[m.subCursor] {
				case "Счета":
					m.manualEntity = "account"
					m.state, m.subCursor = stateManualActionMenu, 0
				case "Категории":
					m.manualEntity = "category"
					m.state, m.subCursor = stateManualActionMenu, 0
				case "Операции":
					m.manualEntity = "operation"
					m.state, m.subCursor = stateManualActionMenu, 0
				case "Назад":
					m.state = stateMain
				}
			case "esc":
				m.state = stateMain
			}

		case stateManualActionMenu:
			switch msg.String() {
			case "up", "k":
				if m.subCursor > 0 {
					m.subCursor--
				}
			case "down", "j":
				if m.subCursor < len(manualActionItems)-1 {
					m.subCursor++
				}
			case "enter":
				switch manualActionItems[m.subCursor] {
				case "Создать":
					m.manualAction = "create"
					m.prepareForm()
					m.state = stateManualFormInput
				case "Редактировать":
					m.manualAction = "update"
					m.prepareForm()
					m.state = stateManualFormInput
				case "Удалить":
					m.manualAction = "delete"
					m.textInput.Placeholder = "Введите ID для удаления"
					m.textInput.SetValue("")
					m.state = stateManualDeleteIDInput
				case "Назад":
					m.state = stateManualEntityMenu
					m.subCursor = 0
				}
			case "esc":
				m.state = stateManualEntityMenu
				m.subCursor = 0
			}

		case stateManualFormInput:
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			if msg.String() == "enter" {
				key := m.formFields[m.formStep].Key
				m.formValues[key] = strings.TrimSpace(m.textInput.Value())
				m.formStep++
				if m.formStep >= len(m.formFields) {
					if err := m.submitForm(); err != nil {
						m.infoMessage = fmt.Sprintf("Ошибка: %v", err)
					} else {
						switch m.manualAction {
						case "create":
							m.infoMessage = "Создание успешно выполнено."
						case "update":
							m.infoMessage = "Обновление успешно выполнено."
						}
					}
					m.state = stateShowMessage
				} else {
					next := m.formFields[m.formStep]
					m.textInput.Placeholder = next.Label
					m.textInput.SetValue("")
				}
			}
			return m, cmd

		case stateManualDeleteIDInput:
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			if msg.String() == "enter" {
				id := strings.TrimSpace(m.textInput.Value())
				var err error
				switch m.manualEntity {
				case "account":
					err = m.ctrl.DeleteAccount(id)
				case "category":
					err = m.ctrl.DeleteCategory(id)
				case "operation":
					err = m.ctrl.DeleteOperation(id)
				}
				if err != nil {
					m.infoMessage = fmt.Sprintf("Ошибка удаления: %v", err)
				} else {
					m.infoMessage = "Удаление успешно выполнено."
				}
				m.state = stateShowMessage
			}
			return m, cmd

		case stateManualAccountIDInput:
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			if msg.String() == "enter" {
				m.manualAccountID = strings.TrimSpace(m.textInput.Value())
				if m.manualAccountID == "" {
					m.infoMessage = "ID счёта не может быть пустым."
					m.state = stateShowMessage
					return m, cmd
				}
				m.textInput.Placeholder = "Введите новый баланс (число)"
				m.textInput.SetValue("")
				m.state = stateManualNewBalanceInput
			}
			return m, cmd

		case stateManualNewBalanceInput:
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			if msg.String() == "enter" {
				raw := strings.TrimSpace(m.textInput.Value())
				raw = strings.ReplaceAll(raw, ",", ".")
				val, err := strconv.ParseFloat(raw, 64)
				if err != nil {
					m.infoMessage = "Неверное число. Введите, например: 1234.56"
					m.state = stateShowMessage
					return m, cmd
				}
				m.manualBalance = val

				m.isBusy = true
				m.busyMsg = "Обновляю баланс…"
				m.state = stateShowMessage
				return m, tea.Batch(m.spin.Tick, m.cmdManualBalance(m.manualAccountID, m.manualBalance))
			}
			return m, cmd

		case stateAnalyticsSub:
			switch msg.String() {
			case "up", "k":
				if m.subCursor > 0 {
					m.subCursor--
				}
			case "down", "j":
				if m.subCursor < len(analyticsSubItems)-1 {
					m.subCursor++
				}
			case "enter":
				switch analyticsSubItems[m.subCursor] {
				case "Подсчет разницы доходов и расходов за период":
					m.textInput.Placeholder = "Введите дату начала (YYYY-MM-DD)"
					m.textInput.SetValue("")
					m.state = stateAnalyticsDateFromInput
				case "Группировка по категориям":
					ops, _ := m.bankManager.GetAllOperations()
					cats, _ := m.bankManager.GetAllCategories()
					m.analyticsFacade = analitic.NewAnalyticsFacade(ops, cats, m.logStore)
					m.infoMessage = m.analyticsFacade.GroupByCategory()
					m.state = stateShowMessage
				case "Назад":
					m.state = stateMain
				}
			case "esc":
				m.state = stateMain
			}

		case stateAnalyticsDateFromInput:
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			if msg.String() == "enter" {
				m.analyticsFrom = strings.TrimSpace(m.textInput.Value())
				m.textInput.Placeholder = "Введите дату конца (YYYY-MM-DD)"
				m.textInput.SetValue("")
				m.state = stateAnalyticsDateToInput
			}
			return m, cmd

		case stateAnalyticsDateToInput:
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			if msg.String() == "enter" {
				m.analyticsTo = strings.TrimSpace(m.textInput.Value())
				from, err1 := time.Parse("2006-01-02", m.analyticsFrom)
				to, err2 := time.Parse("2006-01-02", m.analyticsTo)
				if err1 != nil || err2 != nil {
					m.infoMessage = "Неверный формат даты. Используйте YYYY-MM-DD"
				} else {
					ops, _ := m.bankManager.GetAllOperations()
					cats, _ := m.bankManager.GetAllCategories()
					m.analyticsFacade = analitic.NewAnalyticsFacade(ops, cats, m.logStore)
					m.infoMessage = m.analyticsFacade.CalcBalance(from, to)
				}
				m.state = stateShowMessage
			}
			return m, cmd

		case stateImportPathInput:
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			if msg.String() == "enter" {
				path := strings.TrimSpace(m.textInput.Value())
				isCSV := strings.Contains(m.textInput.Placeholder, "CSV")

				m.lastImportErrors = nil
				m.isBusy = true
				m.busyMsg = "Импортирую данные…"
				m.state = stateShowMessage
				return m, tea.Batch(m.spin.Tick, m.cmdImport(path, isCSV))
			}
			return m, cmd

		case stateExportNameInput:
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			if msg.String() == "enter" {
				name := strings.TrimSpace(m.textInput.Value())
				if name == "" {
					m.infoMessage = "Имя файла не может быть пустым"
					m.state = stateShowMessage
					return m, cmd
				}
				isCSV := strings.Contains(m.textInput.Placeholder, "CSV")

				m.isBusy = true
				m.busyMsg = "Экспортирую данные…"
				m.state = stateShowMessage
				return m, tea.Batch(m.spin.Tick, m.cmdExport(name, isCSV))
			}
			return m, cmd

		case stateShowMessage, stateLogs:
			switch msg.String() {
			case "enter", "esc":
				m.state = stateMain
			}

		case stateShowImportErrors:
			switch msg.String() {
			case "enter", "esc":
				m.state = stateImportSub
				m.subCursor = 2
			}
		}
	}

	return m, nil
}
