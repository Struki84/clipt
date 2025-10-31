package tui

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var defualtCmds = []list.Item{
	ChatCmd{
		title: "/models",
		desc:  "List available models",
		exe: func(m ChatModel) (ChatModel, tea.Cmd) {
			items := []list.Item{}
			for _, provider := range m.Providers {

				log.Printf("Provider: %s", provider.Name())
				// TODO: The provider type should be some kind of enum type, not a free entry string
				if strings.ToLower(provider.Type()) == "llm" || strings.ToLower(provider.Type()) == "model" {
					provider := provider
					items = append(items,
						ChatCmd{
							title: fmt.Sprintf("/%s", provider.Name()),
							desc:  provider.Description(),
							exe: func(model ChatModel) (ChatModel, tea.Cmd) {
								model.Provider = provider
								model.Layout.Provider = provider
								model.Layout.CurrentMenuItems = model.Layout.MenuItems
								model.Layout.FilteredMenuItems = model.Layout.MenuItems
								model.Layout.ChatInput.SetValue("/")
								return model, nil
							},
						},
					)
				}
			}

			m.Layout.CurrentMenuItems = items
			m.Layout.FilteredMenuItems = items
			m.Layout.ChatInput.SetValue("/")

			return m, nil
		},
	},

	ChatCmd{
		title: "/agents",
		desc:  "List available agents",
		exe: func(m ChatModel) (ChatModel, tea.Cmd) {
			items := []list.Item{}
			for _, provider := range m.Providers {
				if strings.ToLower(provider.Type()) == "agent" {
					provider := provider
					items = append(items, ChatCmd{
						title: fmt.Sprintf("/%s", provider.Name()),
						desc:  provider.Description(),
						exe: func(model ChatModel) (ChatModel, tea.Cmd) {
							model.Provider = provider
							model.Layout.Provider = provider
							model.Layout.CurrentMenuItems = model.Layout.MenuItems
							model.Layout.FilteredMenuItems = model.Layout.MenuItems
							model.Layout.ChatInput.SetValue("/")
							return model, nil
						},
					})
				}
			}

			m.Layout.CurrentMenuItems = items
			m.Layout.FilteredMenuItems = items
			m.Layout.ChatInput.SetValue("/")

			return m, nil
		},
	},

	ChatCmd{
		title: "/sessions",
		desc:  "List session history",
		exe: func(model ChatModel) (ChatModel, tea.Cmd) {
			sessions := model.Storage.ListSessions()
			items := []list.Item{}

			for _, session := range sessions {
				session := session
				items = append(items, ChatCmd{
					title: fmt.Sprintf("/%s", session.Title),
					desc:  fmt.Sprintf("%v", session.CreatedAt),
					exe: func(model ChatModel) (ChatModel, tea.Cmd) {
						model.Layout.Session = session
						model.Layout.Msgs = session.Msgs
						model.Layout.CurrentMenuItems = model.Layout.MenuItems
						model.Layout.FilteredMenuItems = model.Layout.MenuItems
						model.Layout.ChatInput.SetValue("/")
						return model, nil
					},
				})
			}

			model.Layout.CurrentMenuItems = items
			model.Layout.FilteredMenuItems = items
			model.Layout.ChatInput.SetValue("/")

			return model, nil
		},
	},

	ChatCmd{
		title: "/new",
		desc:  "Create new session",
		exe: func(model ChatModel) (ChatModel, tea.Cmd) {
			currentSession, err := model.Storage.NewSession()
			if err != nil {
				log.Printf("Error creating new session")
			}

			model.Layout.Session = currentSession

			return model, nil

		},
	},
	ChatCmd{title: "/exit", desc: "Exit", exe: func(model ChatModel) (ChatModel, tea.Cmd) {
		return model, tea.Quit
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
