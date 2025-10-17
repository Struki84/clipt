package clipt

import (
	"context"

	"github.com/charmbracelet/bubbles/list"
)

var Providers = []ChatProvider{}
var Commands = []ChatCmd{}

type ChatMsg struct {
	Role      string
	Content   string
	Timestamp int64
}

type ChatCmd interface {
	list.Item
	Exec() func() error
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
