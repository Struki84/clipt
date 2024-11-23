package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	menuStyle = lipgloss.NewStyle().
			Width(30).
			Padding(1, 0)

	buttonStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("7")).
			Foreground(lipgloss.Color("0")).
			Width(28).
			Height(1).
			Align(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Padding(1, 1).
			MarginBottom(1)

	activeButtonStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("0")).
				Foreground(lipgloss.Color("7")).
				Width(28).
				Height(1).
				Align(lipgloss.Center).
				AlignVertical(lipgloss.Center).
				Padding(1, 1).
				MarginBottom(1)
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
	var buttons []string

	for i, btnText := range menu.Items {
		var btn string
		if i == menu.Selected {
			btn = activeButtonStyle.Render(btnText)
		} else {
			btn = buttonStyle.Render(btnText)
		}

		buttons = append(buttons, btn)
	}

	menuView := lipgloss.JoinVertical(lipgloss.Left, buttons...)

	return menu.Style.Render(menuView)
}

func (menu Menu) Update(msg tea.Msg) (Menu, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "k", "up":
			if menu.Selected > 0 {
				menu.Selected--
			}
		case "j", "down":
			if menu.Selected < len(menu.Items)-1 {
				menu.Selected++
			}
		}
	}
	return menu, nil
}
