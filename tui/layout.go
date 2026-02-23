package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/struki84/clipt/tui/chat"
	"github.com/struki84/clipt/tui/menu"
	"github.com/struki84/clipt/tui/schema"
)

type LayoutView struct {
	WindowSize tea.WindowSizeMsg

	Style Styles
	Menu  menu.ChatMenu
	Chat  chat.ChatView
}

func NewLayoutView(provider schema.ChatProvider) LayoutView {
	return LayoutView{
		Style: DefaultStyles(),
		Menu:  menu.NewChatMenu(defaultCmds),
		Chat:  chat.New(provider),
	}
}

func (layout LayoutView) Init() tea.Cmd {
	cmds := []tea.Cmd{}

	cmds = append(cmds, textarea.Blink)
	cmds = append(cmds, layout.Chat.Init())

	return tea.Batch(cmds...)
}

func (layout LayoutView) View() string {
	elements := []string{}

	// top bar, session tittle
	date := time.Unix(layout.Chat.Session.CreatedAt, 0).Format("2 Jan 2006")
	title := fmt.Sprintf("%s \n%v", layout.Chat.Session.Title, date)
	sessionBar := layout.Style.SessionBar.
		Width(layout.WindowSize.Width - 6).
		Render(title)

	elements = append(elements, sessionBar)

	layout.Chat.View()

	if layout.Menu.Active {
		menuHeight := len(layout.Menu.FilteredItems)

		layout.Chat.Viewport.Height = layout.WindowSize.Height - menuHeight - 8

		elements = append(elements, layout.Chat.Viewport.View())
		elements = append(elements, layout.Menu.View())
	} else {
		elements = append(elements, layout.Chat.Viewport.View())
	}

	elements = append(elements, layout.Chat.Input.View())

	// the bottom bar, status line
	providerType := layout.Style.StatusLine.ProviderType.Render(layout.Chat.Provider.Type().String())
	providerName := layout.Style.StatusLine.ProviderName.Render(layout.Chat.Provider.Name())
	// modeLabel := layout.Style.StatusLine.Tab.Render("tab")
	tab := layout.Style.StatusLine.Tab.Render("tab")
	mode := layout.Style.StatusLine.Mode.Render("CHAT")

	leftPart := lipgloss.JoinHorizontal(lipgloss.Top, providerType, providerName)
	rightPart := lipgloss.JoinHorizontal(lipgloss.Top, tab, mode)

	fillerWidth := layout.WindowSize.Width - lipgloss.Width(leftPart) - lipgloss.Width(rightPart)
	filler := layout.Style.StatusLine.BaseStyle.Width(fillerWidth).Render("")

	statusLine := lipgloss.JoinHorizontal(lipgloss.Top, leftPart, filler, rightPart)

	elements = append(elements, statusLine)

	return lipgloss.JoinVertical(lipgloss.Center, elements...)
}

func (layout LayoutView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		layout.WindowSize = msg
	}

	prompt := layout.Chat.Input.Value()

	layout.Menu.Active = strings.HasPrefix(prompt, "/")
	if layout.Menu.Active {
		layout.Menu.SearchString = strings.TrimPrefix(prompt, "/")
	}

	menuModel, cmd := layout.Menu.Update(msg)
	layout.Menu = menuModel.(menu.ChatMenu)
	cmds = append(cmds, cmd)

	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.Type == tea.KeyEnter {
		if layout.Menu.Active && len(layout.Menu.FilteredItems) > 0 {
			selected, ok := layout.Menu.List.SelectedItem().(schema.CmdItem)
			if ok && selected != nil {
				cmds = append(cmds, func() tea.Msg {
					return schema.ExecuteCmd{Cmd: selected}
				})
			}
		}
	}

	chatModel, cmd := layout.Chat.Update(msg)
	layout.Chat = chatModel.(chat.ChatView)
	cmds = append(cmds, cmd)

	return layout, tea.Batch(cmds...)
}
