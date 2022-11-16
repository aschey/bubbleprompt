package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	prompt "github.com/aschey/bubbleprompt"
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/input/simpleinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type appModel struct {
	suggestions    []input.Suggestion[any]
	textInput      *simpleinput.Model[any]
	initCounter    int
	numSuggestions int
}

func (m appModel) Complete(promptModel prompt.Model[any]) ([]input.Suggestion[any], error) {
	return m.suggestions, nil
}

func (m appModel) Execute(input string, promptModel *prompt.Model[any]) (tea.Model, error) {
	return nil, nil
}

type updateMsg struct{}

func (m appModel) Update(msg tea.Msg) (prompt.AppModel[any], tea.Cmd) {
	switch msg.(type) {
	case updateMsg:
		m.numSuggestions++
		m.suggestions = []input.Suggestion[any]{}
		for i := 0; i < m.numSuggestions; i++ {
			m.suggestions = append(m.suggestions, input.Suggestion[any]{Text: fmt.Sprintf("suggestion%d", i)})
		}
		return m, prompt.OneShotCompleter(0)
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
	suggestions := []input.Suggestion[any]{
		{Text: "banana", Description: "good with peanut butter"},
		{Text: "\"sugar apple\"", CompletionText: "sugar apple", Description: "spherical...ish"},
		{Text: "jackfruit", Description: "the jack of all fruits"},
		{Text: "snozzberry", Description: "tastes like snozzberries"},
		{Text: "lychee", Description: "better than leeches"},
		{Text: "mangosteen", Description: "it's not a mango"},
		{Text: "durian", Description: "stinky"},
	}

	appModel := appModel{
		suggestions: suggestions,
		textInput:   textInput,
	}

	promptModel, err := prompt.New[any](
		appModel,
		textInput,
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render("Type something and watch the completions update asynchronously"))
	fmt.Println()
	if _, err := tea.NewProgram(promptModel, tea.WithFilter(prompt.MsgFilter)).Run(); err != nil {
		fmt.Printf("Could not start program :(\n%v\n", err)
		os.Exit(1)
	}
}
