package style

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/struki84/clipt/tui/schema"
)

func Default() (style schema.LayoutStyle) {
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

	style.ContentView = lipgloss.NewStyle().
		Background(lipgloss.Color("#1E1E2E"))

	style.Chat.Header = lipgloss.NewStyle().
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#11111b")).
		PaddingLeft(1).
		PaddingRight(1).
		BorderLeft(true).
		BorderRight(true).
		MarginTop(1)

	style.Menu.ContentView = lipgloss.NewStyle().
		Background(lipgloss.Color("#11111b")).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#ffffff")).
		BorderLeft(true).
		BorderRight(true).
		BorderTop(false).
		BorderBottom(false).
		PaddingLeft(1).
		PaddingRight(1)

	style.Menu.ItemNormal = lipgloss.NewStyle().
		Background(lipgloss.Color("#11111b")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0).
		Width(30)

	style.Menu.ItemSelected = lipgloss.NewStyle().
		Background(lipgloss.Color("#11111b")).
		Foreground(lipgloss.Color("#b4befe")).
		Padding(0).
		Width(30)

	style.Menu.Description = lipgloss.NewStyle().
		Background(lipgloss.Color("#11111b")).
		Foreground(lipgloss.Color("#6c7086")).
		Width(60)

	style.Chat.Input = lipgloss.NewStyle().
		Background(lipgloss.Color("#11111b")).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#ffffff")).
		BorderLeft(true).
		BorderRight(true).
		BorderTop(false).
		BorderBottom(false).
		Padding(1, 1, 0, 1)

	style.StatusLine.BaseStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("#181825")).
		Foreground(lipgloss.Color("#ebdbb2"))

	style.StatusLine.ModeLabel = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7f849c")).
		PaddingRight(1)

	style.StatusLine.ModeName = lipgloss.NewStyle().
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

	style.StatusLine.ProviderType = lipgloss.NewStyle().
		Background(lipgloss.Color("#b4befe")).
		Foreground(lipgloss.Color("#181825")).
		PaddingLeft(1).
		PaddingRight(1).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#b4befe")).
		BorderRight(true)

	style.StatusLine.ProviderName = lipgloss.NewStyle().
		Background(lipgloss.Color("#45475a")).
		Foreground(lipgloss.Color("#fff")).
		PaddingLeft(1).
		PaddingRight(1).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#45475a")).
		BorderLeft(true)

	style.StatusLine.Loader = lipgloss.NewStyle().
		Background(lipgloss.Color("#181825"))

	style.InfoLine = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7f849c")).
		Padding(0, 2, 0, 2).
		MarginBottom(1).
		Align(lipgloss.Left)

	style.Chat.Msg.AI = lipgloss.NewStyle().
		Align(lipgloss.Left)

	style.Chat.Msg.User = lipgloss.NewStyle().
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

	style.Chat.Msg.Sys = lipgloss.NewStyle().
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

	style.Chat.Msg.Err = lipgloss.NewStyle().
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

	style.Chat.Msg.Internal = lipgloss.NewStyle().
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
