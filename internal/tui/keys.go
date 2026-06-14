package tui

import "github.com/charmbracelet/bubbles/key"

// keyMap defines the operator console keybindings and powers the help bar.
type keyMap struct {
	Up      key.Binding
	Down    key.Binding
	Enter   key.Binding
	Back    key.Binding
	Refresh key.Binding
	Abbrev  key.Binding
	Config  key.Binding
	Help    key.Binding
	Quit    key.Binding
}

// ShortHelp is shown in the single-line help bar.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.Back, k.Refresh, k.Help, k.Quit}
}

// FullHelp is shown when the help bar is expanded with "?".
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Enter, k.Back},
		{k.Refresh, k.Abbrev, k.Config},
		{k.Help, k.Quit},
	}
}

func defaultKeyMap() keyMap {
	return keyMap{
		Up:      key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
		Down:    key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
		Enter:   key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "open account")),
		Back:    key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
		Refresh: key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "refresh")),
		Abbrev:  key.NewBinding(key.WithKeys("h"), key.WithHelp("h", "balance guide")),
		Config:  key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "mcp client config")),
		Help:    key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
		Quit:    key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	}
}
