package menu

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/struki84/clipt/tui/schema"
)

type MenuDelegate struct {
	Style schema.Styles
}

func NewMenuDelegate(style schema.Styles) MenuDelegate {
	return MenuDelegate{
		Style: style,
	}
}

func (delegate MenuDelegate) Height() int                             { return 1 }
func (delegate MenuDelegate) Spacing() int                            { return 0 }
func (delegate MenuDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (delegate MenuDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(schema.CmdItem)
	if !ok {
		return
	}

	titleStyle := delegate.Style.ChatMenu.TitleNormal

	if index == m.Index() {
		titleStyle = delegate.Style.ChatMenu.TitleSelected
	}

	title := titleStyle.Render(i.Title())
	desc := delegate.Style.ChatMenu.Description.Render(i.Description())

	fmt.Fprint(w, title+desc)
}
