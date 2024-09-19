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
	file, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

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

type responseMsg struct {
	Content string
	Done    bool
}

type model struct {
	agent        *Agent
	viewport     viewport.Model
	messages     []string
	textarea     textarea.Model
	senderStyle  lipgloss.Style
	err          error
	streamChan   chan string
	windowWidth  int
	windowHeight int
}

func (m model) Init() tea.Cmd {
	m.agent.Stream(context.Background(), func(ctx context.Context, chunk []byte) {
		// log.Println("Stream chunk:", string(chunk))
		m.streamChan <- string(chunk)
		// m.streamChan <- responseMsg{
		// 	Content: string(chunk),
		// 	Done:    false,
		// }
	})

	return textarea.Blink
}

func (m model) handleStream() tea.Msg {
	log.Println("handleStream:", <-m.streamChan)
	return responseMsg{
		Content: <-m.streamChan,
		Done:    false,
	}
	// return <-m.streamChan
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
		streamChan:   make(chan string),
		windowWidth:  120,
		windowHeight: 30,
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		contentWidth := min(m.windowWidth-20, 120)
		m.viewport.Width = contentWidth
		m.viewport.Height = m.windowHeight - m.textarea.Height() - 10
		m.textarea.SetWidth(contentWidth + 4)
		m.viewport.SetContent(strings.Join(m.messages, "\n"))

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.textarea.Focused() && strings.TrimSpace(m.textarea.Value()) != "" {
				userMessage := m.textarea.Value()
				m.messages = append(m.messages, m.senderStyle.Render("You: ")+userMessage)
				m.messages = append(m.messages, m.senderStyle.Render("Clipt: "))
				m.viewport.SetContent(strings.Join(m.messages, "\n"))
				m.textarea.Reset()
				m.viewport.GotoBottom()

				// was going to trasnfer this to (m *model)Init() - but then my cat destryed my mbp
				// go func() {
				// ctx, cancel := context.WithCancel(context.Background())
				// defer cancel()
				//
				// m.agent.Stream(ctx, func(ctx context.Context, chunk []byte) {
				// 	log.Println("Stream: ", string(chunk))
				// 	m.streamChan <- responseMsg{Content: string(chunk), Done: false}
				// })

				// m.streamChan <- responseMsg{Done: true}

				// log.Println("Stream is done")

				// }()

				go func() {
					log.Println("Run: ", userMessage)
					err := m.agent.Run(context.Background(), userMessage)
					if err != nil {
						log.Println("Run error:", err)
					}
				}()

				return m, m.handleStream
			}
		}

	case responseMsg:
		if !msg.Done {
			// m.messages[len(m.messages)-1] += msg.Content
			// m.viewport.SetContent(strings.Join(m.messages, "\n"))
			// m.viewport.GotoBottom()

			return m, m.handleStream
		}
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
