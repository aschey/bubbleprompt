package prompt

import "github.com/charmbracelet/lipgloss"

type Suggestion struct {
	Name        string
	Description string
	Placeholder string
}

func (s Suggestion) render(selected bool, leftPadding string, maxNameLen int, maxDescLen int, formatters Formatters) string {
	name := formatters.Name.format(s.Name, maxNameLen, selected)
	description := formatters.Description.format(s.Description, maxDescLen, selected)

	line := lipgloss.JoinHorizontal(lipgloss.Bottom, leftPadding, name, description)
	return line
}

type Suggestions []Suggestion

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