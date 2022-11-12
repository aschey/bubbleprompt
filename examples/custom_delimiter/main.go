package main

import (
	"fmt"
	"os"
	"strings"

	prompt "github.com/aschey/bubbleprompt"
	completers "github.com/aschey/bubbleprompt/completer"
	executors "github.com/aschey/bubbleprompt/executor"
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/input/simpleinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type metadata struct {
	children []input.Suggestion[completers.Metadata]
}

func (m metadata) Children() []input.Suggestion[completers.Metadata] {
	return m.children
}

type model struct {
	promptModel prompt.Model[completers.Metadata]
}

type completerModel struct {
	suggestions []input.Suggestion[completers.Metadata]
	textInput   *simpleinput.Model[completers.Metadata]
	outputStyle lipgloss.Style
}

func (m model) Init() tea.Cmd {
	return m.promptModel.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	p, cmd := m.promptModel.Update(msg)
	m.promptModel = p
	return m, cmd
}

func (m model) View() string {
	return m.promptModel.View()
}

func (m *completerModel) completer(promptModel prompt.Model[completers.Metadata]) ([]input.Suggestion[completers.Metadata], error) {
	return completers.GetRecursiveCompletions(m.textInput.Tokens(), m.textInput.CursorIndex(), m.suggestions), nil
}

func (m *completerModel) executor(input string, selectedSuggestion *input.Suggestion[completers.Metadata]) (tea.Model, error) {
	allValues := strings.Join(m.textInput.TokenValues(), " → ")
	return executors.NewStringModel("You picked: " + m.outputStyle.Render(allValues)), nil
}

func main() {
	textInput := simpleinput.New[completers.Metadata](
		simpleinput.WithDelimiterRegex(`\s*\.\s*`),
		simpleinput.WithTokenRegex(`[^\s\.]+`))
	suggestions := []input.Suggestion[completers.Metadata]{
		{Text: "germany",
			Metadata: metadata{
				children: []input.Suggestion[completers.Metadata]{
					{
						Text: "bavaria",
						Metadata: metadata{
							children: []input.Suggestion[completers.Metadata]{
								{Text: "munich"},
								{Text: "dachau"},
								{Text: "würzburg"},
							},
						},
					},
					{
						Text: "saxony",
						Metadata: metadata{
							children: []input.Suggestion[completers.Metadata]{
								{Text: "leipzig"},
								{Text: "dresden"},
								{Text: "freiberg"},
							},
						},
					},
					{
						Text: "baden-württemberg",
						Metadata: metadata{
							children: []input.Suggestion[completers.Metadata]{
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
				children: []input.Suggestion[completers.Metadata]{
					{
						Text: "ontario",
						Metadata: metadata{
							children: []input.Suggestion[completers.Metadata]{
								{Text: "toronto"},
								{Text: "ottowa"},
								{Text: "windsor"},
							},
						},
					},
					{
						Text: "quebec",
						Metadata: metadata{
							children: []input.Suggestion[completers.Metadata]{
								{Text: "montreal"},
								{Text: "gatineau"},
								{Text: "alma"},
							},
						},
					},
					{
						Text: "alberta",
						Metadata: metadata{
							children: []input.Suggestion[completers.Metadata]{
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
				children: []input.Suggestion[completers.Metadata]{
					{
						Text: "lombardy",
						Metadata: metadata{
							children: []input.Suggestion[completers.Metadata]{
								{Text: "milan"},
								{Text: "brescia"},
								{Text: "varese"},
							},
						},
					},
					{
						Text: "campania",
						Metadata: metadata{
							children: []input.Suggestion[completers.Metadata]{
								{Text: "naples"},
								{Text: "pompeii"},
								{Text: "salerno"},
							},
						},
					},
					{
						Text: "sicily",
						Metadata: metadata{
							children: []input.Suggestion[completers.Metadata]{
								{Text: "palermo"},
								{Text: "catania"},
								{Text: "ragusa"},
							},
						},
					},
				},
			}},
	}

	completerModel := completerModel{
		suggestions: suggestions,
		textInput:   textInput,
		outputStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("13")),
	}

	promptModel, err := prompt.New(
		completerModel.completer,
		completerModel.executor,
		textInput,
	)
	if err != nil {
		panic(err)
	}

	m := model{promptModel}
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render("Pick a fruit!"))
	fmt.Println()
	if _, err := tea.NewProgram(m, tea.WithFilter(prompt.MsgFilter)).Run(); err != nil {
		fmt.Printf("Could not start program :(\n%v\n", err)
		os.Exit(1)
	}
}
