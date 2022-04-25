package input

import "github.com/charmbracelet/lipgloss"

type Datatype int

const (
	Bool Datatype = iota
	String
	Int
	Float
)

type Suggestion[T any] struct {
	Text           string
	CompletionText string
	Description    string
	Metadata       T
	CursorOffset   int
}

func (s Suggestion[T]) Render(selected bool, leftPadding string, maxNameLen int, maxDescLen int, formatters Formatters, scrollbar string) string {
	completionText := s.CompletionText
	if completionText == "" {
		completionText = s.Text
	}
	name := formatters.Name.Format(completionText, maxNameLen, selected)
	description := ""
	if len(s.Description) > 0 {
		description = formatters.Description.Format(s.Description, maxDescLen, selected)
	}

	line := lipgloss.JoinHorizontal(lipgloss.Bottom, leftPadding, name, description, scrollbar)
	return line
}

func (s Suggestion[T]) Key() *string {
	key := s.Text + s.Description
	return &key
}
