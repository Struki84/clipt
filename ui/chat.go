package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	customBorder = lipgloss.Border{
		Left: "█", Right: "",
		Top: "", Bottom: "",
		TopLeft: "", TopRight: "",
		BottomLeft: "", BottomRight: "",
	}
)

type ContentView struct {
	Style lipgloss.Style
}

func NewContentView() ContentView {
	return ContentView{
		Style: lipgloss.NewStyle().Border(customBorder),
	}
}

func (content ContentView) Init() tea.Cmd {
	return nil
}

func (content ContentView) View() string {
	return content.Style.Render("Content View")
}
