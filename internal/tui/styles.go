package tui

import "github.com/charmbracelet/lipgloss"

// Shared color palette and lip gloss styles used across the dashboard and the
// setup wizard. Centralizing them keeps the TUI's visual language consistent.
var (
	accentColor  = lipgloss.Color("63") // Violet
	textColor    = lipgloss.Color("255")
	errorColor   = lipgloss.Color("196")
	successColor = lipgloss.Color("40")
	grayColor    = lipgloss.Color("242")
	amberColor   = lipgloss.Color("214")

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(textColor).
			Background(accentColor).
			Padding(0, 1).
			MarginBottom(1)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accentColor).
			Underline(true)

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(textColor).
			Background(accentColor)

	normalStyle = lipgloss.NewStyle().
			Foreground(textColor)

	errorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(errorColor)

	successStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(successColor)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(accentColor).
			Padding(1, 2).
			MarginBottom(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(grayColor).
			Italic(true)

	tipStyle = lipgloss.NewStyle().
			Foreground(amberColor).
			Italic(true)

	labelStyle = lipgloss.NewStyle().
			Foreground(grayColor).
			Bold(true)

	statusBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(grayColor).
			Padding(0, 1).
			MarginBottom(1)
)
