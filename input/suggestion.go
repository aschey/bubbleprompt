package input

import "github.com/charmbracelet/lipgloss"

type Datatype int

const (
	Bool Datatype = iota
	String
	Int
	Float
)

type PositionalArg struct {
	Placeholder      string
	PlaceholderStyle Text
	ArgStyle         Text
}

type Flag struct {
	Short            string
	Long             string
	Placeholder      string
	Datatype         Datatype
	PlaceholderStyle Text
}

type Suggestions []Suggestion

type Suggestion struct {
	Text           string
	CompletionText string
	Description    string
	Metadata       interface{}
	CursorOffset   int
	PositionalArgs []PositionalArg
	Flags          []Flag
}

func (s Suggestion) render(selected bool, leftPadding string, maxNameLen int, maxDescLen int, formatters Formatters, scrollbar string) string {
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

func (s Suggestion) Key() *string {
	key := s.Text + s.Description
	return &key
}

func (s Suggestions) Render(paddingSize int, listPosition int, formatters Formatters, scrollbar string, scrollbarThumb string) []string {
	maxNameLen := 0
	maxDescLen := 0

	// Determine longest name and description to calculate padding
	for _, cur := range s {
		if len(cur.Text) > maxNameLen {
			maxNameLen = len(cur.Text)
		}
		if len(cur.Description) > maxDescLen {
			maxDescLen = len(cur.Description)
		}
	}

	// Add left offset
	leftPadding := lipgloss.
		NewStyle().
		PaddingLeft(paddingSize).
		Render("")

	prompts := []string{}
	for i, cur := range s {
		selected := i == listPosition
		scrollbarView := scrollbar
		if selected {
			// TODO do actual scrollbar calculation
			scrollbarView = scrollbarThumb
		}
		line := cur.render(selected, leftPadding, maxNameLen, maxDescLen, formatters, scrollbarView)
		prompts = append(prompts, line)
	}
	return prompts
}
