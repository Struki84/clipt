package ui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

type VUMeter struct {
	LeftLevel  int
	RightLevel int
	MaxLevel   int
	BarWidth   int
	Character  string
}

func NewVUMeter() VUMeter {
	return VUMeter{
		LeftLevel:  0,
		RightLevel: 0,
		MaxLevel:   15,
		BarWidth:   12,
		Character:  "▇",
	}
}

func (v *VUMeter) SetLevels(left, right int) {
	v.LeftLevel = left
	v.RightLevel = right
}

func (v *VUMeter) renderBar(level int) []string {
	redStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	yellowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("220"))
	greenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("46"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("242"))

	var lines []string
	bar := strings.Repeat(v.Character, v.BarWidth)

	for i := v.MaxLevel - 1; i >= 0; i-- {
		var segment string
		switch {
		case i >= level:
			segment = dimStyle.Render(bar)
		case i >= v.MaxLevel-3:
			segment = redStyle.Render(bar)
		case i >= v.MaxLevel-7:
			segment = yellowStyle.Render(bar)
		default:
			segment = greenStyle.Render(bar)
		}
		lines = append(lines, fmt.Sprintf("│%s│", segment))
	}
	lines = append(lines, "└"+strings.Repeat("─", v.BarWidth)+"┘")
	return lines
}

func (v *VUMeter) View() string {
	leftBar := v.renderBar(v.LeftLevel)
	rightBar := v.renderBar(v.RightLevel)

	var lines []string
	for i := 0; i < len(leftBar); i++ {
		lines = append(lines, fmt.Sprintf("%s%s", leftBar[i], rightBar[i]))
	}

	return strings.Join(lines, "\n")
}
