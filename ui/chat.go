package ui

import (
	"context"
	"log"
	"regexp"
	"strings"

	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/struki84/clipt/internal"
)

var (
	inputStyle    = lipgloss.NewStyle()
	viewportStyle = lipgloss.NewStyle().Padding(1)
	senderStyle   = lipgloss.NewStyle()
	agentStyle    = lipgloss.NewStyle()
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
	mdRenderer *glamour.TermRenderer
	windowSize tea.WindowSizeMsg
}

func NewChatView(agent *internal.Agent) ChatView {
	ta := textarea.New()

	ta.Placeholder = "Send a message..."
	ta.ShowLineNumbers = false
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.SetHeight(10)
	ta.Focus()

	mdRenderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(120),
	)

	return ChatView{
		agent:      agent,
		messages:   make([]string, 0),
		streamChan: make(chan string),
		viewport:   viewport.New(120, 35),
		textarea:   ta,
		mdRenderer: mdRenderer,
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
	return ChatMsgs{
		Content: <-chat.streamChan,
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

func (chat ChatView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		chat.windowSize = msg
		chat.viewport.Width = msg.Width
		chat.viewport.Height = msg.Height - chat.textarea.Height() - 2
		chat.textarea.SetWidth(msg.Width)
		chat.viewport.SetContent(strings.Join(chat.messages, "\n"))
	case ChatMsgs:
		formatted, _ := chat.formatCodeBlock(msg.Content)
		chat.messages[len(chat.messages)-1] += formatted
		chat.viewport.SetContent(strings.Join(chat.messages, "\n"))
		chat.viewport.GotoBottom()

		return chat, chat.handleStream
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if chat.textarea.Focused() && strings.TrimSpace(chat.textarea.Value()) != "" {
				userMsg := chat.textarea.Value()

				chat.messages = append(chat.messages, senderStyle.Render("You: "+userMsg))
				chat.messages = append(chat.messages, agentStyle.Render("Clipt: "))

				chat.viewport.SetContent(strings.Join(chat.messages, "\n"))
				chat.viewport.GotoBottom()
				chat.textarea.Reset()

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

	ta, cmd := chat.textarea.Update(msg)
	chat.textarea = ta
	cmds = append(cmds, cmd)

	vp, cmd := chat.viewport.Update(msg)
	chat.viewport = vp
	cmds = append(cmds, cmd)

	return chat, tea.Batch(cmds...)
}

// Helper function to detect and highlight code blocks
func (chat ChatView) formatCodeBlock(content string) (string, error) {
	// Regular expression to find code blocks with language specification
	// Matches ```language\n...code...```
	codeBlockRegex := regexp.MustCompile("```([a-zA-Z0-9]+)\n([^`]+)```")

	formatted := codeBlockRegex.ReplaceAllStringFunc(content, func(block string) string {
		matches := codeBlockRegex.FindStringSubmatch(block)
		if len(matches) < 3 {
			return block
		}

		language := matches[1]
		code := matches[2]

		// Get lexer for the specified language
		lexer := lexers.Get(language)
		if lexer == nil {
			lexer = lexers.Fallback
		}

		// Tokenize the code
		iterator, err := lexer.Tokenise(nil, code)
		if err != nil {
			return block
		}

		formatter := formatters.Get("terminal256")
		if formatter == nil {
			formatter = formatters.Fallback
		}

		style := styles.Get("monokai")
		if style == nil {
			style = styles.Fallback
		}

		var buf strings.Builder
		err = formatter.Format(&buf, style, iterator)
		if err != nil {
			return block
		}

		return buf.String()
	})

	return formatted, nil
}
