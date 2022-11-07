package test

import (
	"strings"
	"time"

	"github.com/aschey/bubbleprompt/input/commandinput"
	tuitest "github.com/aschey/tui-tester"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Filter", func() {
	var console *tuitest.Console
	textInput := commandinput.New[cmdMetadata]()
	suggestions := suggestions(textInput)

	BeforeEach(OncePerOrdered, func() {
		console, _ = cmdTester.CreateConsole()

		// Wait for prompt to initialize
		console.TrimOutput = true
		_, _ = console.WaitFor(func(state tuitest.TermState) bool {
			return strings.Contains(state.NthOutputLine(6), suggestions[5].Description)
		})
	})

	AfterEach(OncePerOrdered, func() {
		console.SendString(tuitest.KeyCtrlC)
		_ = console.WaitForTermination()
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
				return state.ForegroundColor(0, 2).Int() == tuitest.DefaultFG
			}, 100*time.Millisecond)
		})
	})

	When("the user filters a suggestion with completion text override", Ordered, func() {
		BeforeAll(func() {
			console.SendString("completion")
		})

		It("filters the suggestions", func() {
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return strings.Contains(state.NthOutputLine(1), suggestions[6].CompletionText)
			})
		})

		When("the user chooses the suggestion", Ordered, func() {
			BeforeAll(func() {
				console.SendString(tuitest.KeyDown)
			})

			It("shows the regular text in the input", func() {
				_, _ = console.WaitFor(func(state tuitest.TermState) bool {
					return strings.Contains(state.NthOutputLine(0), suggestions[6].Text)
				})
			})
		})

	})

})
