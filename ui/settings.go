package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	settingsStyle = lipgloss.NewStyle().
		Padding(1)
)

type SettingsView struct {
}

func NewSettingsView() SettingsView {
	return SettingsView{}
}

func (settings SettingsView) Init() tea.Cmd {
	return nil
}

func (settings SettingsView) View() string {
	return settingsStyle.Render("Settings View")
}

func (settings SettingsView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		settingsStyle.Width(msg.Width)
		settingsStyle.Height(msg.Height)
	}
	return settings, nil
}
