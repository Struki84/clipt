package menu

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MenuDelegate struct{}

func NewMenuDelegate() MenuDelegate {
	return MenuDelegate{}
}

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
	i, ok := item.(list.Item)
	if !ok {
		return
	}
	style := normalStyle

	if index == m.Index() {
		style = selectedStyle
	}

	fmt.Fprint(w, style.Render(string(i.Title())))
}

type ChatMenu struct {
	WindowSize tea.WindowSizeMsg
	Style      lipgloss.Style

	List          list.Model
	DefaultItems  []list.Item
	CurrentItems  []list.Item
	FilteredItems []list.Item
	SearchString  string

	MenuActive bool
}

func NewChatMenu(cmds []list.Item) ChatMenu {
	list := list.New(cmds, NewMenuDelegate(), 0, 0)
	return ChatMenu{
		List:          list,
		DefaultItems:  cmds,
		CurrentItems:  cmds,
		FilteredItems: cmds,
	}
}

func (menu ChatMenu) Init() tea.Cmd {
	menu.List.SetShowTitle(false)
	menu.List.SetShowHelp(false)
	menu.List.SetShowPagination(false)
	menu.List.SetShowFilter(false)
	menu.List.SetShowStatusBar(false)
	menu.List.SetFilteringEnabled(false)
	menu.List.KeyMap.CursorDown = key.NewBinding(key.WithKeys("down"))
	menu.List.KeyMap.CursorUp = key.NewBinding(key.WithKeys("up"))

	return nil
}

func (menu ChatMenu) View() string {
	menuHeight := len(menu.FilteredItems)
	if menuHeight == 0 {
		menuHeight = 1
	}

	if menuHeight > 10 {
		menuHeight = 10
	}

	menu.Style = lipgloss.NewStyle().
		Background(lipgloss.Color("#11111b")).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#ffffff")).
		BorderLeft(true).
		BorderRight(true).
		BorderTop(false).
		BorderBottom(false).
		PaddingLeft(1).
		PaddingRight(1).
		Width(menu.WindowSize.Width - 6).
		Height(menuHeight)

	menu.List.SetItems(menu.FilteredItems)
	menu.List.SetSize(menu.WindowSize.Width-4, menuHeight)

	return menu.Style.Render(menu.List.View())
}

func (menu ChatMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		menu.WindowSize = msg
	}

	if menu.MenuActive {
		menu.FilteredItems = []list.Item{}

		for _, item := range menu.CurrentItems {
			if strings.Contains(strings.ToLower(item.FilterValue()), strings.ToLower(menu.SearchString)) {
				menu.FilteredItems = append(menu.FilteredItems, item)
			}
		}

		list, cmd := menu.List.Update(msg)
		menu.List = list
		cmds = append(cmds, cmd)

	} else {
		menu.SearchString = ""
		menu.FilteredItems = menu.DefaultItems
		menu.CurrentItems = menu.DefaultItems
	}

	return menu, tea.Batch(cmds...)
}
