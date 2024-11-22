package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	menuStyle = lipgloss.NewStyle().
		Width(30).
		Padding(1, 0)
)

type Menu struct {
	Selected int
	Items    []string
	Style    lipgloss.Style
}

func NewMenu(items []string) Menu {
	return Menu{
		Items:    items,
		Selected: 0,
		Style:    menuStyle,
	}
}

func (menu Menu) Init() tea.Cmd {
	return nil
}

func (menu Menu) View() string {
	var s string

	for i, item := range menu.Items {
		if i == menu.Selected {
			s += fmt.Sprintf("> %s\n", item)
		} else {
			s += fmt.Sprintf("  %s\n", item)
		}
	}

	return menu.Style.Render(s)
}
