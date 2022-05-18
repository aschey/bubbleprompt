package prompt

import (
	"strings"
	"testing"

	completers "github.com/aschey/bubbleprompt/completer"
	executors "github.com/aschey/bubbleprompt/executor"
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/input/commandinput"
	tuitest "github.com/aschey/tui-tester"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sclevine/spec"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type cmdMetadata = commandinput.CmdMetadata

type model struct {
	prompt Model[cmdMetadata]
}

type testCompleterModel struct {
	suggestions []input.Suggestion[cmdMetadata]
}

type testExecutorModel struct{}

type TestSuite struct {
	suite.Suite
	suggestions  []input.Suggestion[cmdMetadata]
	tester       tuitest.Tester
	initialLines []string
	model        model
	textInput    *commandinput.Model[cmdMetadata]
}

func (suite *TestSuite) SetupTest() {
	suggestions := []input.Suggestion[cmdMetadata]{
		{Text: "first-option", Description: "test desc", Metadata: commandinput.NewCmdMetadata([]commandinput.PositionalArg{{Placeholder: "[test placeholder]"}}, commandinput.Placeholder{})},
		{Text: "second-option", Description: "test desc2"},
		{Text: "third-option", Description: "test desc3"},
		{Text: "fourth-option", Description: "test desc4"},
		{Text: "fifth-option", Description: "test desc5"},
	}

	completerModel := testCompleterModel{suggestions: suggestions}
	executorModel := testExecutorModel{}

	var textInput input.Input[cmdMetadata] = commandinput.New[cmdMetadata]()
	model := model{
		prompt: New(
			completerModel.completer,
			executorModel.executor,
			textInput,
		),
	}
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

	require.NoError(suite.T(), err)
	suite.suggestions = suggestions
	suite.tester = tester
	suite.initialLines = initialLines
	suite.model = model
	suite.textInput = textInput.(*commandinput.Model[cmdMetadata])
}

func (suite *TestSuite) TearDownTest() {
	suite.tester.SendByte(tuitest.KeyCtrlC)
	err := suite.tester.WaitForTermination()
	require.NoError(suite.T(), err)
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

func (m testCompleterModel) completer(document Document, promptModel Model[cmdMetadata]) []input.Suggestion[cmdMetadata] {
	return completers.FilterHasPrefix(document.TextBeforeCursor(), m.suggestions)
}

func (m testExecutorModel) executor(input string) (tea.Model, error) {
	return executors.NewStringModel("result is " + input), nil
}

func TestAll(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) TestBasicCompleter() {
	spec.Run(suite.T(), "basicCompleter", func(t *testing.T, when spec.G, it spec.S) {
		// Check that all prompts show up
		for i := 1; i < len(suite.suggestions); i++ {
			it("should have the initial text for "+suite.initialLines[i], func() {
				require.Contains(t, suite.initialLines[i], suite.suggestions[i-1].Text)
				require.Contains(t, suite.initialLines[i], suite.suggestions[i-1].Description)
			})
		}
	})
}

func (suite *TestSuite) TestFilter() {
	spec.Run(suite.T(), "filter", func(t *testing.T, when spec.G, it spec.S) {
		when("the user types to filter the prompt", func() {
			it.Before(func() {
				suite.tester.SendString("fi")
			})

			it("should filter the suggestions", func() {
				_, lines, err := suite.tester.WaitFor(func(out string, outputLines []string) bool {
					return strings.Contains(outputLines[0], "fi")
				})
				require.NoError(t, err)

				require.Equal(t, 3, len(lines))
				require.Contains(t, lines[1], suite.suggestions[0].Text)
				require.Contains(t, lines[1], suite.suggestions[0].Description)
				require.Contains(t, lines[2], suite.suggestions[4].Text)
				require.Contains(t, lines[2], suite.suggestions[4].Description)
			})
		})
	})
}

func (suite *TestSuite) testExecutor(text string, in *string, backspace bool, doubleEnter bool, expectedOut string) {
	spec.Run(suite.T(), "executor"+text, func(t *testing.T, when spec.G, it spec.S) {
		when("the user presses enter", func() {
			it.Before(func() {
				if in != nil {
					suite.tester.SendString(*in)
					// Wait for typed input to render
					_, _, err := suite.tester.WaitFor(func(out string, outputLines []string) bool {
						return strings.Contains(outputLines[0], *in)
					})
					require.NoError(t, err)

				}
				if in != nil && backspace {
					// Send input twice quickly to test completer ignoring first input
					suite.tester.SendByte(tuitest.KeyBackspace)
					i := *in
					suite.tester.SendString(string(i[len(i)-1]))
				}
				suite.tester.SendByte(tuitest.KeyEnter)
				if doubleEnter {
					// Hit enter twice to re-trigger completer
					suite.tester.SendByte(tuitest.KeyEnter)
				}

			})
			it("should display the executor output", func() {
				_, _, err := suite.tester.WaitFor(func(out string, outputLines []string) bool {
					return len(outputLines) > 1 && strings.Contains(outputLines[1], expectedOut)
				})
				require.NoError(t, err)
			})
		})
	})
}

func (suite *TestSuite) TestExecutor() {
	in := "fi"
	suite.testExecutor("NoInput", nil, false, false, "result is")
	suite.SetupTest()
	suite.testExecutor("NoInputDoubleEnter", nil, false, true, "result is")
	suite.SetupTest()
	suite.testExecutor("WithInput", &in, false, false, "result is fi")
	suite.SetupTest()
	suite.testExecutor("WithInputDoubleEnter", &in, false, true, "result is fi")
	suite.SetupTest()
	suite.testExecutor("WithInputBackspace", &in, true, false, "result is fi")
	suite.SetupTest()
	suite.testExecutor("WithInputDoubleEnterAndBackspace", &in, true, true, "result is fi")
}

func (suite *TestSuite) TestChoosePrompt() {
	suite.tester.RemoveAnsi = false
	spec.Run(suite.T(), "choosePrompt", func(t *testing.T, when spec.G, it spec.S) {
		when("the user presses the down arrow", func() {
			it.Before(func() {
				suite.tester.SendString(tuitest.KeyDown)
			})

			it("selects the first prompt", func() {
				// Wait for first prompt to be selected
				_, lines, err := suite.tester.WaitFor(func(out string, outputLines []string) bool {
					return strings.Contains(outputLines[0], suite.suggestions[0].Text)
				})

				require.NoError(t, err)
				// Check that proper styles are applied
				require.Contains(t, lines[0], suite.textInput.SelectedTextStyle.Render(suite.suggestions[0].Text))
				require.Contains(t, lines[0], suite.suggestions[0].Metadata.PositionalArgs()[0].PlaceholderStyle.Format(suite.suggestions[0].Metadata.PositionalArgs()[0].Placeholder))
				maxNameLen := len("second-option")
				require.Contains(t, lines[1], suite.model.prompt.Formatters.Name.Format(suite.suggestions[0].Text, maxNameLen, true))
				maxDescLen := len("test desc1")
				require.Contains(t, lines[1], suite.model.prompt.Formatters.Description.Format(suite.suggestions[0].Description, maxDescLen, true))
			})
		})

		when("the user chooses the prompt", func() {
			it.Before(func() {
				suite.tester.SendByte(tuitest.KeyEnter)
			})

			it("renders the executor result", func() {
				_, _, err := suite.tester.WaitFor(func(out string, outputLines []string) bool {
					// Check that the selected text gets sent to the executor without the placeholder
					return len(outputLines) > 1 &&
						strings.Contains(outputLines[1], "result is "+suite.suggestions[0].Text) &&
						!strings.Contains(outputLines[1], suite.suggestions[0].Metadata.PositionalArgs()[0].Placeholder)
				})
				require.NoError(t, err)
			})
		})
	})

}

func (suite *TestSuite) TestTypeAfterCompleting() {
	spec.Run(suite.T(), "typeAfterCompleting", func(t *testing.T, when spec.G, it spec.S) {
		when("the user presses the down arrow", func() {
			it.Before(func() {
				suite.tester.SendString(tuitest.KeyDown)
			})

			it("selects the first prompt", func() {
				// Wait for first prompt to be selected
				_, _, err := suite.tester.WaitFor(func(out string, outputLines []string) bool {
					return strings.Contains(outputLines[0], suite.suggestions[0].Text)
				})
				require.NoError(t, err)
			})
		})

		when("the user types input", func() {
			it.Before(func() {
				suite.tester.SendString("a")
			})

			it("should filter the prompts", func() {
				_, lines, err := suite.tester.WaitFor(func(out string, outputLines []string) bool {
					return strings.Contains(outputLines[0], suite.suggestions[0].Text+"a")
				})
				require.NoError(t, err)
				// Check that prompts were filtered
				require.Equal(t, 1, len(lines))
				// Check that selected text formatting was removed
				require.NotContains(t, lines[0], suite.textInput.SelectedTextStyle.Render(suite.suggestions[0].Text+"a"))
			})
		})
	})
}
