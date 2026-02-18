package tui

import (
	"context"

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
	filter schema.ProviderType
}

func (cmd ProvidersCmd) Title() string       { return cmd.title }
func (cmd ProvidersCmd) Description() string { return cmd.desc }
func (cmd ProvidersCmd) FilterValue() string { return cmd.title }
func (cmd ProvidersCmd) Execute(m tea.Model) (tea.Model, tea.Cmd) {
	model := m.(ChatModel)
	items := []list.Item{}

	for _, provider := range model.Providers {
		if provider.Type() == cmd.filter {
			items = append(items, ProviderCmd{provider: provider})
		}
	}

	model.Layout.Menu = model.Layout.Menu.PushMenu(items)
	model.Layout.Chat.Input.SetValue("/")

	return model, nil
}

type ProviderCmd struct {
	provider schema.ChatProvider
}

func (cmd ProviderCmd) Title() string       { return "/" + cmd.provider.Name() }
func (cmd ProviderCmd) Description() string { return cmd.provider.Description() }
func (cmd ProviderCmd) FilterValue() string { return cmd.provider.Name() }
func (cmd ProviderCmd) Execute(m tea.Model) (tea.Model, tea.Cmd) {
	model := m.(ChatModel)

	model.Layout.Chat.Provider = cmd.provider
	model.Layout.Chat.Provider.Stream(context.TODO(), func(ctx context.Context, msg schema.Msg) error {
		model.Layout.Chat.Stream <- msg
		return nil
	})

	model.Layout.Menu = model.Layout.Menu.Close()
	model.Layout.Chat.Input.SetValue("")

	return model, nil
}

type SessionsCmd struct {
	title string
	desc  string
}

func (cmd SessionsCmd) Title() string       { return cmd.title }
func (cmd SessionsCmd) Description() string { return cmd.desc }
func (cmd SessionsCmd) FilterValue() string { return cmd.title }
func (cmd SessionsCmd) Execute(m tea.Model) (tea.Model, tea.Cmd) {
	model := m.(ChatModel)
	sessions := model.Storage.ListSessions()
	items := []list.Item{}

	for _, session := range sessions {
		items = append(items, SessionCmd{session: session})
	}

	model.Layout.Menu = model.Layout.Menu.PushMenu(items)
	model.Layout.Chat.Input.SetValue("/")

	return model, nil
}

type SessionCmd struct {
	session schema.ChatSession
}

func (cmd SessionCmd) Title() string       { return "/" + cmd.session.Title }
func (cmd SessionCmd) Description() string { return "" }
func (cmd SessionCmd) FilterValue() string { return cmd.session.Title }
func (cmd SessionCmd) Execute(m tea.Model) (tea.Model, tea.Cmd) {
	model := m.(ChatModel)
	model.Layout.Chat.Session = cmd.session
	model.Layout.Chat.Msgs = cmd.session.Msgs
	model.Layout.Chat.Viewport.SetContent(model.Layout.Chat.RenderMsgs())
	model.Layout.Chat.Viewport.GotoBottom()
	model.Layout.Chat.Input.SetValue("")

	model.Layout.Menu = model.Layout.Menu.Close()

	return model, nil
}

type NewSessionCmd struct {
	title string
	desc  string
}

func (cmd NewSessionCmd) Title() string       { return cmd.title }
func (cmd NewSessionCmd) Description() string { return cmd.desc }
func (cmd NewSessionCmd) FilterValue() string { return cmd.title }
func (cmd NewSessionCmd) Execute(m tea.Model) (tea.Model, tea.Cmd) {
	model := m.(ChatModel)
	session, _ := model.Storage.NewSession()

	model.Layout.Chat.Session = session
	model.Layout.Chat.Msgs = []schema.Msg{}
	model.Layout.Chat.Viewport.SetContent(model.Layout.Chat.RenderMsgs())
	model.Layout.Chat.Viewport.GotoBottom()
	model.Layout.Chat.Input.SetValue("")

	model.Layout.Menu = model.Layout.Menu.Close()

	return model, nil

}

type ExitCmd struct {
	title string
	desc  string
}

func (cmd ExitCmd) Title() string       { return cmd.title }
func (cmd ExitCmd) Description() string { return cmd.desc }
func (cmd ExitCmd) FilterValue() string { return cmd.title }
func (cmd ExitCmd) Execute(model tea.Model) (tea.Model, tea.Cmd) {
	return model, tea.Quit
}

var defaultCmds = []list.Item{
	ProvidersCmd{title: "/models", desc: "List available models", filter: schema.LLM},
	ProvidersCmd{title: "/agents", desc: "List available agents", filter: schema.Agent},
	SessionsCmd{title: "/sessions", desc: "List saved sessions"},
	NewSessionCmd{title: "/new", desc: "Start new session"},
	ExitCmd{title: "/exit", desc: "Close tui chat app"},
}
