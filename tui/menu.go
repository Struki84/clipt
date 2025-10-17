package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var defualtCmds = []list.Item{
	ChatCmd{
		title: "/models",
		desc:  "List available models",
		exe: func(layout LayoutView) (LayoutView, tea.Cmd) {
			items := []list.Item{}
			for _, provider := range Providers {
				if strings.ToLower(provider.Type()) == "llm" || strings.ToLower(provider.Type()) == "model" {
					provider := provider
					items = append(items,
						ChatCmd{
							title: provider.Name(),
							desc:  provider.Description(),
							exe: func(l LayoutView) (LayoutView, tea.Cmd) {
								l.Provider = provider
								l.CurrentMenuItems = l.MenuItems
								l.FilteredMenuItems = l.MenuItems
								l.ChatInput.SetValue("/")
								return l, nil
							},
						},
					)
				}
			}

			layout.CurrentMenuItems = items
			layout.FilteredMenuItems = items
			layout.ChatInput.SetValue("/")

			return layout, nil
		},
	},

	ChatCmd{
		title: "/agents",
		desc:  "List available agents",
		exe: func(layout LayoutView) (LayoutView, tea.Cmd) {
			items := []list.Item{}
			for _, provider := range Providers {
				if strings.ToLower(provider.Type()) == "agent" {
					provider := provider
					items = append(items, ChatCmd{
						title: provider.Name(),
						desc:  provider.Description(),
						exe: func(l LayoutView) (LayoutView, tea.Cmd) {
							l.Provider = provider
							l.CurrentMenuItems = l.MenuItems
							l.FilteredMenuItems = l.MenuItems
							l.ChatInput.SetValue("/")
							return l, nil
						},
					})
				}
			}

			layout.CurrentMenuItems = items
			layout.FilteredMenuItems = items
			layout.ChatInput.SetValue("/")

			return layout, nil
		},
	},
	// ChatCmd{title: "/sessions", desc: "List session history"},
	// ChatCmd{title: "/new", desc: "Create new session"},
	ChatCmd{title: "/exit", desc: "Exit", exe: func(layout LayoutView) (LayoutView, tea.Cmd) {
		return layout, tea.Quit
	}},
}

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
	i, ok := item.(ChatCmd)
	if !ok {
		return
	}
	style := normalStyle

	if index == m.Index() {
		style = selectedStyle
	}

	fmt.Fprint(w, style.Render(string(i.Title())))
}
