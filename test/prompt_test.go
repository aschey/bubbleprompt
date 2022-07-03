package test

import (
	"fmt"
	"strings"
	"time"

	prompt "github.com/aschey/bubbleprompt"
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/input/commandinput"
	"github.com/aschey/bubbleprompt/test/testapp"
	tuitest "github.com/aschey/tui-tester"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func testExecutor(console *tuitest.Console, in *string, backspace bool, doubleEnter bool, outStr string) {
	if in != nil {
		console.SendString(*in)
		// Wait for typed input to render
		_, _ = console.WaitFor(func(state tuitest.TermState) bool {
			return strings.Contains(state.NthOutputLine(0), *in)
		})
	}

	if in != nil && backspace {
		// Send input twice quickly to test completer ignoring first input
		console.SendString(tuitest.KeyBackspace)
		in := *in
		console.SendString(string(in[len(in)-1]))
	}
	console.SendString(tuitest.KeyEnter)
	if doubleEnter {
		// Hit enter twice to re-trigger completer
		console.SendString(tuitest.KeyEnter)
	}

	_, _ = console.WaitFor(func(state tuitest.TermState) bool {
		return strings.Contains(state.NthOutputLine(1), outStr)
	})
}

var tester tuitest.Tester = tuitest.Tester{}
var _ = BeforeSuite(func() {
	var err error
	tester, err = tuitest.NewTester("./testapp")
	Expect(err).ShouldNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	err := tester.TearDown()
	Expect(err).ShouldNot(HaveOccurred())
})

var suggestions []input.Suggestion[commandinput.CmdMetadata] = testapp.Suggestions

var _ = Describe("Prompt", FlakeAttempts(2), func() {
	leftPadding := 2
	margin := 1
	longestNameLength := len("seventh-option")
	longestDescLength := len("test desc2")
	promptWidth := leftPadding + margin + longestNameLength + 2*margin + longestDescLength + margin

	var console *tuitest.Console
	var initialLines []string
	_ = initialLines

	BeforeEach(OncePerOrdered, func() {
		var err error
		console, err = tester.NewConsole([]string{})
		Expect(err).ShouldNot(HaveOccurred())

		console.OnError = func(err error) error {
			defer GinkgoRecover()
			Expect(err).Error().ShouldNot(HaveOccurred())
			return err
		}

		Expect(err).Error().ShouldNot(HaveOccurred())

		// Wait for prompt to initialize
		console.TrimOutput = true
		state, _ := console.WaitFor(func(state tuitest.TermState) bool {
			return strings.Contains(state.NthOutputLine(6), suggestions[5].Description)
		})
		initialLines = state.OutputLines()
	})

	AfterEach(OncePerOrdered, func() {
		console.SendString(tuitest.KeyCtrlC)
		err := console.WaitForTermination()
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
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				for i := 1; i < 6; i++ {
					if fmt.Sprint(state.BgColor(1, promptWidth)) != prompt.DefaultScrollbarThumbColor {
						return false
					}
				}
				return true
			})

			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return fmt.Sprint(state.BgColor(6, promptWidth)) == prompt.DefaultScrollbarColor
			})

		})
	})

	When("the user types to filter the prompt", Ordered, func() {
		var lines []string
		_ = lines
		BeforeAll(func() {
			console.SendString("fi")
		})

		It("shows the typed filter", func() {
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return strings.Contains(state.NthOutputLine(0), "> fi")
			})
		})

		It("filters the suggestions", func() {
			state, _ := console.WaitFor(func(state tuitest.TermState) bool {
				return state.NumLines() == 3 && strings.Contains(state.NthOutputLine(2), suggestions[4].Description)
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
			testExecutor(console, nil, false, false, "result is")
		})
	})

	When("the user presses enter twice without filtering", func() {
		It("shows the output", func() {
			testExecutor(console, nil, false, true, "result is")
		})
	})

	When("the user presses enter after typing", func() {
		It("shows the output", func() {
			in := "fi"
			testExecutor(console, &in, false, false, "result is fi")
		})
	})

	When("the user presses enter twice after typing", func() {
		It("shows the output", func() {
			in := "fi"
			testExecutor(console, &in, false, true, "result is fi")
		})
	})

	When("the user presses enter after typing and pressing backspace", func() {
		It("shows the output", func() {
			in := "fi"
			testExecutor(console, &in, true, false, "result is fi")
		})
	})

	When("the user presses enter twice after typing and pressing backspace", func() {
		It("shows the output", func() {
			in := "fi"
			testExecutor(console, &in, true, true, "result is fi")
		})
	})

	When("the user presses the down arrow", Ordered, func() {
		BeforeAll(func() {
			console.SendString(tuitest.KeyDown)
		})

		It("selects the first prompt", func() {
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return strings.Contains(state.NthOutputLine(0), suggestions[0].Text)
			})
		})

		It("applies the selected text styling", func() {
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return fmt.Sprint(state.FgColor(0, leftPadding)) == commandinput.DefaultSelectedTextColor
			})
		})

		It("applies the selected placeholder styling", func() {
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return fmt.Sprint(state.FgColor(0, leftPadding+margin+len(suggestions[0].Text))) == commandinput.DefaultPlaceholderForeground
			})
		})

		It("applies the correct background for the suggestion name so it covers the longest name", func() {
			maxNameLen := len("seventh-option")
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return fmt.Sprint(state.BgColor(1, leftPadding+maxNameLen+margin)) == input.DefaultNameBackground
			})
		})

		It("applies the correct background for the suggestion description so it covers the longest description", func() {
			maxNameLen := len("seventh-option")
			maxDescLen := len("test desc2")
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return fmt.Sprint(state.BgColor(1, leftPadding+maxNameLen+2*margin+maxDescLen+margin)) == input.DefaultDescriptionBackground
			})
		})
	})

	When("the user chooses the first prompt", Ordered, func() {
		BeforeAll(func() {
			console.SendString(tuitest.KeyDown)
			console.SendString(tuitest.KeyEnter)
		})

		It("renders the executor result", func() {
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				secondLine := state.NthOutputLine(1)
				return strings.Contains(secondLine, "result is "+suggestions[0].Text) &&
					!strings.Contains(secondLine, suggestions[0].Metadata.PositionalArgs()[0].Placeholder)
			})
		})
	})

	When("the user filters all the prompts", Ordered, func() {
		BeforeAll(func() {
			console.SendString(tuitest.KeyDown)
			console.SendString("a")
		})

		It("displays the input", func() {
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return strings.Contains(state.NthOutputLine(0), suggestions[0].Text+"a")
			})
		})

		It("does not display any prompts", func() {
			_, _ = console.WaitForDuration(func(state tuitest.TermState) bool {
				return state.NumLines() == 1
			}, 100*time.Millisecond)
		})

		It("removes the selected text styling", func() {
			_, _ = console.WaitForDuration(func(state tuitest.TermState) bool {
				return state.FgColor(0, 2) == tuitest.DefaultFG
			}, 100*time.Millisecond)
		})
	})

	When("the user scrolls down", Ordered, func() {
		BeforeAll(func() {
			for i := 0; i < 7; i++ {
				console.SendString(tuitest.KeyDown)
			}
		})

		It("updates the scrollbar", func() {
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				for i := 2; i < 7; i++ {
					if fmt.Sprint(state.BgColor(i, promptWidth)) != prompt.DefaultScrollbarThumbColor {
						return false
					}
				}
				return true
			})

			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return fmt.Sprint(state.BgColor(1, promptWidth)) == prompt.DefaultScrollbarColor
			})
		})
	})

	When("then user scrolls back up", Ordered, func() {
		BeforeAll(func() {
			for i := 0; i < 7; i++ {
				console.SendString(tuitest.KeyDown)
			}
			for i := 0; i < 7; i++ {
				console.SendString(tuitest.KeyUp)
			}
		})

		It("updates the scrollbar", func() {
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				for i := 1; i < 6; i++ {
					if fmt.Sprint(state.BgColor(1, promptWidth)) != prompt.DefaultScrollbarThumbColor {
						return false
					}
				}
				return true
			})

			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return fmt.Sprint(state.BgColor(6, promptWidth)) == prompt.DefaultScrollbarColor
			})
		})
	})

	When("the user moves the cursor left", Ordered, func() {
		BeforeAll(func() {
			console.SendString(tuitest.KeyDown)
			for i := 0; i < 10; i++ {
				console.SendString(tuitest.KeyLeft)
			}
		})

		It("shows the completions matching the prefix", func() {
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {

				return state.NumLines() == 3 &&
					strings.Contains(state.NthOutputLine(1), suggestions[0].Text) &&
					strings.Contains(state.NthOutputLine(2), suggestions[4].Text)
			})
		})
	})

	When("the user chooses an item and presses the spacebar", Ordered, func() {
		BeforeAll(func() {
			console.SendString(tuitest.KeyDown)
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return fmt.Sprint(state.FgColor(1, leftPadding+margin)) == input.DefaultSelectedForeground
			})
			console.SendString(" ")
		})

		It("unselects the suggestion", func() {
			_, _ = console.WaitForDuration(func(state tuitest.TermState) bool {
				return state.FgColor(1, leftPadding+margin+len(suggestions[0].Text)+margin) == tuitest.DefaultFG
			}, 100*time.Millisecond)
		})
	})

	When("the user cycles through all suggestions", Ordered, func() {
		BeforeAll(func() {
			for i := 0; i < 7; i++ {
				console.SendString(tuitest.KeyDown)
			}
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return fmt.Sprint(state.FgColor(6, leftPadding+margin)) == input.DefaultSelectedForeground
			})
			console.SendString(tuitest.KeyDown)
		})

		It("unselects the suggestion", func() {
			_, _ = console.WaitForDuration(func(state tuitest.TermState) bool {
				return state.FgColor(6, leftPadding+margin) == tuitest.DefaultFG
			}, 100*time.Millisecond)
		})
	})

	When("the user types the full suggestion", Ordered, func() {
		BeforeAll(func() {
			console.SendString(suggestions[0].Text)
		})

		It("selects the suggestion", func() {
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return fmt.Sprint(state.FgColor(0, leftPadding)) == commandinput.DefaultSelectedTextColor
			})
		})
	})
})
