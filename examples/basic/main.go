package main

import (
	"fmt"
	"os"
	"time"

	prompt "github.com/aschey/bubbleprompt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	prompt prompt.Model
}

type completerModel struct {
	suggestions []prompt.Suggestion
}

func (m model) Init() tea.Cmd {
	return m.prompt.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	p, cmd := m.prompt.Update(msg)
	m.prompt = p
	return m, cmd
}

func (m model) View() string {
	return m.prompt.View()
}

func (m completerModel) completer(input string) []prompt.Suggestion {
	time.Sleep(100 * time.Millisecond)
	return prompt.FilterHasPrefix(input, m.suggestions)
}

func main() {
	suggestions := []prompt.Suggestion{
		{Name: "first option", Description: "test desc"},
		{Name: "second option", Description: "test desc2"},
		{Name: "third option", Description: "test desc2"},
		{Name: "fourth option", Description: "test desc2"},
		{Name: "fifth option", Description: "test desc2"},
	}

	completerModel := completerModel{suggestions: suggestions}

	defaultStyle := lipgloss.
		NewStyle().
		PaddingLeft(1)

	m := model{prompt: prompt.New(
		completerModel.completer,
		prompt.OptionPrompt(">>> "),
		prompt.OptionNameFormatter(func(name string, columnWidth int, selected bool) string {
			foreground := ""
			if selected {
				foreground = "240"
			}
			return defaultStyle.
				Copy().
				PaddingRight(columnWidth - len(name) + 1).
				Foreground(lipgloss.Color(foreground)).
				Background(lipgloss.Color("14")).
				Render(name)
		}),
		prompt.OptionDescriptionFormatter(func(description string, columnWidth int, selected bool) string {
			foreground := ""
			if selected {
				foreground = "240"
			}
			return defaultStyle.
				Copy().
				PaddingRight(columnWidth - len(description) + 1).
				Foreground(lipgloss.Color(foreground)).
				Background(lipgloss.Color("37")).
				Render(description)
		}),
	)}

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Printf("Could not start program :(\n%v\n", err)
		os.Exit(1)
	}
}
