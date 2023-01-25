package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	prompt "github.com/aschey/bubbleprompt"
	"github.com/aschey/bubbleprompt/input/simpleinput"
	"github.com/aschey/bubbleprompt/suggestion"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	suggestions    []suggestion.Suggestion[any]
	textInput      *simpleinput.Model[any]
	initCounter    int
	numSuggestions int
}

func (m model) Complete(promptModel prompt.Model[any]) ([]suggestion.Suggestion[any], error) {
	time.Sleep(time.Second)
	return m.suggestions, nil
}

func (m model) Execute(input string, promptModel *prompt.Model[any]) (tea.Model, error) {
	return nil, nil
}

type updateMsg struct{}

func (m model) Update(msg tea.Msg) (prompt.InputHandler[any], tea.Cmd) {
	switch msg.(type) {
	case updateMsg:
		m.numSuggestions++
		m.suggestions = []suggestion.Suggestion[any]{}
		for i := 0; i < m.numSuggestions; i++ {
			m.suggestions = append(m.suggestions, suggestion.Suggestion[any]{Text: fmt.Sprintf("suggestion%d", i)})
		}
		return m, suggestion.OneShotCompleter(0)
	case tea.KeyMsg:
		m.numSuggestions = rand.Intn(5)
		return m, func() tea.Msg {
			time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
			return updateMsg{}
		}
	}

	if m.initCounter < 7 {
		m.initCounter++
		return m, func() tea.Msg {
			time.Sleep(time.Duration(rand.Intn(m.initCounter)) * time.Second)
			return updateMsg{}
		}
	}
	return m, nil
}

func main() {
	textInput := simpleinput.New[any]()

	model := model{
		textInput: textInput,
	}

	promptModel := prompt.New[any](
		model,
		textInput,
	)

	fmt.Println(
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("6")).
			Render("Type something and watch the suggestions update asynchronously"),
	)
	fmt.Println()
	if _, err := tea.NewProgram(promptModel, tea.WithFilter(prompt.MsgFilter)).Run(); err != nil {
		fmt.Printf("Could not start program\n%v\n", err)
		os.Exit(1)
	}
}
