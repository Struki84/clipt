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

	layout.content.Style.Width(layout.windowSize.Width - leftColumnStyle.GetWidth() - 3).
		Height(layout.windowSize.Height - 2)

	return lipgloss.JoinHorizontal(lipgloss.Left, leftColmun, layout.content.View())
}

func (layout layout) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		layout.windowSize = msg
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return layout, tea.Quit
		}
	}

	menu, cmd := layout.menu.Update(msg)
	layout.menu = menu
	cmds = append(cmds, cmd)

	content, cmd := layout.content.Update(msg)
	layout.content = content
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
