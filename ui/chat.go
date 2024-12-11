package ui

import (
	"context"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/struki84/clipt/internal"
)

type ChatMsgs struct {
	Content string
}

type ChatView struct {
	agent      *internal.Agent
	messages   []string
	streamChan chan string
	viewport   viewport.Model
	textarea   textarea.Model
	windowSize tea.WindowSizeMsg
	renderer   *glamour.TermRenderer
}

func NewChatView(agent *internal.Agent) ChatView {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.ShowLineNumbers = false
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.SetHeight(10)
	ta.Focus()

	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(120),
	)

	return ChatView{
		agent:      agent,
		messages:   make([]string, 0),
		streamChan: make(chan string),
		viewport:   viewport.New(120, 35),
		textarea:   ta,
		renderer:   renderer,
	}
}

func (chat ChatView) Init() tea.Cmd {
	chat.agent.Stream(context.Background(), func(ctx context.Context, chunk []byte) {
		chat.streamChan <- string(chunk)
	})
	cmds := []tea.Cmd{}
	cmds = append(cmds, textarea.Blink)
	cmds = append(cmds, chat.handleStream)
	return tea.Batch(cmds...)
}

func (chat ChatView) handleStream() tea.Msg {
	// log.Println("handleStream")
	content := <-chat.streamChan

	return ChatMsgs{
		Content: content,
	}
}

func (chat ChatView) View() string {
	joinVertical := lipgloss.JoinVertical(
		lipgloss.Left,
		chat.viewport.View(),
		chat.textarea.View(),
	)
	return joinVertical
}

func (chat ChatView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// log.Printf("ChatView handling WindowSizeMsg")
		chat.windowSize = msg
		chat.viewport.Width = msg.Width
		chat.viewport.Height = msg.Height - chat.textarea.Height() - 1
		chat.textarea.SetWidth(msg.Width)
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
				chat.messages = append(chat.messages, "Clipt:")

				chat.viewport.SetContent(chat.renderMessages())
				chat.textarea.Reset()
				chat.viewport.GotoBottom()

				go func() {
					log.Println("Run: ", userMsg)
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
	// log.Println("renderMessages called")
	// Join messages with double newlines for proper separation
	messageContent := strings.Join(chat.messages, "\n\n")

	// Add proper markdown formatting for messages
	formattedContent := strings.ReplaceAll(messageContent, "You: ", "### You:\n")
	formattedContent = strings.ReplaceAll(formattedContent, "Clipt: ", "### Clipt:\n")

	// formattedContent := messageContent

	rendered, err := chat.renderer.Render(formattedContent)
	if err != nil {
		log.Printf("Error rendering messages: %v", err)
		return messageContent
	}
	return rendered

	// return formattedContent
}
