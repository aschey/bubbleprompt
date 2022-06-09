package prompt

import (
	"fmt"
	"os"
	"strings"
	"time"

	//"github.com/Netflix/go-expect"

	completers "github.com/aschey/bubbleprompt/completer"
	executors "github.com/aschey/bubbleprompt/executor"
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/input/commandinput"
	tuitest "github.com/aschey/tui-tester"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type cmdMetadata = commandinput.CmdMetadata

type model struct {
	prompt Model[cmdMetadata]
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

type testCompleterModel struct {
	suggestions []input.Suggestion[cmdMetadata]
}

func (m testCompleterModel) completer(document Document, promptModel Model[cmdMetadata]) []input.Suggestion[cmdMetadata] {
	time.Sleep(2 * time.Millisecond)
	return completers.FilterHasPrefix(document.TextBeforeCursor(), m.suggestions)
}

type testExecutorModel struct{}

func (m testExecutorModel) executor(input string) (tea.Model, error) {
	time.Sleep(2 * time.Millisecond)
	return executors.NewStringModel("result is " + input), nil
}

func testExecutor(tester *tuitest.Tester, in *string, backspace bool, doubleEnter bool, outStr string) {
	if in != nil {
		tester.SendString(*in)
		// Wait for typed input to render
		_, _ = tester.WaitFor(func(state tuitest.TermState) bool {
			return strings.Contains(state.OutputLines()[0], *in)
		})
	}

	if in != nil && backspace {
		// Send input twice quickly to test completer ignoring first input
		tester.SendByte(tuitest.KeyBackspace)
		in := *in
		tester.SendString(string(in[len(in)-1]))
	}
	tester.SendByte(tuitest.KeyEnter)
	if doubleEnter {
		// Hit enter twice to re-trigger completer
		tester.SendByte(tuitest.KeyEnter)
	}

	_, _ = tester.WaitFor(func(state tuitest.TermState) bool {
		outputLines := state.OutputLines()
		return len(outputLines) > 1 && strings.Contains(outputLines[1], outStr)
	})
}

var _ = Describe("Prompt", func() {
	suggestions := []input.Suggestion[cmdMetadata]{
		{Text: "first-option", Description: "test desc", Metadata: commandinput.NewCmdMetadata([]commandinput.PositionalArg{commandinput.NewPositionalArg("[test placeholder]")}, commandinput.Placeholder{})},
		{Text: "second-option", Description: "test desc2"},
		{Text: "third-option", Description: "test desc3"},
		{Text: "fourth-option", Description: "test desc4"},
		{Text: "fifth-option", Description: "test desc5"},
		{Text: "sixth-option", Description: "test desc6"},
		{Text: "seventh-option", Description: "test desc7"},
	}
	leftPadding := 2
	margin := 1
	longestNameLength := len("seventh-option")
	longestDescLength := len("test desc2")
	promptWidth := leftPadding + margin + longestNameLength + 2*margin + longestDescLength + margin

	var promptModel *model
	var tester *tuitest.Tester
	var initialLines []string

	BeforeEach(OncePerOrdered, func() {
		completerModel := testCompleterModel{suggestions: suggestions}
		executorModel := testExecutorModel{}

		var input input.Input[cmdMetadata] = commandinput.New[cmdMetadata]()
		promptModel = &model{
			prompt: New(
				completerModel.completer,
				executorModel.executor,
				input,
			),
		}

		var err error
		tester, err = tuitest.New(func(tty *os.File) error {
			return tea.NewProgram(promptModel, tea.WithInput(tty), tea.WithOutput(tty)).Start()
		})
		tester.OnError = func(err error) error {
			defer GinkgoRecover()
			Expect(err).Error().ShouldNot(HaveOccurred())
			return err
		}

		Expect(err).Error().ShouldNot(HaveOccurred())

		// Wait for prompt to initialize
		tester.TrimOutput = true
		state, _ := tester.WaitFor(func(state tuitest.TermState) bool {
			outputLines := state.OutputLines()
			return len(outputLines) > 6 && strings.Contains(outputLines[6], suggestions[5].Description)
		})
		initialLines = state.OutputLines()
	})

	AfterEach(OncePerOrdered, func() {
		tester.SendByte(byte(tea.KeyCtrlC))
		err := tester.WaitForTermination()
		Expect(err).Error().ShouldNot(HaveOccurred())
	})

	When("the prompt is loaded", Ordered, func() {
		It("shows the completion text", func() {
			for i := 0; i < 6; i++ {
				Expect(initialLines[i+1]).To(ContainSubstring(suggestions[i].Text))
			}
		})

		It("shows the completion description", func() {
			for i := 0; i < 6; i++ {
				Expect(initialLines[i+1]).To(ContainSubstring(suggestions[i].Text))
			}
		})

		It("shows the scrollbar", func() {
			_, _ = tester.WaitFor(func(state tuitest.TermState) bool {
				for i := 1; i < 6; i++ {
					if fmt.Sprint(state.BgColor(1, promptWidth)) != DefaultScrollbarThumbColor {
						return false
					}
				}
				return true
			})

			_, _ = tester.WaitFor(func(state tuitest.TermState) bool {
				return fmt.Sprint(state.BgColor(6, promptWidth)) == DefaultScrollbarColor
			})

		})
	})

	When("the user types to filter the prompt", Ordered, func() {
		var lines []string
		BeforeAll(func() {
			tester.SendString("fi")
		})

		It("shows the typed filter", func() {
			_, _ = tester.WaitFor(func(state tuitest.TermState) bool {
				return strings.Contains(state.OutputLines()[0], "> fi")
			})
		})

		It("filters the suggestions", func() {
			state, _ := tester.WaitFor(func(state tuitest.TermState) bool {
				outputLines := state.OutputLines()
				return len(outputLines) == 3 && strings.Contains(outputLines[2], suggestions[4].Description)
			})
			lines = state.OutputLines()
		})

		It("shows the correct suggestions", func() {
			Expect(lines[1]).To(ContainSubstring(suggestions[0].Text))
			Expect(lines[1]).To(ContainSubstring(suggestions[0].Description))
			Expect(lines[2]).To(ContainSubstring(suggestions[4].Text))
			Expect(lines[2]).To(ContainSubstring(suggestions[4].Description))
		})
	})

	When("the user presses enter without filtering", func() {
		It("shows the output", func() {
			testExecutor(tester, nil, false, false, "result is")
		})
	})

	When("the user presses enter twice without filtering", func() {
		It("shows the output", func() {
			testExecutor(tester, nil, false, true, "result is")
		})
	})

	When("the user presses enter after typing", func() {
		It("shows the output", func() {
			in := "fi"
			testExecutor(tester, &in, false, false, "result is fi")
		})
	})

	When("the user presses enter twice after typing", func() {
		It("shows the output", func() {
			in := "fi"
			testExecutor(tester, &in, false, true, "result is fi")
		})
	})

	When("the user presses enter after typing and pressing backspace", func() {
		It("shows the output", func() {
			in := "fi"
			testExecutor(tester, &in, true, false, "result is fi")
		})
	})

	When("the user presses enter twice after typing and pressing backspace", func() {
		It("shows the output", func() {
			in := "fi"
			testExecutor(tester, &in, true, true, "result is fi")
		})
	})

	When("the user presses the down arrow", Ordered, func() {
		BeforeAll(func() {
			tester.SendString(tuitest.KeyDown)
		})

		It("selects the first prompt", func() {
			_, _ = tester.WaitFor(func(state tuitest.TermState) bool {
				return strings.Contains(state.OutputLines()[0], suggestions[0].Text)
			})
		})

		It("applies the selected text styling", func() {
			_, _ = tester.WaitFor(func(state tuitest.TermState) bool {
				return fmt.Sprint(state.FgColor(0, leftPadding)) == commandinput.DefaultSelectedTextColor
			})
		})

		It("applies the selected placeholder styling", func() {
			_, _ = tester.WaitFor(func(state tuitest.TermState) bool {
				return fmt.Sprint(state.FgColor(0, leftPadding+margin+len(suggestions[0].Text))) == commandinput.DefaultPlaceholderForeground
			})
		})

		It("applies the correct background for the suggestion name so it covers the longest name", func() {
			maxNameLen := len("seventh-option")
			_, _ = tester.WaitFor(func(state tuitest.TermState) bool {
				return fmt.Sprint(state.BgColor(1, leftPadding+maxNameLen+margin)) == input.DefaultNameBackground
			})
		})

		It("applies the correct background for the suggestion description so it covers the longest description", func() {
			maxNameLen := len("seventh-option")
			maxDescLen := len("test desc2")
			_, _ = tester.WaitFor(func(state tuitest.TermState) bool {
				return fmt.Sprint(state.BgColor(1, leftPadding+maxNameLen+2*margin+maxDescLen+margin)) == input.DefaultDescriptionBackground
			})
		})
	})

	When("the user chooses the first prompt", Ordered, func() {
		BeforeAll(func() {
			tester.SendString(tuitest.KeyDown)
			tester.SendByte(tuitest.KeyEnter)
		})

		It("renders the executor result", func() {
			_, _ = tester.WaitFor(func(state tuitest.TermState) bool {
				outputLines := state.OutputLines()
				return len(outputLines) > 1 &&
					strings.Contains(outputLines[1], "result is "+suggestions[0].Text) &&
					!strings.Contains(outputLines[1], suggestions[0].Metadata.PositionalArgs()[0].Placeholder)
			})
		})
	})

	When("the user filters all the prompts", Ordered, func() {
		BeforeAll(func() {
			tester.SendString(tuitest.KeyDown)
			tester.SendString("a")
		})

		It("displays the input", func() {
			_, _ = tester.WaitFor(func(state tuitest.TermState) bool {
				return strings.Contains(state.OutputLines()[0], suggestions[0].Text+"a")
			})
		})

		It("does not display any prompts", func() {
			_, _ = tester.WaitForDuration(func(state tuitest.TermState) bool {
				return len(state.OutputLines()) == 1
			}, 100*time.Millisecond)
		})

		It("removes the selected text styling", func() {
			_, _ = tester.WaitForDuration(func(state tuitest.TermState) bool {
				return state.FgColor(0, 2) == tuitest.DefaultFG
			}, 100*time.Millisecond)
		})
	})

	When("the user scrolls down", Ordered, func() {
		BeforeAll(func() {
			for i := 0; i < 8; i++ {
				tester.SendString(tuitest.KeyDown)
				time.Sleep(10 * time.Millisecond)
			}
		})

		It("updates the scrollbar", func() {

			_, _ = tester.WaitFor(func(state tuitest.TermState) bool {
				for i := 2; i < 7; i++ {
					if fmt.Sprint(state.BgColor(1, promptWidth)) != DefaultScrollbarThumbColor {
						return false
					}
				}
				return true
			})

			// _, _ = tester.WaitFor(func(state tuitest.TermState) bool {
			// 	return fmt.Sprint(state.BgColor(1, promptWidth)) == DefaultScrollbarColor
			// })

		})
	})
})
