package tui

import (
	"context"
	"time"

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
func (cmd ProvidersCmd) Execute(model tea.Model) (tea.Model, tea.Cmd) {
	layout := model.(LayoutView)
	items := []list.Item{}

	for _, provider := range layout.Providers {
		if provider.Type() == cmd.filter {
			items = append(items, ProviderCmd{provider: provider})
		}
	}

	layout.Menu = layout.Menu.PushMenu(items)
	layout.Chat.Input.SetValue("/")

	return layout, nil
}

type ProviderCmd struct {
	provider schema.ChatProvider
}

func (cmd ProviderCmd) Title() string       { return "/" + cmd.provider.Name() }
func (cmd ProviderCmd) Description() string { return cmd.provider.Description() }
func (cmd ProviderCmd) FilterValue() string { return cmd.provider.Name() }
func (cmd ProviderCmd) Execute(model tea.Model) (tea.Model, tea.Cmd) {
	layout := model.(LayoutView)

	layout.Chat.Provider = cmd.provider
	layout.Chat.Provider.Stream(context.TODO(), func(ctx context.Context, msg schema.Msg) error {
		layout.Chat.Stream <- msg
		return nil
	})

	layout.Menu = layout.Menu.Close()
	layout.Chat.Input.SetValue("")

	return layout, nil
}

type SessionsCmd struct {
	title string
	desc  string
}

func (cmd SessionsCmd) Title() string       { return cmd.title }
func (cmd SessionsCmd) Description() string { return cmd.desc }
func (cmd SessionsCmd) FilterValue() string { return cmd.title }
func (cmd SessionsCmd) Execute(model tea.Model) (tea.Model, tea.Cmd) {
	layout := model.(LayoutView)
	sessions := layout.Storage.ListSessions()
	items := []list.Item{}

	for _, session := range sessions {
		items = append(items, SessionCmd{session: session})
	}

	layout.Menu = layout.Menu.PushMenu(items)
	layout.Chat.Input.SetValue("/")

	return layout, nil
}

type SessionCmd struct {
	session schema.ChatSession
}

func (cmd SessionCmd) Title() string { return "/" + cmd.session.Title }
func (cmd SessionCmd) Description() string {
	return time.Unix(cmd.session.CreatedAt, 0).Format("2 Jan 2006")
}
func (cmd SessionCmd) FilterValue() string { return cmd.session.Title }
func (cmd SessionCmd) Execute(model tea.Model) (tea.Model, tea.Cmd) {
	layout := model.(LayoutView)
	layout.Chat.Session = cmd.session
	layout.Chat.Msgs = cmd.session.Msgs
	layout.Chat.Viewport.SetContent(layout.Chat.RenderMsgs())
	layout.Chat.Viewport.GotoBottom()
	layout.Chat.Input.SetValue("")

	layout.Menu = layout.Menu.Close()

	return layout, nil
}

type NewSessionCmd struct {
	title string
	desc  string
}

func (cmd NewSessionCmd) Title() string       { return cmd.title }
func (cmd NewSessionCmd) Description() string { return cmd.desc }
func (cmd NewSessionCmd) FilterValue() string { return cmd.title }
func (cmd NewSessionCmd) Execute(model tea.Model) (tea.Model, tea.Cmd) {
	layout := model.(LayoutView)
	session, _ := layout.Storage.NewSession()

	layout.Chat.Session = session
	layout.Chat.Msgs = []schema.Msg{}
	layout.Chat.Viewport.SetContent(layout.Chat.RenderMsgs())
	layout.Chat.Viewport.GotoBottom()
	layout.Chat.Input.SetValue("")

	layout.Menu = layout.Menu.Close()

	return layout, nil

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

var DefaultCmds = []list.Item{
	ProvidersCmd{title: "/models", desc: "List available models", filter: schema.LLM},
	ProvidersCmd{title: "/agents", desc: "List available agents", filter: schema.Agent},
	SessionsCmd{title: "/sessions", desc: "List saved sessions"},
	NewSessionCmd{title: "/new", desc: "Start new session"},
	ExitCmd{title: "/exit", desc: "Close tui chat app"},
}
