package lexerinput

import "github.com/charmbracelet/lipgloss"

var DefaultCurrentPlaceholderSuggestion = "240"

// Formatters handles styling for the input.
type Formatters struct {
	// Placeholder handles styling for placeholder that's shown as the user types the current argument.
	Placeholder lipgloss.Style

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
	}
}
