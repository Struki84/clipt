package tui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MenuItem struct {
	title, desc string
}

func (item MenuItem) Title() string       { return item.title }
func (item MenuItem) Description() string { return item.desc }
func (item MenuItem) FilterValue() string { return item.title }

type MenuDelegate struct{}

func (delegate MenuDelegate) Height() int                             { return 1 }
func (delegate MenuDelegate) Spacing() int                            { return 0 }
func (delegate MenuDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (delegate MenuDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	var (
		normalStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#11111b")).
				Foreground(lipgloss.Color("#FFFFFF")).
				Padding(0, 0, 0, 0).
				Width(120)

		selectedStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#11111b")).
				Foreground(lipgloss.Color("#b4befe")).
				Padding(0).
				Width(120)
	)
	i, ok := item.(MenuItem)
	if !ok {
		return
	}
	style := normalStyle

	if index == m.Index() {
		style = selectedStyle
	}

	fmt.Fprint(w, style.Render(string(i.Title())))
}
