package ui

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
}

func initLayout() layout {
	menuItems := []string{"CHAT", "HISTORY", "SETTINGS"}
	views := map[string]tea.Model{
		"CHAT":     NewChatView(),
		"HISTORY":  NewHistoryView(),
		"SETTINGS": NewSettingsView(),
	}
	return layout{
		menu:       NewMenu(menuItems),
		views:      views,
		activeView: menuItems[0],
	}
}

func (layout layout) Init() tea.Cmd {
	return nil
}

func (layout layout) View() string {
	leftColumnHeight := layout.windowSize.Height - 2

	// Calculate approximate section heights
	menuHeight := leftColumnHeight / 3
	infoHeight := leftColumnHeight / 4
	volumeHeight := leftColumnHeight - menuHeight - infoHeight

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
			Height: msg.Height - 2,
		}

		for i, view := range layout.views {
			view, cmd := view.Update(msg)
			layout.views[i] = view
			cmds = append(cmds, cmd)
		}

		// log.Printf("View: %#v", msg)
		//
		// if view, ok := layout.views[layout.activeView]; ok {
		// 	view, cmd := view.Update(msg)
		// 	layout.views[layout.activeView] = view
		// 	cmds = append(cmds, cmd)
		// }

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return layout, tea.Quit
		case tea.KeyEnter:
			newActiveView := layout.menu.Items[layout.menu.Selected]
			if newActiveView != layout.activeView {
				layout.activeView = newActiveView
			}

			if view, ok := layout.views[layout.activeView]; ok {
				view, cmd := view.Update(msg)
				layout.views[layout.activeView] = view
				cmds = append(cmds, cmd)
			}
		}
	}

	menu, cmd := layout.menu.Update(msg)
	layout.menu = menu
	cmds = append(cmds, cmd)

	return layout, tea.Batch(cmds...)
}

func ShowUI() {
	file, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	p := tea.NewProgram(initLayout(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}