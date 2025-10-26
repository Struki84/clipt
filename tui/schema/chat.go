package schema

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
)

type ChatCmd struct {
	title, desc string
	exe         func(tea.Model) (tea.Model, tea.Cmd)
}

func (item ChatCmd) Title() string                            { return item.title }
func (item ChatCmd) Description() string                      { return item.desc }
func (item ChatCmd) FilterValue() string                      { return item.title }
func (item ChatCmd) Execute(m tea.Model) (tea.Model, tea.Cmd) { return item.exe(m) }

type ChatProvider interface {
	Name() string
	Type() string
	Description() string
	Run(ctx context.Context, input string) error
	Stream(ctx context.Context, callback func(ctx context.Context, chunk []byte) error)
}

type SessionStorage interface {
	NewSession() (ChatSession, error)
	ListSessions() []ChatSession
	LoadSession(string) (ChatSession, error)
	SaveSession(ChatSession) (ChatSession, error)
	DeleteSession(string) error
}

type ChatSession struct {
	ID        string
	title     string
	msgs      []ChatMsg
	createdAt int64
}

type ChatMsg struct {
	Role      string
	Content   string
	Timestamp int64
}
