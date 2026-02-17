package tui

import (
	"context"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/struki84/clipt/tui/schema"
)

type MenuExecuteMsg struct {
	Item schema.CmdItem
}

type ProvidersCmd struct {
	title  string
	desc   string
	filter string
}

func (cmd ProvidersCmd) Title() string       { return cmd.title }
func (cmd ProvidersCmd) Description() string { return cmd.desc }
func (cmd ProvidersCmd) FilterValue() string { return cmd.title }
func (cmd ProvidersCmd) Execute(model tea.Model) (tea.Model, tea.Cmd, []list.Item) {
	chat := model.(ChatModel)
	items := []list.Item{}

	for _, provider := range chat.Providers {
		if strings.ToLower(provider.Type()) == cmd.filter {
			items = append(items, ProviderCmd{provider: provider})
		}
	}

	return chat, nil, items
}

type ProviderCmd struct {
	provider schema.ChatProvider
}

func (cmd ProviderCmd) Title() string       { return "/" + cmd.provider.Name() }
func (cmd ProviderCmd) Description() string { return cmd.provider.Description() }
func (cmd ProviderCmd) FilterValue() string { return cmd.provider.Name() }
func (cmd ProviderCmd) Execute(model tea.Model) (tea.Model, tea.Cmd, []list.Item) {
	chat := model.(ChatModel)

	chat.Layout.Chat.Provider = cmd.provider
	chat.Layout.Chat.Provider.Stream(context.TODO(), func(ctx context.Context, msg schema.Msg) error {
		chat.Layout.Chat.Stream <- msg
		return nil
	})

	return chat, nil, nil
}

type SessionsCmd struct {
	title string
	desc  string
}

func (cmd SessionsCmd) Title() string       { return cmd.title }
func (cmd SessionsCmd) Description() string { return cmd.desc }
func (cmd SessionsCmd) FilterValue() string { return cmd.title }
func (cmd SessionsCmd) Execute(model tea.Model) (tea.Model, tea.Cmd, []list.Item) {
	chat := model.(ChatModel)
	sessions := chat.Storage.ListSessions()
	items := []list.Item{}

	for _, session := range sessions {
		items = append(items, SessionCmd{session: session})
	}

	return chat, nil, items
}

type SessionCmd struct {
	session schema.ChatSession
}

func (cmd SessionCmd) Title() string       { return "/" + cmd.session.Title }
func (cmd SessionCmd) Description() string { return "" }
func (cmd SessionCmd) FilterValue() string { return cmd.session.Title }
func (cmd SessionCmd) Execute(model tea.Model) (tea.Model, tea.Cmd, []list.Item) {
	chat := model.(ChatModel)
	chat.Layout.Chat.Session = cmd.session
	chat.Layout.Chat.Msgs = cmd.session.Msgs
	chat.Layout.Chat.Viewport.SetContent(chat.Layout.Chat.RenderMsgs())
	chat.Layout.Chat.Viewport.GotoBottom()

	return chat, nil, nil
}

type NewSessionCmd struct {
	title string
	desc  string
}

func (cmd NewSessionCmd) Title() string       { return cmd.title }
func (cmd NewSessionCmd) Description() string { return cmd.desc }
func (cmd NewSessionCmd) FilterValue() string { return cmd.title }
func (cmd NewSessionCmd) Execute(model tea.Model) (tea.Model, tea.Cmd, []list.Item) {
	chat := model.(ChatModel)
	session, _ := chat.Storage.NewSession()

	chat.Layout.Chat.Session = session
	chat.Layout.Chat.Msgs = []schema.Msg{}
	chat.Layout.Chat.Viewport.SetContent(chat.Layout.Chat.RenderMsgs())
	chat.Layout.Chat.Viewport.GotoBottom()

	return chat, nil, nil

}

type ExitCmd struct {
	title string
	desc  string
}

func (cmd ExitCmd) Title() string       { return cmd.title }
func (cmd ExitCmd) Description() string { return cmd.desc }
func (cmd ExitCmd) FilterValue() string { return cmd.title }
func (cmd ExitCmd) Execute(model tea.Model) (tea.Model, tea.Cmd, []list.Item) {
	return model, tea.Quit, nil
}

var defaultCmds = []list.Item{
	ProvidersCmd{title: "/models", desc: "List available models", filter: "llm"},
	ProvidersCmd{title: "/agents", desc: "List available agents", filter: "agent"},
	SessionsCmd{title: "/sessions", desc: "List saved sessions"},
	NewSessionCmd{title: "/new", desc: "Start new session"},
	ExitCmd{title: "/exit", desc: "Close tui chat app"},
}
