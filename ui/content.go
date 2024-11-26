package ui

import (
	"log"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ContentView struct {
	Style        lipgloss.Style
	viewport     viewport.Model
	textarea     textarea.Model
	selectedView string
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
	log.Println("View:", content.selectedView)
	if content.selectedView == "CHAT" {
		chatView := lipgloss.JoinVertical(
			lipgloss.Left,
			viewportStyle.Render(content.viewport.View()),
			inputStyle.Render(content.textarea.View()),
		)
		return content.Style.PaddingLeft(1).Render(chatView)
	}

	if content.selectedView == "HISTORY" {
		return content.Style.PaddingLeft(1).Render("History View")
	}

	if content.selectedView == "SETTINGS" {
		return content.Style.PaddingLeft(1).Render("Settings View")
	}

	return ""

}

func (content ContentView) Update(msg tea.Msg) (ContentView, tea.Cmd) {

	return content, nil
}

func (content ContentView) SetContent(view string) {
	content.selectedView = view
}
