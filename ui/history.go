package ui

import tea "github.com/charmbracelet/bubbletea"

type HistoryView struct {
}

func NewHistoryView() HistoryView {
	return HistoryView{}
}

func (history HistoryView) Init() tea.Cmd {
	return nil
}

func (history HistoryView) View() string {
	return "History View"
}

func (history HistoryView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return history, nil
}
