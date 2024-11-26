package ui

import (
	"log"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	viewportStyle = lipgloss.NewStyle().Padding(1)
	inputStyle    = lipgloss.NewStyle().
			BorderStyle(lipgloss.ThickBorder()).
			BorderTop(true).
			Height(10)
)

type ChatView struct {
	windowSize tea.WindowSizeMsg
	viewport   viewport.Model
	textarea   textarea.Model
}

func NewChatView() ChatView {
	ta := textarea.New()

	ta.Placeholder = "Send a message..."
	ta.ShowLineNumbers = false
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.SetHeight(10)
	ta.Focus()

	return ChatView{
		viewport: viewport.New(120, 35),
		textarea: ta,
	}
}

func (chat ChatView) View() string {
	chat.textarea.BlurredStyle.Base.BorderTop(true)

	joinVertical := lipgloss.JoinVertical(
		lipgloss.Center,
		viewportStyle.Render(chat.viewport.View()),
		chat.textarea.View(),
	)

	return joinVertical
}

func (chat ChatView) Init() tea.Cmd {
	return nil
}

func (chat ChatView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		log.Printf("ChatView: %#v", msg)
		chat.windowSize = msg
		viewportStyle.Width(msg.Width)
		viewportStyle.Height(msg.Height - chat.textarea.Height() - 1)
		chat.textarea.SetWidth(msg.Width)
	}

	return chat, nil
}
