package ui

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	leftColumnStyle = lipgloss.NewStyle().
			Width(30)

	sectionStyle = lipgloss.NewStyle().
			Width(leftColumnStyle.GetWidth())
)

type layout struct {
	menu       Menu
	content    ContentView
	windowSize tea.WindowSizeMsg
}

func initLayout() layout {
	return layout{
		menu:    NewMenu([]string{"CHAT", "HISTORY", "SETTINGS"}),
		content: NewContentView(),
	}
}

func (layout layout) Init() tea.Cmd {
	return nil
}

func (layout layout) View() string {
	leftColumnHeight := layout.windowSize.Height - 2

	// Calculate approximate section heights
	menuHeight := leftColumnHeight / 6
	infoHeight := leftColumnHeight / 3
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

	layout.content.Style.Width(layout.windowSize.Width - layout.menu.Style.GetWidth()).
		Height(layout.windowSize.Height - 2)

	return lipgloss.JoinHorizontal(lipgloss.Left, leftColmun, layout.content.View())
}

func (layout layout) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		layout.windowSize = msg
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return layout, tea.Quit
		case "up", "k":
			if layout.menu.Selected > 0 {
				layout.menu.Selected--
			}
		case "down", "j":
			if layout.menu.Selected < len(layout.menu.Items)-1 {
				layout.menu.Selected++
			}
		}
	}
	return layout, nil
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
