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

func generateSuggestions(maxSeconds int) func() []suggestion.Suggestion[any] {
	return func() []suggestion.Suggestion[any] {
		time.Sleep(time.Duration(rand.Intn(maxSeconds)) * time.Second)
		suggestions := []suggestion.Suggestion[any]{}
		for i := 0; i < rand.Intn(5); i++ {
			suggestions = append(suggestions, suggestion.Suggestion[any]{Text: fmt.Sprintf("suggestion%d", i)})
		}
		return suggestions
	}
}

func (m model) Init() tea.Cmd {
	return suggestion.RefreshSuggestions(generateSuggestions(1))
}

func (m model) Update(msg tea.Msg) (prompt.InputHandler[any], tea.Cmd) {
	switch msg := msg.(type) {
	case suggestion.RefreshSuggestionsMessage[any]:
		m.suggestions = msg
		return m, tea.Batch(suggestion.Complete, suggestion.RefreshSuggestions(generateSuggestions(5)))

	case tea.KeyMsg:
		return m, suggestion.RefreshSuggestions(generateSuggestions(1))
	}

	// if m.initCounter < 7 {
	// 	m.initCounter++
	// 	return m, func() tea.Msg {
	// 		time.Sleep(time.Duration(rand.Intn(m.initCounter)) * time.Second)
	// 		return updateMsg{}
	// 	}
	// }
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
