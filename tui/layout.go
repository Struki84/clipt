package tui

import (
	"context"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type StreamMsg struct {
	Chunk string
}

type LayoutView struct {
	Provider ChatProvider
	Session  ChatHistory
	Msgs     []ChatMsg
}

func NewLayoutView(provider ChatProvider) LayoutView {
	return LayoutView{
		Provider: provider,
	}
}

func (layout LayoutView) Init() tea.Cmd {
	cmds := []tea.Cmd{}

	return tea.Batch(cmds...)
}

func (layout LayoutView) View() string {
	// generate a container for all the elements and join them?

	return ""
}

func (layout LayoutView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return layout, nil
}

// Components
// - Session Bar
// - Chat View
// - Chat Input
// - Status Bar
// - Chat Menu

type ChatView struct {
	Style    lipgloss.Style
	ViewPort viewport.Model
}

func NewChatView() ChatView {
	return ChatView{}
}

func (chat ChatView) Init() tea.Cmd {
	return nil
}

func (chat ChatView) View() string {
	return ""
}

func (chat ChatView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	return chat, nil
}

type ChatInput struct {
	Style  *lipgloss.Style
	Prompt textarea.Model
}

func NewChatInput() ChatInput {
	style := lipgloss.NewStyle().
		Background(lipgloss.Color("#11111b")).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#ffffff")).
		BorderLeft(true).
		BorderRight(true).
		BorderTop(false).
		BorderBottom(false).
		Padding(1).
		MarginBottom(1)

	ta := textarea.New()
	ta.Prompt = ""
	ta.ShowLineNumbers = false
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.SetHeight(1)
	ta.Focus()
	ta.FocusedStyle.Base = style

	return ChatInput{
		Style:  &style,
		Prompt: ta,
	}
}

func (input ChatInput) Init() tea.Cmd {
	return textarea.Blink
}

func (input ChatInput) View() string {
	return input.Prompt.View()
}

func (input ChatInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		input.Prompt.SetWidth(msg.Width - 4)
	}

	return input, nil
}
