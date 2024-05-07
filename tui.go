package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

func ShowUI() {
	p := tea.NewProgram(intialModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type model struct {
	viewport      viewport.Model
	messages      []string
	textarea      textarea.Model
	senderStyle   lipgloss.Style
	viewportStyle lipgloss.Style
}

func intialModel() model {
	ta := textarea.New()
	ta.Placeholder = "Enter text here"
	ta.Focus()

	ta.Prompt = "> "
	ta.SetWidth(120)
	ta.SetHeight(8)

	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false

	ta.KeyMap.InsertNewline.SetEnabled(false)

	vp := viewport.New(120, 8)
	vp.SetContent("This is the content of the viewport")

	vpStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.Border{Left: "|", Right: "|"}).
		BorderStyle(lipgloss.DoubleBorder())

	return model{
		viewport:      vp,
		messages:      []string{},
		textarea:      ta,
		senderStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		viewportStyle: vpStyle,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			m.messages = append(m.messages, m.senderStyle.Render("You > "+m.textarea.Value()))
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.textarea.Reset()
			m.viewport.GotoBottom()
		}
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	return m.viewportStyle.Render(m.viewport.View())
}
