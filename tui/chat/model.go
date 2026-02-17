package chat

import (
	"context"
	"fmt"
	"log"
	"os/user"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/struki84/clipt/tui/schema"
)

type ChatView struct {
	WindowSize tea.WindowSizeMsg

	Style    Styles
	Msgs     []schema.Msg
	Stream   chan schema.Msg
	Provider schema.ChatProvider
	Session  schema.ChatSession

	IsLoading bool

	Input    *textarea.Model
	Viewport *viewport.Model
	Loader   spinner.Model
}

func New(provider schema.ChatProvider) ChatView {
	input := textarea.New()
	input.Focus()

	view := viewport.New(0, 0)
	loader := spinner.New()
	loader.Spinner = spinner.Dot

	return ChatView{
		Provider:  provider,
		Input:     &input,
		Viewport:  &view,
		Style:     DefaultStyles(),
		Loader:    loader,
		IsLoading: false,
		Stream:    make(chan schema.Msg),
	}
}

func (chat ChatView) Init() tea.Cmd {
	cmds := []tea.Cmd{}
	cmds = append(cmds, textarea.Blink)
	cmds = append(cmds, chat.HandleStream)

	chat.Provider.Stream(context.TODO(), func(ctx context.Context, msg schema.Msg) error {
		chat.Stream <- msg
		return nil
	})

	return tea.Batch(cmds...)
}

func (chat ChatView) View() string {
	chat.Viewport.Width = chat.WindowSize.Width - 2
	chat.Viewport.Height = chat.WindowSize.Height - 8

	chat.Viewport.KeyMap = viewport.KeyMap{
		PageDown: key.NewBinding(key.WithKeys("pgdown")),
		PageUp:   key.NewBinding(key.WithKeys("pgup")),
		Down:     key.NewBinding(key.WithKeys("down")),
		Up:       key.NewBinding(key.WithKeys("up")),
	}

	chat.Viewport.SetContent(chat.RenderMsgs())
	// chat.Viewport.GotoBottom()

	chat.Input.Prompt = ""
	chat.Input.SetHeight(1)
	chat.Input.SetWidth(chat.WindowSize.Width - 4)
	chat.Input.FocusedStyle.CursorLine = lipgloss.NewStyle()
	chat.Input.FocusedStyle.Base = chat.Style.Input
	chat.Input.ShowLineNumbers = false

	chat.Loader.Style.PaddingBottom(1)

	return ""
}

func (chat ChatView) RenderMsgs() string {
	var styledMessages []string

	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(120),
	)

	width := chat.Viewport.Width - 4

	for i, msg := range chat.Msgs {
		switch msg.Role {
		case schema.InternalMsg:
			fullMsg := fmt.Sprintf("%s", msg.Content)
			chatMsg := chat.Style.Msg.Internal.Width(width).Render(fullMsg)

			styledMessages = append(styledMessages, chatMsg)
		case schema.ErrMsg:
			fullMsg := fmt.Sprintf("%s", msg.Content)
			chatMsg := chat.Style.Msg.Err.Width(width).Render(fullMsg)

			styledMessages = append(styledMessages, chatMsg)
		case schema.SysMsg:
			fullMsg := fmt.Sprintf("%s", msg.Content)
			chatMsg := chat.Style.Msg.Sys.Width(width).Render(fullMsg)

			styledMessages = append(styledMessages, chatMsg)
		case schema.UserMsg:
			date := time.Unix(msg.Timestamp, 0).Format("2 Jan 2006 15:04")
			username := user.Username
			fullMsg := fmt.Sprintf("%s\n\n%s (%s) ", msg.Content, username, date)
			chatMsg := chat.Style.Msg.User.Width(width).Render(fullMsg)

			styledMessages = append(styledMessages, chatMsg)
		case schema.AIMsg:
			var renderedTxt string

			if chat.IsLoading && msg.Content == "" && i == len(chat.Msgs)-1 {
				renderedTxt = chat.Loader.View() + " Working..."
			} else {
				renderedTxt, _ = renderer.Render(msg.Content)
			}

			chatMsg := chat.Style.Msg.AI.Width(width).Render(renderedTxt)

			styledMessages = append(styledMessages, chatMsg)
		}
	}

	return lipgloss.JoinVertical(lipgloss.Center, styledMessages...)
}

func (chat ChatView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		chat.WindowSize = msg
	case spinner.TickMsg:
		loader, cmd := chat.Loader.Update(msg)
		chat.Loader = loader
		cmds = append(cmds, cmd)
	case schema.Msg:
		if chat.IsLoading {
			chat.IsLoading = false
		}

		lastMsg := chat.Msgs[len(chat.Msgs)-1]
		if msg.Stream {
			if !lastMsg.Stream {
				aiMsg := schema.Msg{
					Stream:    true,
					Content:   "",
					Role:      schema.AIMsg,
					Timestamp: time.Now().Unix(),
				}

				chat.Msgs = append(chat.Msgs, aiMsg)
				lastMsg = aiMsg
			}

			lastMsg.Role = msg.Role
			lastMsg.Content += msg.Content
			lastMsg.Timestamp = time.Now().Unix()

			chat.Msgs[len(chat.Msgs)-1] = lastMsg
		} else {
			chat.Msgs = append(
				chat.Msgs,
				schema.Msg{
					Role:      msg.Role,
					Content:   msg.Content,
					Timestamp: time.Now().Unix(),
				},
			)
		}

		chat.Viewport.SetContent(chat.RenderMsgs())
		chat.Viewport.GotoBottom()

		return chat, chat.HandleStream

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			prompt := chat.Input.Value()
			menuActive := strings.HasPrefix(prompt, "/")
			if !menuActive && chat.Input.Focused() {
				input := chat.Input.Value()

				chat.Input.Reset()
				chat.IsLoading = true

				userMsg := schema.Msg{
					Stream:    false,
					Content:   input,
					Role:      schema.UserMsg,
					Timestamp: time.Now().Unix(),
				}

				chat.Msgs = append(chat.Msgs, userMsg)

				aiMsg := schema.Msg{
					Stream:    true,
					Content:   "",
					Role:      schema.AIMsg,
					Timestamp: time.Now().Unix(),
				}

				chat.Msgs = append(chat.Msgs, aiMsg)

				chat.Viewport.SetContent(chat.RenderMsgs())
				chat.Viewport.GotoBottom()

				go func() {
					err := chat.Provider.Run(context.TODO(), input, chat.Session)
					if err != nil {
						log.Printf("Error: %v", err)
					}
				}()
			}

			return chat, chat.Loader.Tick
		}
	}

	input, cmd := chat.Input.Update(msg)
	chat.Input = &input
	cmds = append(cmds, cmd)

	vp, cmd := chat.Viewport.Update(msg)
	chat.Viewport = &vp
	cmds = append(cmds, cmd)

	return chat, tea.Batch(cmds...)
}

func (chat ChatView) HandleStream() tea.Msg {
	return <-chat.Stream
}
