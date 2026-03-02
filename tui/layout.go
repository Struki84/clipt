package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/struki84/clipt/tui/chat"
	"github.com/struki84/clipt/tui/menu"
	"github.com/struki84/clipt/tui/schema"
)

type LayoutView struct {
	WindowSize tea.WindowSizeMsg

	Style schema.Styles
	Menu  menu.ChatMenu
	Chat  chat.ChatView

	Storage   schema.SessionStorage
	Providers []schema.ChatProvider
}

func NewLayout(conf schema.Config) LayoutView {
	layout := LayoutView{
		Menu:      menu.New(conf.Cmds, conf.Style),
		Chat:      chat.New(conf.Providers[0], conf.Style),
		Style:     conf.Style,
		Storage:   conf.Storage,
		Providers: conf.Providers,
	}

	if layout.Storage != nil {
		session, err := layout.Storage.LoadRecentSession()
		if err != nil {
			layout.Chat.Msgs = []schema.Msg{}
		} else {
			layout.Chat.Session = session
			layout.Chat.Msgs = session.Msgs
		}
	}

	return layout
}

func (layout LayoutView) Init() tea.Cmd {
	return tea.Batch(layout.Chat.Init())
}

func (layout LayoutView) View() string {
	elements := []string{}

	// Render Chat Header - Session title, date, and info
	// layout.Chat.View() will return only the header section and configure
	// the chat view port and the input, since the layout is dynamic due to
	// open-closing of the menu section
	header := layout.Chat.View()
	elements = append(elements, header)

	// Render Chat viewport and/or chat menu and modify the viewport height based on menu height
	if layout.Menu.Active {
		menuHeight := len(layout.Menu.FilteredItems)
		layout.Chat.Viewport.Height = layout.WindowSize.Height - menuHeight - 8

		elements = append(elements, layout.Chat.Viewport.View())
		elements = append(elements, layout.Menu.View())
	} else {
		elements = append(elements, layout.Chat.Viewport.View())
	}

	// Render Chat input
	elements = append(elements, layout.Chat.Input.View())

	// Render the status line
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
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if layout.Menu.Active {
				layout.Menu = layout.Menu.Close()
				layout.Chat.Input.SetValue("")
				return layout, nil
			}
		case tea.KeyCtrlC:
			return layout, tea.Quit
		}
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
				return selected.Execute(layout)
			}
		}
	}

	chatModel, cmd := layout.Chat.Update(msg)
	layout.Chat = chatModel.(chat.ChatView)
	cmds = append(cmds, cmd)

	return layout, tea.Batch(cmds...)
}
