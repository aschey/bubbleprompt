package test

import (
	"fmt"
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
