package test

import (
	"strings"

	"github.com/aschey/bubbleprompt/input/commandinput"
	tuitest "github.com/aschey/tui-tester"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Command Input", func() {
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

	When("the user presses the down arrow", Ordered, func() {
		BeforeAll(func() {
			console.SendString(tuitest.KeyDown)
		})

		It("selects the first prompt", func() {
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return strings.Contains(state.NthOutputLine(0), suggestions[0].Text)
			})
		})

		It("renders the placeholder", func() {
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return strings.Contains(state.NthOutputLine(0),
					suggestions[0].Text+" "+suggestions[0].Metadata.PositionalArgs[0].Placeholder+" "+suggestions[0].Metadata.PositionalArgs[1].Placeholder)
			})
		})

		It("applies the selected text styling", func() {
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return state.ForegroundColor(0, leftPadding).String() == commandinput.DefaultSelectedTextColor
			})
		})

		It("applies the selected placeholder styling", func() {
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return state.ForegroundColor(0, leftPadding+margin+len(suggestions[0].Text)).String() == commandinput.DefaultPlaceholderForeground
			})
		})

		When("the user presses the down arrow again", Ordered, func() {
			BeforeAll(func() {
				console.SendString(tuitest.KeyDown)
			})

			It("renders the updated suggestion with placeholders", func() {
				_, _ = console.WaitFor(func(state tuitest.TermState) bool {
					return strings.Contains(state.NthOutputLine(0),
						suggestions[1].Text+" "+suggestions[1].Metadata.PositionalArgs[0].Placeholder)
				})
			})
		})
	})

	When("the user types the full suggestion", Ordered, func() {
		BeforeAll(func() {
			console.SendString(suggestions[0].Text)
		})

		It("selects the suggestion", func() {
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return state.ForegroundColor(0, leftPadding).String() == commandinput.DefaultSelectedTextColor
			})
		})
	})

	When("the user views a subcommand suggestion", Ordered, func() {
		BeforeAll(func() {
			console.SendString(tuitest.KeyDown)
			console.SendString(tuitest.KeyDown)
			console.SendString(" ")
		})

		It("shows the subcommand suggestion", func() {
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return strings.Contains(state.NthOutputLine(0),
					suggestions[1].Text+" "+secondLevelSuggestions[0].Text+" "+secondLevelSuggestions[0].Metadata.PositionalArgs[0].Placeholder)
			})
		})

		It("shows the correct styles", func() {
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				firstSuggestionLen := leftPadding + len(suggestions[1].Text) + margin
				secondSuggestionLen := firstSuggestionLen + len(secondLevelSuggestions[0].Text) + margin

				return state.ForegroundColor(0, leftPadding).String() == commandinput.DefaultSelectedTextColor &&
					state.ForegroundColor(0, firstSuggestionLen).String() == commandinput.DefaultCurrentPlaceholderSuggestion &&
					state.ForegroundColor(0, secondSuggestionLen).String() == commandinput.DefaultPlaceholderForeground
			})
		})
	})
})
