package suggestion

import (
	"strings"

	"github.com/aschey/bubbleprompt/formatter"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

type Suggestion[T any] struct {
	Text           string
	SuggestionText string
	Description    string
	Metadata       T
	CursorOffset   int
}

func (s Suggestion[T]) GetSuggestionText() string {
	if len(s.SuggestionText) > 0 {
		return s.SuggestionText
	}
	return s.Text
}

func (s Suggestion[T]) Render(
	selected bool,
	maxNameLen int,
	maxDescLen int,
	formatters formatter.Formatters,
	scrollbar string,
	indicator string,
) string {
	name := formatters.Name.Format(s.GetSuggestionText(), maxNameLen, selected)
	selectedIndicator := formatters.SelectedIndicator.Render(indicator)
	if !selected {
		selectedIndicator = strings.Repeat(" ", runewidth.StringWidth(indicator))
	}
	description := ""
	middlePadding := ""
	if maxDescLen > 0 {
		description = formatters.Description.Format(s.Description, maxDescLen, selected)
		if !formatters.Name.HasBackground() && !formatters.Description.HasBackground() {
			middlePadding = " "
		}
	}

	line := lipgloss.JoinHorizontal(
		lipgloss.Bottom,
		selectedIndicator,
		name,
		middlePadding,
		description,
		scrollbar,
	)
	return line
}

func (s Suggestion[T]) Key() *string {
	key := s.Text + s.Description
	return &key
}
