package prompt

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
	PositionalArgs []PositionalArg
	Flags          []Flag
}

func (s Suggestion) render(selected bool, leftPadding string, maxNameLen int, maxDescLen int, formatters Formatters) string {
	completionText := s.CompletionText
	if completionText == "" {
		completionText = s.Text
	}
	name := formatters.Name.format(completionText, maxNameLen, selected)
	description := ""
	if len(s.Description) > 0 {
		description = formatters.Description.format(s.Description, maxDescLen, selected)
	}

	line := lipgloss.JoinHorizontal(lipgloss.Bottom, leftPadding, name, description)
	return line
}

func (s Suggestion) key() *string {
	key := s.Text + s.Description
	return &key
}

func (s Suggestions) render(paddingSize int, listPosition int, formatters Formatters) []string {
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
		line := cur.render(selected, leftPadding, maxNameLen, maxDescLen, formatters)
		prompts = append(prompts, line)
	}
	return prompts
}
