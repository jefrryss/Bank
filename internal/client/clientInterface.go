package client

import (
	"fmt"
	"task2/internal/bankmanager"

	tea "github.com/charmbracelet/bubbletea"
)

type MenuItem struct {
	ID    string
	Label string
}

type model struct {
	manager       *bankmanager.BankManager
	menu          []MenuItem
	cursor        int
	activeID      string
	state         string
	input         string
	message       string
	PerformAction func(activeID string) string
}

func NewModel(manager *bankmanager.BankManager) model {
	menu := []MenuItem{
		{ID: "import", Label: "Импорт данных"},
		{ID: "export", Label: "Экспорт данных"},
		{ID: "analytics_balance", Label: "Аналитика: Баланс за период"},
		{ID: "analytics_category", Label: "Аналитика: Группировка по категориям"},
		{ID: "view_errors", Label: "Показать ошибки"},
	}
	return model{
		manager: manager,
		menu:    menu,
		state:   "menu",
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.menu)-1 {
				m.cursor++
			}
		case "enter":
			m.activeID = m.menu[m.cursor].ID
			if m.PerformAction != nil {
				m.message = m.PerformAction(m.activeID)
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	s := "Меню:\n"
	for i, item := range m.menu {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, item.Label)
	}
	s += "\nИспользуйте стрелки ↑↓ и Enter для выбора, q для выхода\n"
	if m.message != "" {
		s += "\n" + m.message + "\n"
	}
	return s
}
