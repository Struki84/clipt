package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	historyStyle = lipgloss.NewStyle().Padding(1)
)

type HistoryView struct {
}

func NewHistoryView() HistoryView {
	return HistoryView{}
}

func (history HistoryView) Init() tea.Cmd {
	return nil
}

func (history HistoryView) View() string {
	return historyStyle.Render("History View")
}

func (history HistoryView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		historyStyle.Width(msg.Width)
		historyStyle.Height(msg.Height)
	}

	return history, nil
}
