// TODO
// - load chat history
// - top bar with session name / timestamp
// - add user msg timestapm and handle
// - add error msgs
// - tool msgs
// - add some kind of debug view for reading reasoning steps
// - bottom bar with active mode and active engine
// - commands menu

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type ChatMsgs struct {
	Content string
	Role    string
}

type ChatView struct {
	agent      AIEngine
	messages   []string
	streamChan chan string
	viewport   viewport.Model
	textarea   textarea.Model
	windowSize tea.WindowSizeMsg
	renderer   *glamour.TermRenderer
}

func NewChatViewLight(agent AIEngine) ChatView {
	ta := textarea.New()
	ta.Prompt = ""
	ta.ShowLineNumbers = false
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.FocusedStyle.Base = lipgloss.NewStyle().
		Background(lipgloss.Color("#11111b")).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#ffffff")).
		BorderLeft(true).
		BorderRight(true).
		BorderTop(false).
		BorderBottom(false).
		Padding(1).
		MarginBottom(1)
	ta.SetHeight(1)
	ta.Focus()

	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(120),
	)

	vp := viewport.New(120, 35)
	vp.KeyMap = viewport.KeyMap{
		PageDown: key.NewBinding(key.WithKeys("pgdown", "space")),
		PageUp:   key.NewBinding(key.WithKeys("pgup")),
		Down:     key.NewBinding(key.WithKeys("down")),
		Up:       key.NewBinding(key.WithKeys("up")),
	}

	return ChatView{
		agent:      agent,
		messages:   make([]string, 0),
		streamChan: make(chan string, 100),
		viewport:   vp,
		textarea:   ta,
		renderer:   renderer,
	}
}

func (chat ChatView) Init() tea.Cmd {
	chat.agent.Stream(context.Background(), func(ctx context.Context, chunk []byte) error {
		chat.streamChan <- string(chunk)
		return nil
	})

	cmds := []tea.Cmd{}
	cmds = append(cmds, textarea.Blink)
	cmds = append(cmds, chat.handleStream)

	return tea.Batch(cmds...)
}

func (chat ChatView) handleStream() tea.Msg {
	content := <-chat.streamChan

	return ChatMsgs{
		Content: content,
	}
}

func (chat ChatView) View() string {
	joinVertical := lipgloss.JoinVertical(
		lipgloss.Center,
		chat.viewport.View(),
		chat.textarea.View(),
	)
	return joinVertical
}

func (chat ChatView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		chat.windowSize = msg
		chat.viewport.Width = msg.Width
		chat.viewport.Height = msg.Height - chat.textarea.Height() - 3
		chat.textarea.SetWidth(msg.Width - 4)
		chat.viewport.SetContent(chat.renderMessages())

	case ChatMsgs:
		chat.messages[len(chat.messages)-1] += msg.Content

		chat.viewport.SetContent(chat.renderMessages())
		chat.viewport.GotoBottom()

		return chat, chat.handleStream
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return chat, tea.Quit

		case tea.KeyEnter:
			if chat.textarea.Focused() && strings.TrimSpace(chat.textarea.Value()) != "" {
				userMsg := "You: " + chat.textarea.Value()
				chat.messages = append(chat.messages, userMsg)
				chat.messages = append(chat.messages, "Clipt: ")

				chat.viewport.SetContent(chat.renderMessages())
				chat.textarea.Reset()
				chat.viewport.GotoBottom()

				go func() {
					err := chat.agent.Run(context.Background(), userMsg)
					if err != nil {
						log.Println("Run error:", err)
					}
				}()
				return chat, nil
			}
		}
	}

	cmds := []tea.Cmd{}

	ta, cmd := chat.textarea.Update(msg)
	chat.textarea = ta
	cmds = append(cmds, cmd)

	vp, cmd := chat.viewport.Update(msg)
	chat.viewport = vp
	cmds = append(cmds, cmd)

	return chat, tea.Batch(cmds...)
}

func (chat ChatView) renderMessages() string {

	var styledMessages []string

	userStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#181825")).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#b4befe")).
		BorderLeft(true).
		BorderRight(true).
		BorderTop(false).
		BorderBottom(false).
		Padding(1).
		Margin(1).
		Align(lipgloss.Left).
		Width(chat.viewport.Width - 6)

	aiStyle := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Width(chat.viewport.Width - 2)

	for _, msg := range chat.messages {
		if strings.HasPrefix(msg, "You:") {
			content := strings.TrimPrefix(msg, "You: ")
			content = content + "\n" + "simun (17 Sep 2025 15:00)"
			styledMessages = append(styledMessages, userStyle.Render(content))
		} else if strings.HasPrefix(msg, "Clipt:") {
			renderedTxt, _ := chat.renderer.Render(strings.TrimPrefix(msg, "Clipt: "))
			styledMessages = append(styledMessages, aiStyle.Render(renderedTxt))
		}
	}

	return lipgloss.JoinVertical(lipgloss.Center, styledMessages...)
}

func ShowChatViewLight(agent AIEngine) {
	file, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	p := tea.NewProgram(
		NewChatViewLight(agent),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
