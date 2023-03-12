package main

import (
	"fmt"
	"os"
	"strconv"

	prompt "github.com/aschey/bubbleprompt"
	"github.com/aschey/bubbleprompt/completer"
	"github.com/aschey/bubbleprompt/executor"
	"github.com/aschey/bubbleprompt/input/simpleinput"
	"github.com/aschey/bubbleprompt/renderer"
	"github.com/aschey/bubbleprompt/suggestion"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	viewport    viewport.Model
	promptModel prompt.Model[any]
}

func (m model) Init() tea.Cmd {
	return m.promptModel.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height
		promptSize := tea.WindowSizeMsg{
			Width:  msg.Width,
			Height: 7,
		}
		promptModel, newCmd := m.promptModel.Update(promptSize)
		cmd = newCmd
		m.promptModel = promptModel.(prompt.Model[any])
	default:
		promptModel, newCmd := m.promptModel.Update(msg)
		cmd = newCmd
		m.promptModel = promptModel.(prompt.Model[any])
	}
	search := lipgloss.NewStyle().PaddingRight(50).Border(lipgloss.NormalBorder()).Render("Search:")
	bg := lipgloss.JoinVertical(lipgloss.Left, search, "option 1", "option 2", "option 3", "option 4", "option 5", "option 6", "option 7")
	m.viewport.SetContent(bg)

	return m, cmd
}

func (m model) View() string {
	return renderer.PlaceOverlay(9, 1, m.promptModel.View(), m.viewport.View())
}

type promptModel struct {
	suggestions []suggestion.Suggestion[any]
	textInput   *simpleinput.Model[any]
	outputStyle lipgloss.Style
	numChoices  int64
	filterer    completer.Filterer[any]
}

func (m promptModel) Complete(promptModel prompt.Model[any]) ([]suggestion.Suggestion[any], error) {
	if len(m.textInput.Tokens()) > 1 {
		return nil, nil
	}

	return m.filterer.Filter(m.textInput.CurrentTokenBeforeCursor(), m.suggestions), nil
}

func (m promptModel) Execute(input string, promptModel *prompt.Model[any]) (tea.Model, error) {
	tokens := m.textInput.WordTokenValues()
	if len(tokens) == 0 {
		return nil, fmt.Errorf("No selection")
	}
	return executor.NewStringModel(m.formatOutput(tokens[0])), nil
}

func (m promptModel) formatOutput(choice string) string {
	return fmt.Sprintf("You picked: %s\nYou've entered %s submissions(s)\n\n",
		m.outputStyle.Render(choice),
		m.outputStyle.Render(strconv.FormatInt(m.numChoices, 10)))
}

func (m promptModel) Update(msg tea.Msg) (prompt.InputHandler[any], tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok && msg.Type == tea.KeyEnter {
		m.numChoices++
	}
	return m, nil
}

func main() {
	textInput := simpleinput.New(simpleinput.WithPrompt[any](""))
	suggestions := []suggestion.Suggestion[any]{
		{Text: "banana", Description: "good with peanut butter"},
		{Text: "\"sugar apple\"", SuggestionText: "sugar apple", Description: "spherical...ish"},
		{Text: "jackfruit", Description: "the jack of all fruits"},
		{Text: "snozzberry", Description: "tastes like snozzberries"},
		{Text: "lychee", Description: "better than leeches"},
		{Text: "mangosteen", Description: "it's not a mango"},
		{Text: "durian", Description: "stinky"},
	}

	pmodel := promptModel{
		suggestions: suggestions,
		textInput:   textInput,
		outputStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("13")),
		filterer:    completer.NewPrefixFilter[any](),
	}

	ppModel := prompt.New[any](pmodel, textInput,
		prompt.WithViewportRenderer[any](renderer.WithUseHistory(false)))

	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render("Pick a fruit!"))
	fmt.Println()

	if _, err := tea.NewProgram(model{promptModel: ppModel}, tea.WithFilter(prompt.MsgFilter)).Run(); err != nil {
		fmt.Printf("Could not start program\n%v\n", err)
		os.Exit(1)
	}
}
