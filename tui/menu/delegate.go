package menu

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/struki84/clipt/tui/schema"
)

type MenuDelegate struct{}

func NewMenuDelegate() MenuDelegate {
	return MenuDelegate{}
}

func (delegate MenuDelegate) Height() int                             { return 1 }
func (delegate MenuDelegate) Spacing() int                            { return 0 }
func (delegate MenuDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (delegate MenuDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	titleWidth := 30
	var (
		normalStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#11111b")).
				Foreground(lipgloss.Color("#FFFFFF")).
				Padding(0, 0, 0, 0).
				Width(titleWidth)

		selectedStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#11111b")).
				Foreground(lipgloss.Color("#b4befe")).
				Padding(0).
				Width(titleWidth)

		descStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#11111b")).
				Foreground(lipgloss.Color("#6c7086")).
				Width(60)
	)
	i, ok := item.(schema.CmdItem)
	if !ok {
		return
	}

	titleStyle := normalStyle

	if index == m.Index() {
		titleStyle = selectedStyle
	}

	title := titleStyle.Render(i.Title())
	desc := descStyle.Render(i.Description())

	fmt.Fprint(w, title+desc)
}
