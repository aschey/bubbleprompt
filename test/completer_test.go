package test

import (
	"fmt"
	"strings"
	"time"

	prompt "github.com/aschey/bubbleprompt"
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/input/commandinput"
	tuitest "github.com/aschey/tui-tester"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Completer", func() {
	longestNameLength := len("seventh-option")
	longestDescLength := len("test desc2")
	promptWidth := leftPadding + margin + longestNameLength + 2*margin + longestDescLength + margin

	var console *tuitest.Console
	var initialLines []string
	_ = initialLines

	BeforeEach(OncePerOrdered, func() {
		console, _ = cmdTester.CreateConsole()

		// Wait for prompt to initialize
		console.TrimOutput = true
		state, _ := console.WaitFor(func(state tuitest.TermState) bool {
			return strings.Contains(state.NthOutputLine(6), suggestions[5].Description)
		})
		initialLines = state.OutputLines()
	})

	AfterEach(OncePerOrdered, func() {
		console.SendString(tuitest.KeyCtrlC)
		_ = console.WaitForTermination()
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
				return fmt.Sprint(state.BgColor(1, leftPadding+maxNameLen+margin)) == input.DefaultSelectedNameBackground
			})
		})

		It("applies the correct background for the suggestion description so it covers the longest description", func() {
			maxNameLen := len("seventh-option")
			maxDescLen := len("test desc2")
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return fmt.Sprint(state.BgColor(1, leftPadding+maxNameLen+2*margin+maxDescLen+margin)) == input.DefaultSelectedDescriptionBackground
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
				return fmt.Sprint(state.FgColor(1, leftPadding+margin)) == input.DefaultSelectedNameForeground
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
				return fmt.Sprint(state.FgColor(6, leftPadding+margin)) == input.DefaultSelectedNameForeground
			})
			console.SendString(tuitest.KeyDown)
		})

		It("unselects the suggestion", func() {
			_, _ = console.WaitForDuration(func(state tuitest.TermState) bool {
				return fmt.Sprint(state.FgColor(6, leftPadding+margin)) == input.DefaultNameForeground
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

	When("the update function sends a oneshot completer message", Ordered, func() {
		BeforeAll(func() {
			console.SendString("oneshot")
			console.SendString(tuitest.KeyEnter)
		})

		It("Updates the completions", func() {
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return strings.Contains(state.NthOutputLine(3), "changed text")
			})
		})
	})
})
