package tui

import (
	"fmt"
	"log"
	"os/user"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type LayoutView struct {
	WindowSize tea.WindowSizeMsg

	Provider ChatProvider
	Session  ChatHistory
	Msgs     []ChatMsg

	ChatInput textarea.Model
	ChatView  viewport.Model
	ChatMenu  list.Model

	MenuActive        bool
	MenuItems         []list.Item
	FilteredMenuItems []list.Item

	Style Styles
}

func NewLayoutView(provider ChatProvider) LayoutView {
	return LayoutView{
		Provider:   provider,
		MenuActive: false,
		ChatInput:  textarea.New(),
		ChatView:   viewport.New(120, 35),
		Style:      DefaultStyles(),
	}
}

func (layout LayoutView) Init() tea.Cmd {
	cmds := []tea.Cmd{}

	cmds = append(cmds, textarea.Blink)

	return tea.Batch(cmds...)
}

func (layout LayoutView) View() string {
	// Setup Chat View
	// The session bar and the chat messages(viewport)
	elements := []string{}

	layout.Style.SessionBar.Width(layout.WindowSize.Width - 4)
	sessionBar := layout.Style.SessionBar.Render("New Session \n04 Oct 2025 23:34")

	elements = append(elements, sessionBar)

	layout.ChatView.Width = layout.WindowSize.Width
	layout.ChatView.Height = layout.WindowSize.Height - 9
	layout.ChatView.KeyMap = viewport.KeyMap{
		PageDown: key.NewBinding(key.WithKeys("pgdown")),
		PageUp:   key.NewBinding(key.WithKeys("pgup")),
		Down:     key.NewBinding(key.WithKeys("down")),
		Up:       key.NewBinding(key.WithKeys("up")),
	}

	layout.ChatView.SetContent(layout.RenderMsgs())

	elements = append(elements, layout.ChatView.View())

	//Setup the Chat Input
	//The textarea and the menu list
	if layout.MenuActive {
		menuHeight := len(layout.FilteredMenuItems)
		layout.Style.ChatMenu.Width(layout.WindowSize.Width - 6).Height(menuHeight)

		layout.ChatView.Height = layout.WindowSize.Height - menuHeight - 7
		layout.ChatMenu.SetSize(layout.WindowSize.Width-10, menuHeight)

		menu := layout.Style.ChatMenu.Render(layout.ChatMenu.View())

		elements = append(elements, menu)
	}

	layout.ChatInput.Prompt = ""
	layout.ChatInput.ShowLineNumbers = false
	layout.ChatInput.SetHeight(1)
	layout.ChatInput.SetWidth(layout.WindowSize.Width - 4)
	layout.ChatInput.FocusedStyle.CursorLine = lipgloss.NewStyle()
	layout.ChatInput.FocusedStyle.Base = layout.Style.ChatInput
	layout.ChatInput.Focus()

	elements = append(elements, layout.ChatInput.View())

	//Setup the nvim like status line

	providerType := layout.Style.StatusLine.ProviderType.Render(layout.Provider.Type())
	providerName := layout.Style.StatusLine.ProviderName.Render(layout.Provider.Name())
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

func (layout LayoutView) RenderMsgs() string {
	var styledMessages []string

	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(120),
	)

	width := layout.ChatView.Width - 4

	for _, msg := range layout.Msgs {
		if msg.Role == "User" {
			date := time.Unix(msg.Timestamp, 0).Format("2 Jan 2006 15:04")
			username := user.Username
			fullMsg := fmt.Sprintf("%s\n\n%s (%s) ", msg.Content, username, date)

			layout.Style.Msg.User.Width(width)
			styledMessages = append(styledMessages, layout.Style.Msg.User.Render(fullMsg))

		} else if msg.Role == "AI" {
			renderedTxt, _ := renderer.Render(msg.Content)
			layout.Style.Msg.AI.Width(layout.ChatView.Width - 2)

			styledMessages = append(styledMessages, layout.Style.Msg.AI.Render(renderedTxt))
		}
	}

	return lipgloss.JoinVertical(lipgloss.Center, styledMessages...)
}

func (layout LayoutView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		layout.WindowSize = msg
	}

	return layout, nil
}
