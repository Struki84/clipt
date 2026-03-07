package schema

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	View       lipgloss.Style
	ChatHeader lipgloss.Style
	ChatInput  lipgloss.Style
	InfoLine   lipgloss.Style

	ChatMenu struct {
		View          lipgloss.Style
		TitleNormal   lipgloss.Style
		TitleSelected lipgloss.Style
		Description   lipgloss.Style
	}

	StatusLine struct {
		BaseStyle    lipgloss.Style
		Tab          lipgloss.Style
		Mode         lipgloss.Style
		ProviderType lipgloss.Style
		ProviderName lipgloss.Style
		Loader       lipgloss.Style
	}

	Msg struct {
		User     lipgloss.Style
		AI       lipgloss.Style
		Sys      lipgloss.Style
		Err      lipgloss.Style
		Internal lipgloss.Style
	}
}

type LayoutStyle struct {
	ContentView lipgloss.Style
	InfoLine    lipgloss.Style

	StatusLine struct {
		BaseStyle    lipgloss.Style
		ProviderType lipgloss.Style
		ProviderName lipgloss.Style
		Loader       lipgloss.Style
		ModeLabel    lipgloss.Style
		ModeName     lipgloss.Style
	}

	Menu struct {
		ContentView  lipgloss.Style
		ItemNormal   lipgloss.Style
		ItemSelected lipgloss.Style
		Description  lipgloss.Style
	}

	Chat struct {
		ContentView lipgloss.Style
		Header      lipgloss.Style
		Input       lipgloss.Style

		Msg struct {
			User     lipgloss.Style
			AI       lipgloss.Style
			Sys      lipgloss.Style
			Err      lipgloss.Style
			Internal lipgloss.Style
		}
	}
}
