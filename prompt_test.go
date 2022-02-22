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

type testData struct {
	suggestions  Suggestions
	tester       tuitest.Tester
	initialLines []string
	model        model
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

func (m completerModel) completer(document Document) Suggestions {
	return FilterHasPrefix(document.InputBeforeCursor, m.suggestions)
}

func executor(input string, selected *Suggestion, suggestions Suggestions) tea.Model {
	return NewStringModel("result is " + input)
}

func setup(t *testing.T) testData {
	suggestions := Suggestions{
		{Name: "first option", Description: "test desc", Placeholder: "[test placeholder]"},
		{Name: "second option", Description: "test desc2"},
		{Name: "third option", Description: "test desc3"},
		{Name: "fourth option", Description: "test desc4"},
		{Name: "fifth option", Description: "test desc5"},
	}

	completerModel := completerModel{suggestions: suggestions}

	model := model{prompt: New(
		completerModel.completer,
		executor,
	)}
	model.prompt.ready = true
	model.prompt.viewport = viewport.New(80, 30)

	program := func(tester *tuitest.Tester) {
		if err := tea.NewProgram(model, tea.WithInput(tester), tea.WithOutput(tester)).Start(); err != nil {
			panic(err)
		}
	}

	tester := tuitest.New(program)
	tester.TrimOutput = true
	tester.RemoveAnsi = true

	// Wait for prompt to initialize
	_, initialLines, err := tester.WaitFor(func(out string, outputLines []string) bool {
		return len(outputLines) > 1
	})
	testza.AssertNoError(t, err)

	return testData{suggestions, tester, initialLines, model}
}

func teardown(t *testing.T, tester tuitest.Tester) {
	tester.SendByte(tuitest.KeyCtrlC)
	testza.AssertNoError(t, tester.WaitForTermination())
}

func TestBasicCompleter(t *testing.T) {
	testData := setup(t)

	// Check that all prompts show up
	for i := 1; i < len(testData.suggestions); i++ {
		testza.AssertContains(t, testData.initialLines[i], testData.suggestions[i-1].Name)
		testza.AssertContains(t, testData.initialLines[i], testData.suggestions[i-1].Description)
	}

	teardown(t, testData.tester)
}

func TestFilter(t *testing.T) {
	testData := setup(t)

	testData.tester.SendString("fi")
	// Check that typed input shows up
	_, lines, err := testData.tester.WaitFor(func(out string, outputLines []string) bool {
		return strings.Contains(outputLines[0], "fi")
	})
	testza.AssertNoError(t, err)

	// Check that suggestions filter properly
	testza.AssertEqual(t, 3, len(lines))
	testza.AssertContains(t, lines[1], testData.suggestions[0].Name)
	testza.AssertContains(t, lines[1], testData.suggestions[0].Description)
	testza.AssertContains(t, lines[2], testData.suggestions[4].Name)
	testza.AssertContains(t, lines[2], testData.suggestions[4].Description)

	teardown(t, testData.tester)
}

func testExecutor(t *testing.T, in *string, expectedOut string) {
	testData := setup(t)

	if in != nil {
		testData.tester.SendString(*in)
		// Wait for typed input to render
		_, _, err := testData.tester.WaitFor(func(out string, outputLines []string) bool {
			return strings.Contains(outputLines[0], *in)
		})
		testza.AssertNoError(t, err)
	}

	testData.tester.SendByte(tuitest.KeyEnter)

	// Check that executor output displays
	_, _, err := testData.tester.WaitFor(func(out string, outputLines []string) bool {
		return len(outputLines) > 1 && strings.Contains(outputLines[1], expectedOut)
	})
	testza.AssertNoError(t, err)

	teardown(t, testData.tester)
}

func TestExecutorNoInput(t *testing.T) {
	testExecutor(t, nil, "result is")
}

func TestExecutorWithInput(t *testing.T) {
	in := "fi"
	testExecutor(t, &in, "result is fi")
}

func TestChoosePrompt(t *testing.T) {
	testData := setup(t)
	testData.tester.RemoveAnsi = false

	testData.tester.SendString(tuitest.KeyDown)
	// Wait for first prompt to be selected
	_, lines, err := testData.tester.WaitFor(func(out string, outputLines []string) bool {
		return strings.Contains(outputLines[0], testData.suggestions[0].Name)
	})
	testza.AssertNoError(t, err)
	// Check that proper styles are applied
	testza.AssertContains(t, lines[0], testData.model.prompt.Formatters.SelectedSuggestion.Render(testData.suggestions[0].Name))
	testza.AssertContains(t, lines[0], testData.model.prompt.Formatters.Placeholder.format(testData.suggestions[0].Placeholder))
	maxNameLen := len("second option")
	testza.AssertContains(t, lines[1], testData.model.prompt.Formatters.Name.format(testData.suggestions[0].Name, maxNameLen, true))
	maxDescLen := len("test desc1")
	testza.AssertContains(t, lines[1], testData.model.prompt.Formatters.Description.format(testData.suggestions[0].Description, maxDescLen, true))

	// Check that the selected text gets sent to the executor without the placeholder
	testData.tester.SendByte(tuitest.KeyEnter)
	_, _, err = testData.tester.WaitFor(func(out string, outputLines []string) bool {
		return len(outputLines) > 1 &&
			strings.Contains(outputLines[1], "result is "+testData.suggestions[0].Name) &&
			!strings.Contains(outputLines[1], testData.suggestions[0].Placeholder)
	})
	testza.AssertNoError(t, err)

	teardown(t, testData.tester)
}

func TestTypeAfterCompleting(t *testing.T) {
	testData := setup(t)

	testData.tester.SendString(tuitest.KeyDown)
	// Wait for first prompt to be selected
	_, _, err := testData.tester.WaitFor(func(out string, outputLines []string) bool {
		return strings.Contains(outputLines[0], testData.suggestions[0].Name)
	})
	testza.AssertNoError(t, err)

	testData.tester.SendString("a")
	// Check that text updates
	_, lines, err := testData.tester.WaitFor(func(out string, outputLines []string) bool {
		return strings.Contains(outputLines[0], testData.suggestions[0].Name+"a")
	})
	testza.AssertNoError(t, err)
	// Check that prompts were filtered
	testza.AssertEqual(t, 1, len(lines))
	// Check that selected text formatting was removed
	testza.AssertNotContains(t, lines[0], testData.model.prompt.Formatters.SelectedSuggestion.Render(testData.suggestions[0].Name+"a"))

	teardown(t, testData.tester)
}
