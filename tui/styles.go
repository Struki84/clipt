package tui

import (
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	SessionBar lipgloss.Style
	ChatInput  lipgloss.Style
	ChatMenu   lipgloss.Style

	StatusLine struct {
		BaseStyle    lipgloss.Style
		Tab          lipgloss.Style
		Mode         lipgloss.Style
		ProviderType lipgloss.Style
		ProviderName lipgloss.Style
	}

	Msg struct {
		User     lipgloss.Style
		AI       lipgloss.Style
		Sys      lipgloss.Style
		Err      lipgloss.Style
		Internal lipgloss.Style
	}
}

func DefaultStyles() (s Styles) {
	s.SessionBar = lipgloss.NewStyle().
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#11111b")).
		PaddingLeft(1).
		PaddingRight(1).
		BorderLeft(true).
		BorderRight(true).
		MarginTop(1)

	s.ChatMenu = lipgloss.NewStyle().
		Background(lipgloss.Color("#11111b")).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#ffffff")).
		BorderLeft(true).
		BorderRight(true).
		BorderTop(false).
		BorderBottom(false).
		PaddingLeft(1).
		PaddingRight(1)

	s.ChatInput = lipgloss.NewStyle().
		Background(lipgloss.Color("#11111b")).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#ffffff")).
		BorderLeft(true).
		BorderRight(true).
		BorderTop(false).
		BorderBottom(false).
		Padding(1).
		MarginBottom(1)

	s.StatusLine.BaseStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("#181825")).
		Foreground(lipgloss.Color("#ebdbb2"))

	s.StatusLine.Tab = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7f849c")).
		PaddingRight(1)

	s.StatusLine.Mode = lipgloss.NewStyle().
		Background(lipgloss.Color("#b4befe")).
		Foreground(lipgloss.Color("#181825")).
		PaddingLeft(1).
		PaddingRight(1).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#b4befe")).
		BorderLeft(true).
		BorderRight(false).
		BorderTop(false).
		BorderBottom(false)

	s.StatusLine.ProviderType = lipgloss.NewStyle().
		Background(lipgloss.Color("#b4befe")).
		Foreground(lipgloss.Color("#181825")).
		PaddingLeft(1).
		PaddingRight(1).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#b4befe")).
		BorderRight(true)

	s.StatusLine.ProviderName = lipgloss.NewStyle().
		Background(lipgloss.Color("#45475a")).
		Foreground(lipgloss.Color("#fff")).
		PaddingLeft(1).
		PaddingRight(1).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#45475a")).
		BorderLeft(true)

	s.Msg.AI = lipgloss.NewStyle().
		Align(lipgloss.Left)

	s.Msg.User = lipgloss.NewStyle().
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

	s.Msg.Sys = lipgloss.NewStyle().
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

	s.Msg.Err = lipgloss.NewStyle().
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

	s.Msg.Internal = lipgloss.NewStyle().
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

	return s
}
