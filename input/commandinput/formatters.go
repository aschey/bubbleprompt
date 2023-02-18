package commandinput

import "github.com/charmbracelet/lipgloss"

var (
	DefaultPlaceholderForeground        = "14"
	DefaultCurrentPlaceholderSuggestion = "240"
	DefaultSelectedTextColor            = "10"
	DefaultFlagForeground               = "245"
	DefaultFlagPlaceholderForeground    = "14"
	DefaultBoolFlagForeground           = "13"
	DefaultNumberFlagForeground         = "5"
)

// PositionalArgFormatter handles styling for positional arguments.
type PositionalArgFormatter struct {
	// Placeholder handles styling for the placeholder that appears before the argument is supplied.
	Placeholder lipgloss.Style
	// Arg handles styling for the argument that is supplied.
	Arg lipgloss.Style
}

// FlagFormatter handles styling for flags.
type FlagFormatter struct {
	// Flag handles styling for the flag itself.
	Flag lipgloss.Style
	// Placeholder handles styling for the placeholder that appears before the flag's argument is supplied (if applicable).
	Placeholder lipgloss.Style
}

// FlagValueFormatter handles styling for different flag value data types.
type FlagValueFormatter struct {
	// String handles styling for string values.
	String lipgloss.Style
	// Bool handles styling for boolean values.
	Bool lipgloss.Style
	// Number handles styling for numeric values.
	Number lipgloss.Style
}

// Formatters handles styling for the command input.
type Formatters struct {
	// PositionalArg handles styling for positional arguments.
	PositionalArg PositionalArgFormatter
	// Flag handles styling for flags.
	Flag FlagFormatter
	// FlagValue handles styling for a flag's value (if applicable).
	FlagValue FlagValueFormatter
	// Placeholder handles styling for placeholder that's shown as the user types the current argument.
	Placeholder lipgloss.Style
	// Prompt handles styling for the prompt that's shown before the user input.
	Prompt lipgloss.Style
	// Command handles styling for the command.
	Command lipgloss.Style
	// SelectedText handles styling for the text that's selected by the suggestion manager.
	SelectedText lipgloss.Style
	// Cursor handles styling for the cursor.
	Cursor lipgloss.Style
}

// DefaultFormatters initializes the [Formatters] with sensible defaults.
// You can modify any settings that you wish after calling this function.
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
