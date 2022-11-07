package test

import (
	"fmt"
	"strings"
	"time"

	"github.com/aschey/bubbleprompt/input"
	tuitest "github.com/aschey/tui-tester"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Completer", func() {
	longestNameLength := len("completion text")
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
					if state.BackgroundColor(1, promptWidth).String() != input.DefaultScrollbarThumbColor {
						return false
					}
				}
				return true
			})

			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return state.BackgroundColor(6, promptWidth).String() == input.DefaultScrollbarColor
			})

		})
	})

	When("the user presses the down arrow", Ordered, func() {
		BeforeAll(func() {
			console.SendString(tuitest.KeyDown)
		})

		It("applies the correct background for the suggestion name so it covers the longest name", func() {
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return state.BackgroundColor(1, leftPadding+longestNameLength+margin).String() == input.DefaultSelectedNameBackground
			})
		})

		It("applies the correct background for the suggestion description so it covers the longest description", func() {
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return state.BackgroundColor(1, leftPadding+longestNameLength+2*margin+longestDescLength+margin).String() == input.DefaultSelectedDescriptionBackground
			})
		})
	})

	When("the user presses the tab key", Ordered, func() {
		BeforeAll(func() {
			console.SendString(tuitest.KeyTab)
		})

		It("selects the suggestion", func() {
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return state.BackgroundColor(1, leftPadding).String() == input.DefaultSelectedNameBackground
			})
		})

		When("the user presses the space key", Ordered, func() {
			BeforeAll(func() {
				console.SendString(" ")
			})

			It("moves the suggestions over", func() {
				_, _ = console.WaitFor(func(state tuitest.TermState) bool {
					return state.BackgroundColor(1, leftPadding+len(suggestions[0].Text)+1).String() == input.DefaultNameBackground
				})
			})
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
					if state.BackgroundColor(i, promptWidth).String() != input.DefaultScrollbarThumbColor {
						return false
					}
				}
				return true
			})

			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return state.BackgroundColor(1, promptWidth).String() == input.DefaultScrollbarColor
			})
		})

		When("then user scrolls back up", Ordered, func() {
			BeforeAll(func() {
				for i := 0; i < 7; i++ {
					console.SendString(tuitest.KeyUp)
				}
			})

			It("updates the scrollbar", func() {
				_, _ = console.WaitFor(func(state tuitest.TermState) bool {
					for i := 1; i < 6; i++ {
						if state.BackgroundColor(1, promptWidth).String() != input.DefaultScrollbarThumbColor {
							return false
						}
					}
					return true
				})

				_, _ = console.WaitFor(func(state tuitest.TermState) bool {
					return state.BackgroundColor(6, promptWidth).String() == input.DefaultScrollbarColor
				})
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

		It("moves the suggestions to match the cursor", func() {
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return state.BackgroundColor(1, 3).Int() == tuitest.DefaultBG && fmt.Sprint(state.BackgroundColor(1, 4)) == input.DefaultNameBackground
			})
		})
	})

	When("the user chooses an item and presses the spacebar", Ordered, func() {
		BeforeAll(func() {
			console.SendString(tuitest.KeyDown)
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return state.ForegroundColor(1, leftPadding+margin).String() == input.DefaultSelectedNameForeground
			})
			console.SendString(" ")
		})

		It("unselects the suggestion", func() {
			_, _ = console.WaitForDuration(func(state tuitest.TermState) bool {
				return state.ForegroundColor(1, leftPadding+margin+len(suggestions[0].Text)+margin).Int() == tuitest.DefaultFG
			}, 100*time.Millisecond)
		})
	})

	When("the user cycles through all suggestions", Ordered, func() {
		BeforeAll(func() {
			for i := 0; i < 7; i++ {
				console.SendString(tuitest.KeyDown)
			}
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return state.ForegroundColor(6, leftPadding+margin).String() == input.DefaultSelectedNameForeground
			})
			console.SendString(tuitest.KeyDown)
		})

		It("unselects the suggestion", func() {
			_, _ = console.WaitForDuration(func(state tuitest.TermState) bool {
				return state.ForegroundColor(6, leftPadding+margin).String() == input.DefaultNameForeground
			}, 100*time.Millisecond)
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

	When("the update function sends a periodic completer message", Ordered, func() {
		BeforeAll(func() {
			console.SendString("periodic")
			console.SendString(tuitest.KeyEnter)
		})

		It("Updates the completions", func() {
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return strings.Contains(state.NthOutputLine(3), "changed text2")
			})
		})
	})
})
