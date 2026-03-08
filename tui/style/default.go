package style

import (
	"github.com/charmbracelet/glamour/styles"
	"github.com/charmbracelet/lipgloss"
	"github.com/struki84/clipt/tui/schema"
)

func Default() (style schema.LayoutStyle) {
	// Colorscheme
	var (
		// reused background and foreground colors
		primaryBGcolor   = "#FAF4ED" // Base
		secondaryBGcolor = "#FFFAF3" // Surface
		tertiaryBGcolor  = "#F2E9E1" // Overlay

		primaryFGcolor   = "#907AA9" // Iris
		secondaryFGcolor = "#575279" // Text
		tertiaryFGcolor  = "#9893A5" // Muted

		statusLineFGcolor            = "#6E6A86" // Subtle
		providerNameBGcolor          = "#F2E9E1" // Overlay
		menuDescFGcolor              = "#9893A5" // Muted
		chatMsgErrBorderFGcolor      = "#B4637A" // Love
		chatMsgInternalBorderFGcolor = "#D7827E" // Rose
	)

	//Background color for adding, margin, and border chars
	style.WhitespaceBGcolor = primaryBGcolor

	// Main container view
	style.ContentView = lipgloss.NewStyle().
		Background(lipgloss.Color(primaryBGcolor))

	// Infoline and status line
	style.InfoLine = lipgloss.NewStyle().
		Background(lipgloss.Color(primaryBGcolor)).
		Foreground(lipgloss.Color(tertiaryFGcolor)).
		Padding(0, 2, 0, 2).
		MarginBottom(1).
		Align(lipgloss.Left)

	style.StatusLine.BaseStyle = lipgloss.NewStyle().
		Background(lipgloss.Color(tertiaryBGcolor)).
		Foreground(lipgloss.Color(statusLineFGcolor))

	style.StatusLine.ModeLabel = lipgloss.NewStyle().
		Background(lipgloss.Color(tertiaryBGcolor)).
		Foreground(lipgloss.Color(tertiaryFGcolor)).
		PaddingRight(1)

	style.StatusLine.ModeName = lipgloss.NewStyle().
		Background(lipgloss.Color(primaryFGcolor)).
		Foreground(lipgloss.Color(tertiaryBGcolor)).
		PaddingLeft(1).
		PaddingRight(1).
		BorderStyle(lipgloss.ThickBorder()).
		BorderBackground(lipgloss.Color(tertiaryBGcolor)).
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
		BorderBackground(lipgloss.Color(tertiaryBGcolor)).
		BorderForeground(lipgloss.Color(primaryFGcolor)).
		BorderRight(true)

	style.StatusLine.ProviderName = lipgloss.NewStyle().
		Background(lipgloss.Color(providerNameBGcolor)).
		Foreground(lipgloss.Color(secondaryFGcolor)).
		PaddingLeft(1).
		PaddingRight(1).
		BorderStyle(lipgloss.ThickBorder()).
		BorderBackground(lipgloss.Color(tertiaryBGcolor)).
		BorderForeground(lipgloss.Color(providerNameBGcolor)).
		BorderLeft(true)

	style.StatusLine.Loader = lipgloss.NewStyle().
		Background(lipgloss.Color((tertiaryBGcolor)))

	// Chat menu
	style.Menu.ContentView = lipgloss.NewStyle().
		Background(lipgloss.Color(secondaryBGcolor)).
		BorderStyle(lipgloss.ThickBorder()).
		BorderBackground(lipgloss.Color(secondaryBGcolor)).
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
		Background(lipgloss.Color(primaryBGcolor)).
		Foreground(lipgloss.Color(secondaryFGcolor)).
		BorderStyle(lipgloss.ThickBorder()).
		BorderBackground(lipgloss.Color(primaryBGcolor)).
		BorderForeground(lipgloss.Color(secondaryBGcolor)).
		BorderRight(true).
		PaddingLeft(1).
		PaddingRight(1).
		BorderLeft(true).
		MarginTop(1).
		MarginBackground(lipgloss.Color(primaryBGcolor))

	style.Chat.ContentView = lipgloss.NewStyle().
		Background(lipgloss.Color(primaryBGcolor))

	style.Chat.Msg.AI = lipgloss.NewStyle().
		Background(lipgloss.Color(primaryBGcolor)).
		MarginLeft(3).
		MarginRight(3).
		MarginBackground(lipgloss.Color(primaryBGcolor))

	style.Chat.Msg.User = lipgloss.NewStyle().
		Background(lipgloss.Color(tertiaryBGcolor)).
		Foreground(lipgloss.Color(secondaryFGcolor)).
		BorderStyle(lipgloss.ThickBorder()).
		BorderBackground(lipgloss.Color(primaryBGcolor)).
		BorderForeground(lipgloss.Color(primaryFGcolor)).
		BorderLeft(true).
		BorderRight(true).
		BorderTop(false).
		BorderBottom(false).
		Padding(1).
		Margin(1, 2, 1, 2).
		MarginBackground(lipgloss.Color(primaryBGcolor)).
		Align(lipgloss.Left)

	style.Chat.Msg.Sys = lipgloss.NewStyle().
		Background(lipgloss.Color(tertiaryBGcolor)).
		Foreground(lipgloss.Color(secondaryFGcolor)).
		BorderBackground(lipgloss.Color(primaryBGcolor)).
		BorderStyle(lipgloss.ThickBorder()).
		BorderBackground(lipgloss.Color(primaryBGcolor)).
		BorderForeground(lipgloss.Color(tertiaryBGcolor)).
		BorderLeft(true).
		BorderRight(true).
		BorderTop(false).
		BorderBottom(false).
		Padding(1).
		Margin(1, 2, 1, 2).
		Align(lipgloss.Left)

	style.Chat.Msg.Err = lipgloss.NewStyle().
		Background(lipgloss.Color(tertiaryBGcolor)).
		Foreground(lipgloss.Color(secondaryFGcolor)).
		BorderStyle(lipgloss.ThickBorder()).
		BorderBackground(lipgloss.Color(primaryBGcolor)).
		BorderForeground(lipgloss.Color(chatMsgErrBorderFGcolor)).
		BorderLeft(true).
		BorderRight(true).
		BorderTop(false).
		BorderBottom(false).
		Padding(1).
		Margin(1, 2, 1, 2).
		Align(lipgloss.Left)

	style.Chat.Msg.Internal = lipgloss.NewStyle().
		Background(lipgloss.Color(tertiaryBGcolor)).
		Foreground(lipgloss.Color(secondaryFGcolor)).
		BorderStyle(lipgloss.ThickBorder()).
		BorderBackground(lipgloss.Color(primaryBGcolor)).
		BorderForeground(lipgloss.Color(chatMsgInternalBorderFGcolor)).
		BorderLeft(true).
		BorderRight(true).
		BorderTop(false).
		BorderBottom(false).
		Padding(1).
		Margin(1, 2, 1, 2).
		MarginBackground(lipgloss.Color(primaryBGcolor)).
		Align(lipgloss.Left)

	style.Chat.Input = lipgloss.NewStyle().
		Background(lipgloss.Color(secondaryBGcolor)).
		Foreground(lipgloss.Color(secondaryFGcolor)).
		BorderStyle(lipgloss.ThickBorder()).
		BorderBackground(lipgloss.Color(primaryBGcolor)).
		BorderForeground(lipgloss.Color(secondaryFGcolor)).
		BorderLeft(true).
		BorderRight(true).
		BorderTop(false).
		BorderBottom(false).
		Padding(1, 1, 0, 1)

	// Glamour Styling - WIP
	// Glamour is used for rendering mardkown
	// Match glamour document background
	style.Chat.Msg.Glamour = styles.LightStyleConfig
	// style.Chat.Msg.Glamour.Document.BackgroundColor = &primaryBGcolor
	// style.Chat.Msg.Glamour.CodeBlock.Chroma.Text.BackgroundColor = &primaryBGcolor

	// Remove document margin
	zeroUint := uint(0)
	style.Chat.Msg.Glamour.Document.Margin = &zeroUint
	return style
}
