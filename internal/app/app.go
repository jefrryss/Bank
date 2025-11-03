package app

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"task2/internal/analitic"
	"task2/internal/bankmanager"
	"task2/internal/export"
	"task2/internal/handmaker"
	"task2/internal/importer"
	"task2/internal/logstorage"
)

type app struct {
	state     int
	cursor    int
	subCursor int
	textInput textinput.Model

	infoMessage string

	analyticsFrom string
	analyticsTo   string

	manualAccountID string
	manualBalance   float64

	logStore        *logstorage.LogStorage
	bankManager     *bankmanager.BankManager
	importFacade    *importer.ImportFacade
	exportFacade    *export.ExportFacade
	analyticsFacade *analitic.AnalyticsFacade

	ctrl         *handmaker.Controller
	manualEntity string
	manualAction string
	formFields   []handmaker.Field
	formValues   map[string]string
	formStep     int

	lastImportErrors []string

	spin    spinner.Model
	isBusy  bool
	busyMsg string
}

func NewApp(
	ls *logstorage.LogStorage,
	bm *bankmanager.BankManager,
	imp *importer.ImportFacade,
	ex *export.ExportFacade,
	an *analitic.AnalyticsFacade,
	ctrl *handmaker.Controller,
) tea.Model {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()
	ti.CharLimit = 512
	ti.Width = 50

	sp := spinner.New()
	sp.Spinner = spinner.Dot

	return &app{
		state:           stateMain,
		cursor:          0,
		subCursor:       0,
		textInput:       ti,
		logStore:        ls,
		bankManager:     bm,
		importFacade:    imp,
		exportFacade:    ex,
		analyticsFacade: an,
		ctrl:            ctrl,
		spin:            sp,
	}
}

func (m *app) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.spin.Tick)
}
