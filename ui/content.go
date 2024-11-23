package ui

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
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

	viewportStyle = lipgloss.NewStyle().Padding(1)
	inputStyle    = lipgloss.NewStyle().
			BorderStyle(lipgloss.ThickBorder()).
			BorderTop(true).BorderRight(false).
			BorderLeft(false).BorderBottom(false)
)

type ContentView struct {
	Style    lipgloss.Style
	viewport viewport.Model
	textarea textarea.Model
}

func NewContentView() ContentView {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.ShowLineNumbers = false
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.SetWidth(120)
	ta.SetHeight(15)
	ta.Focus()

	return ContentView{
		Style:    lipgloss.NewStyle().Border(customBorder),
		viewport: viewport.New(120, 56),
		textarea: ta,
	}
}

func (content ContentView) Init() tea.Cmd {
	return nil
}

func (content ContentView) View() string {
	chatView := lipgloss.JoinVertical(
		lipgloss.Left,
		viewportStyle.Render(content.viewport.View()),
		inputStyle.Render(content.textarea.View()),
	)
	return content.Style.PaddingLeft(1).Render(chatView)
}

func (content ContentView) Update(msg tea.Msg) (ContentView, tea.Cmd) {
	return content, nil
}
