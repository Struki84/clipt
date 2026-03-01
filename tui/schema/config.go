package schema

import "github.com/charmbracelet/bubbles/list"

type Mode int

const (
	Chat Mode = iota
	Debug
	Execute
)

type Config struct {
	Providers []ChatProvider
	Style     Styles
	Storage   SessionStorage
	Cmds      []list.Item

	Debug struct {
		Log  bool
		Path string
	}
}
