package test

import (
	"strings"

	tuitest "github.com/aschey/tui-tester"
	. "github.com/onsi/ginkgo/v2"
	//. "github.com/onsi/gomega"
)

var _ = Describe("Parser", func() {
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
				return strings.Contains(state.NthOutputLine(1), "abcd")
			})
		})
	})

	When("the user selects a suggestion", Ordered, func() {
		BeforeAll(func() {
			console.SendString("a")
			console.SendString(tuitest.KeyTab)
		})

		It("updates the input", func() {
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return strings.Contains(state.NthOutputLine(0), "> abcd")
			})
		})
	})

	When("the user selects multiple tokens", Ordered, func() {
		BeforeAll(func() {
			console.SendString("a")
			console.SendString(tuitest.KeyTab)
			console.SendString(",")
			console.SendString("d")
			console.SendString(tuitest.KeyTab)
		})

		It("updates each token separately", func() {
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return strings.Contains(state.NthOutputLine(0), "> abcd,def")
			})
		})
	})

	When("the user types multiple tokens with whitespace", Ordered, func() {
		BeforeAll(func() {
			console.SendString(tuitest.KeyTab)
			console.SendString(" ")
			console.SendString(tuitest.KeyTab)
			console.SendString(tuitest.KeyTab)
		})

		It("updates the correct token", func() {
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return strings.Contains(state.NthOutputLine(0), "> abcd def")
			})
		})
	})

	When("the user types a token and a delimiter", Ordered, func() {
		BeforeAll(func() {
			console.SendString("a")
			console.SendString(tuitest.KeyTab)
			console.SendString(",")
		})

		It("suggests a new token", func() {
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return strings.Contains(state.NthOutputLine(1), "abcd")
			})
		})

		When("the user chooses the suggestion", Ordered, func() {
			BeforeAll(func() {
				console.SendString(tuitest.KeyTab)
			})

			It("updates the correct token", func() {
				_, _ = console.WaitFor(func(state tuitest.TermState) bool {
					return strings.Contains(state.NthOutputLine(0), "> abcd,abc")
				})
			})
		})
	})
	When("the user enters a suggestion between two delimiters", Ordered, func() {
		BeforeAll(func() {
			console.SendString("def,,abcd")
			for i := 0; i < 5; i++ {
				console.SendString(tuitest.KeyLeft)
			}
			console.SendString(tuitest.KeyTab)
			console.SendString(tuitest.KeyTab)
		})

		It("updates the correct token", func() {
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return strings.Contains(state.NthOutputLine(0), "> def,def,abcd")
			})
		})
	})

	When("the user enters a suggestion between two delimiters with spaces", Ordered, func() {
		BeforeAll(func() {
			console.SendString("def , , abcd")
			for i := 0; i < 7; i++ {
				console.SendString(tuitest.KeyLeft)
			}
			console.SendString(tuitest.KeyTab)
			console.SendString(tuitest.KeyTab)
		})

		It("updates the correct token", func() {
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return strings.Contains(state.NthOutputLine(0), "> def ,abcd , abcd")
			})
		})
	})

	When("the user enters invalid input", Ordered, func() {
		BeforeAll(func() {
			console.SendString("'")
			console.SendString(tuitest.KeyEnter)
		})

		It("displays the error message", func() {
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return strings.Contains(state.NthOutputLine(1), "invalid char literal")
			})
		})
	})
})
