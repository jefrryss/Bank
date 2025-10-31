package importer

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// pathInputModel - модель для ввода пути
type pathInputModel struct {
	textInput textinput.Model
	spinner   spinner.Model
	err       error
	filePath  string
	done      bool
	cancelled bool
	loading   bool // Флаг загрузки
}

func newPathInputModel() pathInputModel {
	ti := textinput.New()
	ti.Placeholder = "Введите путь к JSON файлу..."
	ti.Focus()
	ti.CharLimit = 512
	ti.Width = 60

	// Настраиваем спиннер
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return pathInputModel{
		textInput: ti,
		spinner:   s,
		loading:   false,
	}
}

func (m pathInputModel) Init() tea.Cmd {
	return textinput.Blink
}

// Сообщение о завершении валидации
type validateDoneMsg struct {
	filePath string
	err      error
}

// Команда для валидации файла (выполняется в фоне)
func validateFileCmd(filePath string) tea.Cmd {
	return func() tea.Msg {
		// Убираем кавычки
		filePath = strings.Trim(filePath, `"'`)

		if filePath == "" {
			return validateDoneMsg{err: errors.New("путь не может быть пустым")}
		}

		fileInfo, err := os.Stat(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				return validateDoneMsg{err: fmt.Errorf("файл не существует: %s", filePath)}
			}
			return validateDoneMsg{err: fmt.Errorf("ошибка проверки файла: %v", err)}
		}

		if fileInfo.IsDir() {
			return validateDoneMsg{err: fmt.Errorf("указан путь к директории, а не к файлу: %s", filePath)}
		}

		return validateDoneMsg{filePath: filePath, err: nil}
	}
}

func (m pathInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case validateDoneMsg:
		// Валидация завершена
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		// Успех!
		m.filePath = msg.filePath
		m.done = true
		return m, tea.Quit

	case tea.KeyMsg:
		if m.loading {
			// Во время загрузки игнорируем ввод
			return m, nil
		}

		switch msg.Type {
		case tea.KeyEnter:
			filePath := strings.TrimSpace(m.textInput.Value())

			// Включаем режим загрузки
			m.loading = true
			m.err = nil

			// Запускаем валидацию и спиннер
			return m, tea.Batch(
				validateFileCmd(filePath),
				m.spinner.Tick,
			)

		case tea.KeyCtrlC, tea.KeyEsc:
			m.cancelled = true
			return m, tea.Quit
		}

	// Обновление спиннера
	case spinner.TickMsg:
		if m.loading {
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	// Обновление текстового поля
	if !m.loading {
		m.textInput, cmd = m.textInput.Update(msg)
	}
	return m, cmd
}

func (m pathInputModel) View() string {
	if m.loading {
		// Показываем спиннер во время загрузки
		return fmt.Sprintf("\n  %s Проверка файла...\n\n", m.spinner.View())
	}

	s := "\n  Введите путь к файлу:\n\n"
	s += "  " + m.textInput.View() + "\n\n"

	if m.err != nil {
		s += fmt.Sprintf("  \033[31mОшибка: %v\033[0m\n\n", m.err)
	}

	s += "  (Enter - подтвердить | Esc - отмена)\n"
	return s
}
