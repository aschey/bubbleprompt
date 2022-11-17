package commandinput

import "github.com/charmbracelet/lipgloss"

var DefaultPlaceholderForeground = "14"
var DefaultCurrentPlaceholderSuggestion = "240"
var DefaultSelectedTextColor = "10"
var DefaultFlagForeground = "245"
var DefaultFlagPlaceholderForeground = "14"
var DefaultBoolFlagForeground = "13"
var DefaultNumberFlagForeground = "5"

type PositionalArgFormatter struct {
	Placeholder lipgloss.Style
	Arg         lipgloss.Style
}

type FlagFormatter struct {
	Flag        lipgloss.Style
	Placeholder lipgloss.Style
}

type FlagValueFormatter struct {
	String lipgloss.Style
	Bool   lipgloss.Style
	Number lipgloss.Style
}

type Formatters struct {
	PositionalArg PositionalArgFormatter
	Flag          FlagFormatter
	FlagValue     FlagValueFormatter
	Placeholder   lipgloss.Style
	Prompt        lipgloss.Style
	Text          lipgloss.Style
	SelectedText  lipgloss.Style
	Cursor        lipgloss.Style
}

func DefaultFormatters() Formatters {
	return Formatters{
		Placeholder: lipgloss.
			NewStyle().
			Foreground(lipgloss.Color(DefaultCurrentPlaceholderSuggestion)),
		SelectedText: lipgloss.
			NewStyle().
			Foreground(lipgloss.Color(DefaultSelectedTextColor)),
		PositionalArg: PositionalArgFormatter{
			Placeholder: lipgloss.
				NewStyle().
				Foreground(lipgloss.Color(DefaultPlaceholderForeground)),
		},
		Flag: FlagFormatter{
			Flag: lipgloss.
				NewStyle().
				Foreground(lipgloss.Color(DefaultFlagForeground)),
			Placeholder: lipgloss.
				NewStyle().
				Foreground(lipgloss.Color(DefaultFlagPlaceholderForeground)),
		},
		FlagValue: FlagValueFormatter{
			Bool: lipgloss.
				NewStyle().
				Foreground(lipgloss.Color(DefaultBoolFlagForeground)),
			Number: lipgloss.
				NewStyle().
				Foreground(lipgloss.Color(DefaultNumberFlagForeground)),
		},
	}
}
