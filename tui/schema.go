package tui

import (
	"context"
)

var Providers = []ChatProvider{}

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
