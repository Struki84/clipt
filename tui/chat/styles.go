package chat

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	Input lipgloss.Style
	Msg   struct {
		User     lipgloss.Style
		AI       lipgloss.Style
		Sys      lipgloss.Style
		Err      lipgloss.Style
		Internal lipgloss.Style
	}
}

func DefaultStyles() (style Styles) {
	style.Input = lipgloss.NewStyle().
		Background(lipgloss.Color("#11111b")).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#ffffff")).
		BorderLeft(true).
		BorderRight(true).
		BorderTop(false).
		BorderBottom(false).
		Padding(1).
		MarginBottom(1)

	style.Msg.AI = lipgloss.NewStyle().
		Align(lipgloss.Left)

	style.Msg.User = lipgloss.NewStyle().
		Background(lipgloss.Color("#181825")).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#b4befe")).
		BorderLeft(true).
		BorderRight(true).
		BorderTop(false).
		BorderBottom(false).
		Padding(1).
		Margin(1).
		Align(lipgloss.Left)

	style.Msg.Sys = lipgloss.NewStyle().
		Background(lipgloss.Color("#181825")).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#181825")).
		BorderLeft(true).
		BorderRight(true).
		BorderTop(false).
		BorderBottom(false).
		Padding(1).
		Margin(1).
		Align(lipgloss.Left)

	style.Msg.Err = lipgloss.NewStyle().
		Background(lipgloss.Color("#181825")).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#e64553")).
		BorderLeft(true).
		BorderRight(true).
		BorderTop(false).
		BorderBottom(false).
		Padding(1).
		Margin(1).
		Align(lipgloss.Left)

	style.Msg.Internal = lipgloss.NewStyle().
		Background(lipgloss.Color("#181825")).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#fab387")).
		BorderLeft(true).
		BorderRight(true).
		BorderTop(false).
		BorderBottom(false).
		Padding(1).
		Margin(1).
		Align(lipgloss.Left)
	return style
}
