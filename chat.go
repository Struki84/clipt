// TODO
// - load chat history
// - top bar with session name / timestamp
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
	"os/user"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type StreamMsg struct {
	Chunk string
}

type ChatView struct {
	user       *user.User
	provider   ChatProvider
	msgs       []ChatMsg
	streamChan chan string
	viewport   viewport.Model
	textarea   textarea.Model
	renderer   *glamour.TermRenderer
}

func NewChatViewLight(agent ChatProvider) ChatView {
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
		glamour.WithWordWrap(240),
	)

	vp := viewport.New(120, 35)
	vp.KeyMap = viewport.KeyMap{
		PageDown: key.NewBinding(key.WithKeys("pgdown", "space")),
		PageUp:   key.NewBinding(key.WithKeys("pgup")),
		Down:     key.NewBinding(key.WithKeys("down")),
		Up:       key.NewBinding(key.WithKeys("up")),
	}

	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	return ChatView{
		user:       user,
		provider:   agent,
		streamChan: make(chan string, 100),
		viewport:   vp,
		textarea:   ta,
		renderer:   renderer,
	}
}

func (chat ChatView) Init() tea.Cmd {
	chat.provider.Stream(context.Background(), func(ctx context.Context, chunk []byte) error {
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

	return StreamMsg{
		Chunk: content,
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
		chat.viewport.Width = msg.Width
		chat.viewport.Height = msg.Height - chat.textarea.Height() - 3
		chat.textarea.SetWidth(msg.Width - 4)
		chat.viewport.SetContent(chat.renderMessages())

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return chat, tea.Quit

		case tea.KeyEnter:
			if chat.textarea.Focused() && strings.TrimSpace(chat.textarea.Value()) != "" {
				input := chat.textarea.Value()

				usrMsg := ChatMsg{
					Content:   input,
					Role:      "User",
					Timestamp: time.Now().Unix(),
				}

				chat.msgs = append(chat.msgs, usrMsg)

				aiMsg := ChatMsg{
					Content:   "",
					Role:      "AI",
					Timestamp: time.Now().Unix(),
				}

				chat.msgs = append(chat.msgs, aiMsg)

				chat.viewport.SetContent(chat.renderMessages())
				chat.textarea.Reset()
				chat.viewport.GotoBottom()

				go func() {
					err := chat.provider.Run(context.Background(), input)
					if err != nil {
						log.Println("Run error:", err)
					}
				}()

				return chat, nil
			}
		}

	case StreamMsg:
		chat.msgs[len(chat.msgs)-1].Content += msg.Chunk
		chat.msgs[len(chat.msgs)-1].Timestamp = time.Now().Unix()

		chat.viewport.SetContent(chat.renderMessages())
		chat.viewport.GotoBottom()

		return chat, chat.handleStream
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

	for _, msg := range chat.msgs {
		if msg.Role == "User" {
			date := time.Unix(msg.Timestamp, 0).Format("2 Jan 2006 15:04")
			username := chat.user.Username
			fullMsg := fmt.Sprintf("%s\n\n%s (%s) ", msg.Content, username, date)

			styledMessages = append(styledMessages, userStyle.Render(fullMsg))

		} else if msg.Role == "AI" {
			renderedTxt, _ := chat.renderer.Render(msg.Content)

			styledMessages = append(styledMessages, aiStyle.Render(renderedTxt))
		}
	}

	return lipgloss.JoinVertical(lipgloss.Center, styledMessages...)
}

func ShowChatViewLight(agent ChatProvider) {
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
