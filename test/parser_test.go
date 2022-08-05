package test

import (
	"strings"

	tuitest "github.com/aschey/tui-tester"
	. "github.com/onsi/ginkgo/v2"
	//. "github.com/onsi/gomega"
)

var _ = Describe("Prompt", func() {
	var console *tuitest.Console

	BeforeEach(OncePerOrdered, func() {
		console, _ = parserTester.CreateConsole()

		// Wait for prompt to initialize
		console.TrimOutput = true
		_, _ = console.WaitFor(func(state tuitest.TermState) bool {
			return strings.Contains(state.NthOutputLine(2), "def")
		})
	})

	AfterEach(OncePerOrdered, func() {
		console.SendString(tuitest.KeyCtrlC)
		_ = console.WaitForTermination()
	})

	When("the user types a token", Ordered, func() {
		BeforeAll(func() {
			console.SendString("d")
		})

		It("filters the completions", func() {
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return strings.Contains(state.NthOutputLine(1), "def")
			})
		})
	})

	When("the user types multiple tokens", Ordered, func() {
		BeforeAll(func() {
			console.SendString("def.a")
		})

		It("filters the completions", func() {
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return strings.Contains(state.NthOutputLine(1), "abc")
			})
		})
	})
})
