package style

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/struki84/clipt/tui/schema"
)

func Catpuccin() (style schema.LayoutStyle) {
	// Colorscheme
	const (
		// reused background and foreground colors
		primaryBGcolor   = "#1E1E2E"
		secondaryBGcolor = "#11111b"
		tertiaryBGcolor  = "#181825"

		primaryFGcolor   = "#b4befe"
		secondaryFGcolor = "#ffffff"
		tertiaryFGcolor  = "#7f849c"

		// unique colors
		statusLineFGcolor            = "#ebdbb2"
		providerNameBGcolor          = "#45475a"
		menuDescFGcolor              = "#6c7086"
		chatMsgErrBorderFGcolor      = "#e64553"
		chatMsgInternalBorderFGcolor = "#fab387"
	)

	// Main container view
	style.ContentView = lipgloss.NewStyle().
		Background(lipgloss.Color(primaryBGcolor))

	// Infoline and status line
	style.InfoLine = lipgloss.NewStyle().
		Foreground(lipgloss.Color(tertiaryFGcolor)).
		Padding(0, 2, 0, 2).
		MarginBottom(1).
		Align(lipgloss.Left)

	style.StatusLine.BaseStyle = lipgloss.NewStyle().
		Background(lipgloss.Color(tertiaryBGcolor)).
		Foreground(lipgloss.Color(statusLineFGcolor))

	style.StatusLine.ModeLabel = lipgloss.NewStyle().
		Foreground(lipgloss.Color(tertiaryFGcolor)).
		PaddingRight(1)

	style.StatusLine.ModeName = lipgloss.NewStyle().
		Background(lipgloss.Color(primaryFGcolor)).
		Foreground(lipgloss.Color(tertiaryBGcolor)).
		PaddingLeft(1).
		PaddingRight(1).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color(primaryFGcolor)).
		BorderLeft(true).
		BorderRight(false).
		BorderTop(false).
		BorderBottom(false)

	style.StatusLine.ProviderType = lipgloss.NewStyle().
		Background(lipgloss.Color(primaryFGcolor)).
		Foreground(lipgloss.Color(tertiaryBGcolor)).
		PaddingLeft(1).
		PaddingRight(1).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color(primaryFGcolor)).
		BorderRight(true)

	style.StatusLine.ProviderName = lipgloss.NewStyle().
		Background(lipgloss.Color(providerNameBGcolor)).
		Foreground(lipgloss.Color(secondaryFGcolor)).
		PaddingLeft(1).
		PaddingRight(1).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color(tertiaryFGcolor)).
		BorderLeft(true)

	style.StatusLine.Loader = lipgloss.NewStyle().
		Background(lipgloss.Color((tertiaryBGcolor)))

	// Chat menu
	style.Menu.ContentView = lipgloss.NewStyle().
		Background(lipgloss.Color(secondaryBGcolor)).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color(secondaryFGcolor)).
		BorderLeft(true).
		BorderRight(true).
		BorderTop(false).
		BorderBottom(false).
		PaddingLeft(1).
		PaddingRight(1)

	style.Menu.ItemNormal = lipgloss.NewStyle().
		Background(lipgloss.Color(secondaryBGcolor)).
		Foreground(lipgloss.Color(secondaryFGcolor)).
		Padding(0).
		Width(30)

	style.Menu.ItemSelected = lipgloss.NewStyle().
		Background(lipgloss.Color(secondaryBGcolor)).
		Foreground(lipgloss.Color(primaryFGcolor)).
		Padding(0).
		Width(30)

	style.Menu.Description = lipgloss.NewStyle().
		Background(lipgloss.Color(secondaryBGcolor)).
		Foreground(lipgloss.Color(menuDescFGcolor)).
		Width(60)

	// Chat view - viewport, input, and messages
	style.Chat.Header = lipgloss.NewStyle().
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color(secondaryBGcolor)).
		PaddingLeft(1).
		PaddingRight(1).
		BorderLeft(true).
		BorderRight(true).
		MarginTop(1)

	style.Chat.ContentView = lipgloss.NewStyle().
		Background(lipgloss.Color(primaryBGcolor))

	style.Chat.Msg.AI = lipgloss.NewStyle().
		Align(lipgloss.Left)

	style.Chat.Msg.User = lipgloss.NewStyle().
		Background(lipgloss.Color(tertiaryBGcolor)).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color(primaryFGcolor)).
		BorderLeft(true).
		BorderRight(true).
		BorderTop(false).
		BorderBottom(false).
		Padding(1).
		Margin(1).
		Align(lipgloss.Left)

	style.Chat.Msg.Sys = lipgloss.NewStyle().
		Background(lipgloss.Color(tertiaryBGcolor)).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color(tertiaryBGcolor)).
		BorderLeft(true).
		BorderRight(true).
		BorderTop(false).
		BorderBottom(false).
		Padding(1).
		Margin(1).
		Align(lipgloss.Left)

	style.Chat.Msg.Err = lipgloss.NewStyle().
		Background(lipgloss.Color(tertiaryBGcolor)).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color(chatMsgErrBorderFGcolor)).
		BorderLeft(true).
		BorderRight(true).
		BorderTop(false).
		BorderBottom(false).
		Padding(1).
		Margin(1).
		Align(lipgloss.Left)

	style.Chat.Msg.Internal = lipgloss.NewStyle().
		Background(lipgloss.Color(tertiaryBGcolor)).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color(chatMsgInternalBorderFGcolor)).
		BorderLeft(true).
		BorderRight(true).
		BorderTop(false).
		BorderBottom(false).
		Padding(1).
		Margin(1).
		Align(lipgloss.Left)

	style.Chat.Input = lipgloss.NewStyle().
		Background(lipgloss.Color(secondaryBGcolor)).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color(secondaryFGcolor)).
		BorderLeft(true).
		BorderRight(true).
		BorderTop(false).
		BorderBottom(false).
		Padding(1, 1, 0, 1)

	return style
}
