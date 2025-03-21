package ui

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/struki84/clipt/internal"
)

type Mode int

const (
	MenuMode Mode = iota
	InputMode
)

var (
	customBorder = lipgloss.Border{
		Left: "█", Right: "",
		Top: "", Bottom: "",
		TopLeft: "", TopRight: "",
		BottomLeft: "", BottomRight: "",
	}

	leftColumnStyle = lipgloss.NewStyle().
			Width(30)

	sectionStyle = lipgloss.NewStyle().
			Width(leftColumnStyle.GetWidth())

	mainColumnStyle = lipgloss.NewStyle().Border(customBorder)
)

type layout struct {
	menu       Menu
	views      map[string]tea.Model
	activeView string
	windowSize tea.WindowSizeMsg
	mode       Mode
}

func ShowUI(agent *internal.Agent) {
	file, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	p := tea.NewProgram(NewLayout(agent), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func NewLayout(agent *internal.Agent) layout {
	menuItems := []string{"CHAT", "MEMORY", "SETTINGS"}
	views := map[string]tea.Model{
		"CHAT":     NewChatView(agent),
		"MEMORY":   NewHistoryView(),
		"SETTINGS": NewSettingsView(),
	}
	return layout{
		menu:       NewMenu(menuItems),
		views:      views,
		activeView: menuItems[0],
		mode:       InputMode,
	}
}

func (layout layout) Init() tea.Cmd {
	cmds := []tea.Cmd{}

	cmds = append(cmds, layout.menu.Init())

	for _, view := range layout.views {
		cmds = append(cmds, view.Init())
	}

	return tea.Batch(cmds...)
}

func (layout layout) View() string {
	leftColumnHeight := layout.windowSize.Height - 2

	// Calculate approximate section heights
	menuHeight := 13
	volumeHeight := 15
	infoHeight := leftColumnHeight - menuHeight - volumeHeight

	layout.menu.Style.Height(menuHeight)

	infoSection := sectionStyle.Height(infoHeight).Render("Info Section")

	vuMeter := NewVUMeter()
	vuMeterSection := sectionStyle.Height(volumeHeight).AlignVertical(lipgloss.Bottom).Render(vuMeter.View())

	leftColmun := lipgloss.JoinVertical(
		lipgloss.Left,
		layout.menu.View(),
		infoSection,
		vuMeterSection,
	)

	if view, ok := layout.views[layout.activeView]; ok {
		mainColumn := mainColumnStyle.Render(view.View())
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			leftColmun,
			mainColumn,
		)
	}

	return leftColmun
}

func (layout layout) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		layout.windowSize = msg
		msg = tea.WindowSizeMsg{
			Width:  msg.Width - leftColumnStyle.GetWidth() - 3,
			Height: msg.Height - 1,
		}

		activeView := layout.views[layout.activeView]
		view, cmd := activeView.Update(msg)
		layout.views[layout.activeView] = view
		cmds = append(cmds, cmd)

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return layout, tea.Quit
		case tea.KeyTab:
			if layout.mode == InputMode {
				layout.mode = MenuMode
				return layout, nil
			}

			if layout.mode == MenuMode {
				layout.mode = InputMode

				newActiveView := layout.menu.Items[layout.menu.Selected]
				if newActiveView != layout.activeView {
					layout.activeView = newActiveView
				}
			}
		case tea.KeyEnter:
			if layout.mode == MenuMode {
				layout.mode = InputMode

				newActiveView := layout.menu.Items[layout.menu.Selected]
				if newActiveView != layout.activeView {
					layout.activeView = newActiveView
				}
			}
		}

		if layout.mode == InputMode {
			activeView := layout.views[layout.activeView]
			view, cmd := activeView.Update(msg)
			layout.views[layout.activeView] = view
			cmds = append(cmds, cmd)
		} else {
			menu, cmd := layout.menu.Update(msg)
			layout.menu = menu
			cmds = append(cmds, cmd)
		}
	default:
		activeView := layout.views[layout.activeView]
		view, cmd := activeView.Update(msg)
		layout.views[layout.activeView] = view
		cmds = append(cmds, cmd)
	}

	return layout, tea.Batch(cmds...)
}
