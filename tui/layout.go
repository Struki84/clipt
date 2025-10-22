package tui

import (
	"fmt"
	"log"
	"os/user"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/struki84/clipt/tui/menu"
)

type LayoutView struct {
	Style Styles

	WindowSize tea.WindowSizeMsg

	Provider ChatProvider
	Session  ChatSession
	Msgs     []ChatMsg

	ChatInput *textarea.Model
	ChatView  *viewport.Model
	ChatMenu  *list.Model

	Menu menu.ChatMenu

	MenuActive        bool
	MenuItems         []list.Item
	CurrentMenuItems  []list.Item
	FilteredMenuItems []list.Item

	IsLoading bool
	Loader    spinner.Model
}

func NewLayoutView(provider ChatProvider) LayoutView {
	ta := textarea.New()
	ta.Focus()
	vp := viewport.New(0, 0)
	ls := list.New(defualtCmds, MenuDelegate{}, 0, 0)
	l := spinner.New()
	l.Spinner = spinner.Dot

	return LayoutView{
		Menu:              menu.NewChatMenu(defualtCmds),
		Style:             DefaultStyles(),
		Provider:          provider,
		ChatInput:         &ta,
		ChatView:          &vp,
		ChatMenu:          &ls,
		MenuActive:        false,
		MenuItems:         defualtCmds,
		CurrentMenuItems:  defualtCmds,
		FilteredMenuItems: defualtCmds,
		Loader:            l,
		IsLoading:         false,
	}
}

func (layout LayoutView) Init() tea.Cmd {
	log.Printf("Layout Init()")
	cmds := []tea.Cmd{}

	cmds = append(cmds, textarea.Blink)

	return tea.Batch(cmds...)
}

func (layout LayoutView) View() string {
	log.Printf("Layout View()")
	// Setup Chat View
	// The session bar and the chat view(viewport)
	elements := []string{}

	sessionBar := layout.Style.SessionBar.
		Width(layout.WindowSize.Width - 8).
		Render("New Session \n04 Oct 2025 23:34")

	elements = append(elements, sessionBar)

	layout.ChatView.Width = layout.WindowSize.Width - 4
	layout.ChatView.Height = layout.WindowSize.Height - 8
	layout.ChatView.KeyMap = viewport.KeyMap{
		PageDown: key.NewBinding(key.WithKeys("pgdown")),
		PageUp:   key.NewBinding(key.WithKeys("pgup")),
		Down:     key.NewBinding(key.WithKeys("down")),
		Up:       key.NewBinding(key.WithKeys("up")),
	}

	layout.ChatView.SetContent(layout.RenderMsgs())

	//Setup the Chat Input
	//The textarea and the menu list
	if layout.MenuActive {
		menuHeight := len(layout.FilteredMenuItems)
		if menuHeight == 0 {
			menuHeight = 1
		}

		if menuHeight > 10 {
			menuHeight = 10
		}

		layout.ChatView.Height = layout.WindowSize.Height - 8 - menuHeight

		layout.ChatMenu.SetItems(layout.FilteredMenuItems)
		layout.ChatMenu.SetSize(layout.WindowSize.Width-4, menuHeight)
		layout.ChatMenu.SetShowTitle(false)
		layout.ChatMenu.SetShowHelp(false)
		layout.ChatMenu.SetShowPagination(false)
		layout.ChatMenu.SetShowFilter(false)
		layout.ChatMenu.SetShowStatusBar(false)
		layout.ChatMenu.SetFilteringEnabled(false)
		layout.ChatMenu.KeyMap.CursorDown = key.NewBinding(key.WithKeys("down"))
		layout.ChatMenu.KeyMap.CursorUp = key.NewBinding(key.WithKeys("up"))

		menu := layout.Style.ChatMenu.
			Width(layout.WindowSize.Width - 6).
			Height(menuHeight).
			Render(layout.ChatMenu.View())

		elements = append(elements, layout.ChatView.View())
		elements = append(elements, menu)
	} else {
		elements = append(elements, layout.ChatView.View())
	}

	layout.ChatInput.Prompt = ""
	layout.ChatInput.SetHeight(1)
	layout.ChatInput.SetWidth(layout.WindowSize.Width - 4)
	layout.ChatInput.FocusedStyle.CursorLine = lipgloss.NewStyle()
	layout.ChatInput.FocusedStyle.Base = layout.Style.ChatInput
	layout.ChatInput.ShowLineNumbers = false

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

	for i, msg := range layout.Msgs {
		switch msg.Role {
		case "Internal":
			fullMsg := fmt.Sprintf("%s", msg.Content)
			chatMsg := layout.Style.Msg.Internal.Width(width).Render(fullMsg)

			styledMessages = append(styledMessages, chatMsg)
		case "Err":
			fullMsg := fmt.Sprintf("%s", msg.Content)
			chatMsg := layout.Style.Msg.Err.Width(width).Render(fullMsg)

			styledMessages = append(styledMessages, chatMsg)
		case "System":
			fullMsg := fmt.Sprintf("%s", msg.Content)
			chatMsg := layout.Style.Msg.Sys.Width(width).Render(fullMsg)

			styledMessages = append(styledMessages, chatMsg)
		case "User":
			date := time.Unix(msg.Timestamp, 0).Format("2 Jan 2006 15:04")
			username := user.Username
			fullMsg := fmt.Sprintf("%s\n\n%s (%s) ", msg.Content, username, date)
			chatMsg := layout.Style.Msg.User.Width(width).Render(fullMsg)

			styledMessages = append(styledMessages, chatMsg)
		case "AI":
			var renderedTxt string

			if layout.IsLoading && msg.Content == "" && i == len(layout.Msgs)-1 {
				renderedTxt = layout.Loader.View() + " Working..."
			} else {
				renderedTxt, _ = renderer.Render(msg.Content)
			}

			chatMsg := layout.Style.Msg.AI.Width(width).Render(renderedTxt)

			styledMessages = append(styledMessages, chatMsg)
		}
	}

	return lipgloss.JoinVertical(lipgloss.Center, styledMessages...)
}

func (layout LayoutView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Printf("Layout Update()")
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		layout.WindowSize = msg
	case spinner.TickMsg:
		loader, cmd := layout.Loader.Update(msg)
		layout.Loader = loader
		cmds = append(cmds, cmd)
	}

	input, cmd := layout.ChatInput.Update(msg)
	layout.ChatInput = &input
	cmds = append(cmds, cmd)

	chat, cmd := layout.ChatView.Update(msg)
	layout.ChatView = &chat
	cmds = append(cmds, cmd)

	prompt := layout.ChatInput.Value()
	if strings.HasPrefix(prompt, "/") {
		layout.MenuActive = true
		layout.FilteredMenuItems = []list.Item{}
		filterStr := strings.TrimPrefix(prompt, "/")

		for _, item := range layout.CurrentMenuItems {
			if strings.Contains(strings.ToLower(item.FilterValue()), strings.ToLower(filterStr)) {
				layout.FilteredMenuItems = append(layout.FilteredMenuItems, item)
			}
		}

		menu, cmd := layout.ChatMenu.Update(msg)
		layout.ChatMenu = &menu
		cmds = append(cmds, cmd)
	} else {
		layout.MenuActive = false
		layout.CurrentMenuItems = layout.MenuItems
		layout.FilteredMenuItems = layout.MenuItems
	}

	return layout, tea.Batch(cmds...)
}

func (layout LayoutView) ResetMenu() {
	layout.CurrentMenuItems = layout.MenuItems
	layout.FilteredMenuItems = layout.MenuItems
	layout.ChatInput.SetValue("/")
}
