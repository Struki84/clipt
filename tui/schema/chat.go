package schema

import (
	"context"
	"fmt"
)

// Chat schema
const (
	AIMsg MsgRole = iota
	UserMsg
	SysMsg
	ErrMsg
	InternalMsg
)

type MsgRole int

func (r MsgRole) String() string {
	switch r {
	case AIMsg:
		return "AIMsg"
	case UserMsg:
		return "UserMsg"
	case SysMsg:
		return "SysMsg"
	case ErrMsg:
		return "ErrMsg"
	case InternalMsg:
		return "InternalMsg"
	default:
		return fmt.Sprintf("MsgRole(%d)", r)
	}
}

func Enum(s string) MsgRole {
	switch s {
	case "AIMsg":
		return AIMsg
	case "UserMsg":
		return UserMsg
	case "SysMsg":
		return SysMsg
	case "ErrMsg":
		return ErrMsg
	case "InternalMsg":
		return InternalMsg
	default:
		return 0
	}
}

type Msg struct {
	Stream    bool
	Role      MsgRole
	Content   string
	Timestamp int64
}

type SessionStorage interface {
	NewSession() (ChatSession, error)
	ListSessions() []ChatSession
	LoadRecentSession() (ChatSession, error)
	LoadSession(string) (ChatSession, error)
	SaveSession(ChatSession) (ChatSession, error)
	DeleteSession(string) error
}

type ChatSession struct {
	ID        string
	Title     string
	Msgs      []Msg
	CreatedAt int64
}

type ChatProvider interface {
	Name() string
	Type() string
	Description() string
	Run(ctx context.Context, input string, session ChatSession) error
	Stream(ctx context.Context, callback func(ctx context.Context, msg Msg) error)
}
