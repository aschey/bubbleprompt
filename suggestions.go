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
	Name             string
	Placeholder      string
	Datatype         Datatype
	PlaceholderStyle Text
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
	Name           string
	Description    string
	PositionalArgs []PositionalArg
	Flags          []Flag
}

func (s Suggestion) render(selected bool, leftPadding string, maxNameLen int, maxDescLen int, formatters Formatters) string {
	name := formatters.Name.format(s.Name, maxNameLen, selected)
	description := formatters.Description.format(s.Description, maxDescLen, selected)

	line := lipgloss.JoinHorizontal(lipgloss.Bottom, leftPadding, name, description)
	return line
}

func (s Suggestion) key() *string {
	key := s.Name + s.Description
	return &key
}

func (s Suggestions) render(paddingSize int, listPosition int, formatters Formatters) []string {
	maxNameLen := 0
	maxDescLen := 0

	// Determine longest name and description to calculate padding
	for _, cur := range s {
		if len(cur.Name) > maxNameLen {
			maxNameLen = len(cur.Name)
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
