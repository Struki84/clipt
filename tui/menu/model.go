package menu

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ChatMenu struct {
	WindowSize tea.WindowSizeMsg
	Style      lipgloss.Style

	List          list.Model
	DefaultItems  []list.Item
	CurrentItems  []list.Item
	FilteredItems []list.Item
	SearchString  string

	Active bool
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

	menu.List.SetShowTitle(false)
	menu.List.SetShowHelp(false)
	menu.List.SetShowPagination(false)
	menu.List.SetShowFilter(false)
	menu.List.SetShowStatusBar(false)
	menu.List.SetFilteringEnabled(false)
	menu.List.KeyMap.CursorDown = key.NewBinding(key.WithKeys("down"))
	menu.List.KeyMap.CursorUp = key.NewBinding(key.WithKeys("up"))

	menu.List.SetItems(menu.FilteredItems)
	menu.List.SetSize(menu.WindowSize.Width-4, menuHeight)

	return menu.Style.Render(menu.List.View())
}

func (menu ChatMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		menu.WindowSize = msg
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if menu.Active {
				return menu.Close(), tea.Batch(cmds...)
			}
		}
	}

	if menu.Active {
		menu.FilteredItems = []list.Item{}

		for _, item := range menu.CurrentItems {
			if strings.Contains(strings.ToLower(item.FilterValue()), strings.ToLower(menu.SearchString)) {
				menu.FilteredItems = append(menu.FilteredItems, item)
			}
		}

		menu.List.SetItems(menu.FilteredItems)

	} else {
		menu.SearchString = ""
		menu.FilteredItems = menu.DefaultItems
		menu.CurrentItems = menu.DefaultItems
		menu.List.SetItems(menu.DefaultItems)
	}
	list, cmd := menu.List.Update(msg)
	menu.List = list
	cmds = append(cmds, cmd)

	return menu, tea.Batch(cmds...)
}

func (menu ChatMenu) Close() ChatMenu {
	menu.Active = false
	menu.SearchString = ""
	menu.CurrentItems = menu.DefaultItems
	menu.FilteredItems = menu.DefaultItems
	menu.List.SetItems(menu.DefaultItems)
	return menu
}

func (menu ChatMenu) PushMenu(submenu []list.Item) ChatMenu {
	menu.FilteredItems = submenu
	menu.CurrentItems = submenu
	menu.List.SetItems(submenu)

	return menu
}
