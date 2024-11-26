package ui

import tea "github.com/charmbracelet/bubbletea"

type SettingsView struct {
}

func NewSettingsView() SettingsView {
	return SettingsView{}
}

func (settings SettingsView) Init() tea.Cmd {
	return nil
}

func (settings SettingsView) View() string {
	return "Settings View"
}

func (settings SettingsView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return settings, nil
}
