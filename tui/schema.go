package tui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
)

var Providers = []ChatProvider{}

// var Commands = []list.Item{}

type ChatCmd struct {
	title, desc string
	exe         func(LayoutView) (LayoutView, tea.Cmd)
}

func (item ChatCmd) Title() string                              { return item.title }
func (item ChatCmd) Description() string                        { return item.desc }
func (item ChatCmd) FilterValue() string                        { return item.title }
func (item ChatCmd) Execute(l LayoutView) (LayoutView, tea.Cmd) { return item.exe(l) }

type ChatMsg struct {
	Role      string
	Content   string
	Timestamp int64
}

type ChatHistory interface {
	LoadHistory(sessionID string) []ChatMsg
	SaveHistory(sessionID string, msgs []ChatMsg) error
}

type ChatProvider interface {
	Name() string
	Type() string
	Description() string
	ChatHistory(sessionID string) ChatHistory
	Run(ctx context.Context, input string) error
	Stream(ctx context.Context, callback func(ctx context.Context, chunk []byte) error)
}

type Session interface {
	NewSession() error
	LoadSession(ID string) ([]ChatMsg, error)
	SaveSession(ID string, msgs []ChatMsg) error
	DeleteSession(ID string) error
}
