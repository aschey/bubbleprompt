package test

import (
	"strings"
	"time"

	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/input/commandinput"
	tuitest "github.com/aschey/tui-tester"
	. "github.com/onsi/ginkgo/v2"
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

var _ = Describe("Executor", func() {
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

	When("the user chooses the first prompt", Ordered, func() {
		BeforeAll(func() {
			console.SendString(tuitest.KeyDown)
			console.SendString(tuitest.KeyEnter)
		})

		It("renders the executor result", func() {
			_, _ = console.WaitForDuration(func(state tuitest.TermState) bool {
				secondLine := state.NthOutputLine(1)
				return strings.Contains(secondLine, "selected suggestion is "+suggestions[0].Text) &&
					!strings.Contains(secondLine, suggestions[0].Metadata.GetPositionalArgs()[0].Placeholder())
			}, 100*time.Millisecond)
		})
	})

	When("the executor returns an error", Ordered, func() {
		BeforeAll(func() {
			console.SendString("error")
			console.SendString(tuitest.KeyEnter)
		})

		It("displays the error message", func() {
			_, _ = console.WaitFor(func(state tuitest.TermState) bool {
				return strings.Contains(state.NthOutputLine(0), "> error") &&
					strings.Contains(state.NthOutputLine(1), "bad things") &&
					state.BackgroundColor(1, 0).String() == input.DefaultErrorTextBackground
			})
		})
	})
})
