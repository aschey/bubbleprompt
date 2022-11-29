package main

import (
	"fmt"
	"os"
	"strings"

	prompt "github.com/aschey/bubbleprompt"
	"github.com/aschey/bubbleprompt/completer"
	"github.com/aschey/bubbleprompt/executor"
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/input/simpleinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type metadata struct {
	children []input.Suggestion[metadata]
}

func (m metadata) Children() []input.Suggestion[metadata] {
	return m.children
}

type model struct {
	suggestions []input.Suggestion[metadata]
	textInput   *simpleinput.Model[metadata]
	outputStyle lipgloss.Style
}

func (m model) Complete(promptModel prompt.Model[metadata]) ([]input.Suggestion[metadata], error) {
	return completer.GetRecursiveSuggestions(m.textInput.WordTokens(), m.textInput.CursorIndex(), m.suggestions), nil
}

func (m model) Execute(input string, promptModel *prompt.Model[metadata]) (tea.Model, error) {
	allValues := strings.Join(m.textInput.WordTokenValues(), " → ")
	return executor.NewStringModel("You picked: " + m.outputStyle.Render(allValues)), nil
}

func (m model) Update(cmd tea.Msg) (prompt.InputHandler[metadata], tea.Cmd) {
	return m, nil
}

func main() {
	textInput := simpleinput.New(
		simpleinput.WithDelimiterRegex[metadata](`\s*\.\s*`),
		simpleinput.WithTokenRegex[metadata](`[^\s\.]+`))
	suggestions := []input.Suggestion[metadata]{
		{Text: "germany",
			Metadata: metadata{
				children: []input.Suggestion[metadata]{
					{
						Text: "bavaria",
						Metadata: metadata{
							children: []input.Suggestion[metadata]{
								{Text: "munich"},
								{Text: "dachau"},
								{Text: "würzburg"},
							},
						},
					},
					{
						Text: "saxony",
						Metadata: metadata{
							children: []input.Suggestion[metadata]{
								{Text: "leipzig"},
								{Text: "dresden"},
								{Text: "freiberg"},
							},
						},
					},
					{
						Text: "baden-württemberg",
						Metadata: metadata{
							children: []input.Suggestion[metadata]{
								{Text: "stuttgart"},
								{Text: "mannheim"},
								{Text: "heidelberg"},
							},
						},
					},
				},
			}},
		{Text: "canada",
			Metadata: metadata{
				children: []input.Suggestion[metadata]{
					{
						Text: "ontario",
						Metadata: metadata{
							children: []input.Suggestion[metadata]{
								{Text: "toronto"},
								{Text: "ottowa"},
								{Text: "windsor"},
							},
						},
					},
					{
						Text: "quebec",
						Metadata: metadata{
							children: []input.Suggestion[metadata]{
								{Text: "montreal"},
								{Text: "gatineau"},
								{Text: "alma"},
							},
						},
					},
					{
						Text: "alberta",
						Metadata: metadata{
							children: []input.Suggestion[metadata]{
								{Text: "calgary"},
								{Text: "edmonton"},
								{Text: "leduc"},
							},
						},
					},
				},
			}},
		{Text: "italy",
			Metadata: metadata{
				children: []input.Suggestion[metadata]{
					{
						Text: "lombardy",
						Metadata: metadata{
							children: []input.Suggestion[metadata]{
								{Text: "milan"},
								{Text: "brescia"},
								{Text: "varese"},
							},
						},
					},
					{
						Text: "campania",
						Metadata: metadata{
							children: []input.Suggestion[metadata]{
								{Text: "naples"},
								{Text: "pompeii"},
								{Text: "salerno"},
							},
						},
					},
					{
						Text: "sicily",
						Metadata: metadata{
							children: []input.Suggestion[metadata]{
								{Text: "palermo"},
								{Text: "catania"},
								{Text: "ragusa"},
							},
						},
					},
				},
			}},
	}

	model := model{
		suggestions: suggestions,
		textInput:   textInput,
		outputStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("13")),
	}

	promptModel, err := prompt.New[metadata](
		model,
		textInput,
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render("Pick a city!"))
	fmt.Println()
	if _, err := tea.NewProgram(promptModel, tea.WithFilter(prompt.MsgFilter)).Run(); err != nil {
		fmt.Printf("Could not start program\n%v\n", err)
		os.Exit(1)
	}
}
