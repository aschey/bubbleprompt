package prompt

import (
	"strings"
	"testing"

	"github.com/MarvinJWendt/testza"
	tuitest "github.com/aschey/tui-tester"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	prompt Model
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
	return FilterHasPrefix(input, m.suggestions)
}

func executor(input string, selected *Suggestion, suggestions Suggestions) tea.Model {
	return NewStringModel("result is " + input)
}

func setup() (Suggestions, tuitest.Tester) {
	suggestions := Suggestions{
		{Name: "first option", Description: "test desc", Placeholder: "[test placeholder]"},
		{Name: "second option", Description: "test desc2"},
		{Name: "third option", Description: "test desc3"},
		{Name: "fourth option", Description: "test desc4"},
		{Name: "fifth option", Description: "test desc5"},
	}

	completerModel := completerModel{suggestions: suggestions}

	m := model{prompt: New(
		completerModel.completer,
		executor,
	)}
	m.prompt.ready = true
	m.prompt.viewport = viewport.New(80, 30)

	program := func(tester *tuitest.Tester) {
		if err := tea.NewProgram(m, tea.WithInput(tester), tea.WithOutput(tester)).Start(); err != nil {
			panic(err)
		}
	}

	tester := tuitest.New(program)
	tester.TrimOutput = true
	tester.RemoveAnsi = true
	return suggestions, tester
}

func waitFor(t *testing.T, tester tuitest.Tester, condition func(out string, outputLines []string) bool) []string {
	_, lines, err := tester.WaitFor(condition)
	testza.AssertNoError(t, err)
	return lines
}

func TestBasicCompleter(t *testing.T) {
	suggestions, tester := setup()

	lines := waitFor(t, tester, func(out string, outputLines []string) bool {
		return len(outputLines) > 1
	})

	for i := 1; i < len(suggestions); i++ {
		testza.AssertContains(t, lines[i], suggestions[i-1].Name)
		testza.AssertContains(t, lines[i], suggestions[i-1].Description)
	}

	tester.SendByte(tuitest.KeyCtrlC)

	testza.AssertNoError(t, tester.WaitForTermination())
}

func TestFilter(t *testing.T) {
	suggestions, tester := setup()

	tester.SendString("fi")
	lines := waitFor(t, tester, func(out string, outputLines []string) bool {
		return len(outputLines) > 1
	})
	testza.AssertEqual(t, 3, len(lines))
	testza.AssertContains(t, lines[0], "fi")
	testza.AssertContains(t, lines[1], suggestions[0].Name)
	testza.AssertContains(t, lines[1], suggestions[0].Description)
	testza.AssertContains(t, lines[2], suggestions[4].Name)
	testza.AssertContains(t, lines[2], suggestions[4].Description)
}

func testExecutor(t *testing.T, in *string, expectedOut string) {
	_, tester := setup()

	if in != nil {
		tester.SendString("fi")
	}

	_ = waitFor(t, tester, func(out string, outputLines []string) bool {
		return len(outputLines) > 1
	})
	tester.SendByte(tuitest.KeyEnter)
	_ = waitFor(t, tester, func(out string, outputLines []string) bool {
		return len(outputLines) > 1 && strings.Contains(outputLines[1], expectedOut)
	})
}

func TestExecutorNoInput(t *testing.T) {
	testExecutor(t, nil, "result is")
}

func TestExecutorWithInput(t *testing.T) {
	in := "fi"
	testExecutor(t, &in, "result is fi")
}
