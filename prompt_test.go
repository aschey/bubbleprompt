package prompt

import (
	"strings"
	"testing"

	tester "github.com/aschey/tui-tester"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	prompt Model
}

type testCompleterModel struct {
	suggestions []Suggestion
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

func (m completerModel) completer(input string) Suggestions {
	//time.Sleep(100 * time.Millisecond)
	return FilterHasPrefix(input, m.suggestions)
}

func executor(input string, selected *Suggestion, suggestions Suggestions) tea.Model {
	return NewAsyncStringModel(func() string {
		//time.Sleep(10 * time.Millisecond)
		return "test"
	})
}

func Test(t *testing.T) {

	suggestions := []Suggestion{
		{Name: "first option", Description: "test desc", Placeholder: "[hh]"},
		{Name: "second option", Description: "test desc2"},
		{Name: "third option", Description: "test desc2"},
		{Name: "fourth option", Description: "test desc2"},
		{Name: "fifth option", Description: "test desc2"},
	}

	completerModel := completerModel{suggestions: suggestions}

	m := model{prompt: New(
		completerModel.completer,
		executor,
	)}
	m.prompt.ready = true
	m.prompt.viewport = viewport.New(80, 30)

	program := func(tester *tester.Tester) {
		if err := tea.NewProgram(m, tea.WithInput(tester), tea.WithOutput(tester)).Start(); err != nil {
			panic(err)
		}
	}

	tester := tester.New(program)
	tester.Send([]byte("test"))
	tester.WaitFor(func(out string) bool {
		return strings.Contains(out, "test")
	})
	tester.Send([]byte{3})
	tester.WaitForTermination()
}
