package app

import "github.com/charmbracelet/lipgloss"

var (
	clTitle   = lipgloss.Color("#A78BFA")
	clAccent  = lipgloss.Color("#10B981")
	clInfo    = lipgloss.Color("#93C5FD")
	clText    = lipgloss.Color("#E5E7EB")
	clDim     = lipgloss.Color("#6B7280")
	clWarn    = lipgloss.Color("#F59E0B")
	clSuccess = lipgloss.Color("#22C55E")
	clError   = lipgloss.Color("#EF4444")
	clBorder  = lipgloss.Color("#3F3F46")

	styleTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(clTitle)

	styleHint = lipgloss.NewStyle().
			Foreground(clDim)

	styleMenuCursor = lipgloss.NewStyle().
			Foreground(clAccent).
			Bold(true)

	styleMenuItem = lipgloss.NewStyle().
			Foreground(clText)

	styleMenuItemSelected = lipgloss.NewStyle().
				Foreground(clAccent).
				Bold(true)

	stylePrompt = lipgloss.NewStyle().
			Foreground(clInfo).
			Bold(true)

	styleInput = lipgloss.NewStyle().
			Foreground(clText)

	styleSuccess = lipgloss.NewStyle().
			Foreground(clSuccess).
			Bold(true)

	styleError = lipgloss.NewStyle().
			Foreground(clError).
			Bold(true)

	styleSpinner = lipgloss.NewStyle().
			Foreground(clWarn).
			Bold(true)

	styleBox = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.NormalBorder()).
			BorderForeground(clBorder)
)
