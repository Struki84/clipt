package main

import (
	"context"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func ShowUI(agent *Agent) {
	p := tea.NewProgram(initialModel(agent), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

var customBorder = lipgloss.Border{
	Left: "█", Right: "█",
	Top: "", Bottom: "",
	TopLeft: "", TopRight: "",
	BottomLeft: "", BottomRight: "",
}

var viewportStyle = lipgloss.NewStyle().Border(customBorder).Padding(1)
var inputStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.ThickBorder()).
	BorderTop(true).BorderRight(false).
	BorderLeft(false).BorderBottom(false)

type model struct {
	agent        *Agent
	viewport     viewport.Model
	messages     []string
	textarea     textarea.Model
	senderStyle  lipgloss.Style
	err          error
	windowWidth  int
	windowHeight int
}

func initialModel(agent *Agent) model {
	vp := viewport.New(120, 30)
	vp.SetContent(`Welcome to the chat room! Type a message and press Enter to send.`)

	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false

	return model{
		agent:        agent,
		textarea:     ta,
		messages:     []string{},
		viewport:     vp,
		senderStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		windowWidth:  120,
		windowHeight: 30,
	}
}

func (m model) Init() tea.Cmd {
	response := ""
	m.agent.Stream(context.Background(), func(ctx context.Context, chunk []byte) {
		response += string(chunk)
		m.messages = append(m.messages, m.senderStyle.Render("Clipt: ")+response)
		m.viewport.SetContent(strings.Join(m.messages, "\n"))
		m.viewport.GotoBottom()
	})

	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.textarea.Focused() && strings.TrimSpace(m.textarea.Value()) != "" {
				m.messages = append(m.messages, m.senderStyle.Render("You: ")+m.textarea.Value())
				m.viewport.SetContent(strings.Join(m.messages, "\n"))
				m.textarea.Reset()
				m.viewport.GotoBottom()
			}
		}
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		contentWidth := min(m.windowWidth-20, 120)
		m.viewport.Width = contentWidth
		m.viewport.Height = m.windowHeight - m.textarea.Height() - 10
		m.textarea.SetWidth(contentWidth)
		m.viewport.SetContent(strings.Join(m.messages, "\n"))
	}

	ta, cmd := m.textarea.Update(msg)
	m.textarea = ta
	cmds = append(cmds, cmd)

	vp, cmd := m.viewport.Update(msg)
	m.viewport = vp
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return lipgloss.Place(
		m.windowWidth,
		m.windowHeight,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			viewportStyle.Render(m.viewport.View()),
			inputStyle.Render(m.textarea.View()),
		),
	)
}
