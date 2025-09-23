package v2

import "context"

type AIEngine interface {
	Name() string
	Type() string
	Description() string
	GetChatHistory() []string
	Run(ctx context.Context, input string) error
	Stream(ctx context.Context, callback func(ctx context.Context, chunk []byte))
}
