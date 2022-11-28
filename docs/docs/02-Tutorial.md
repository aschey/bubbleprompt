# Tutorial

Let's build a simple app to demonstratate how to use Bubbleprompt.
The app will display a list of fruits and tell the user which one they selected.
The final code can be seen in the [basic example](https://github.com/aschey/bubbleprompt/tree/main/examples/basic/main.go).

## Starting Out

First, we need to choose an input.
We'll use the simple input here because we don't need any fancy features like custom parsing or flags.
By default, the simple input parses input text as a series of whitespace-delimited tokens.
It also supports using double quotes to define a single token, so `"two words"` will be parsed as one token rather than two.

```go
func main() {
    textInput := simpleinput.New[any]()
}
```

The simple input component takes one generic parameter.
This parameter is used to define custom metadata that gets attached to each suggestion.
We don't need any custom metadata here so we'll leave it as `any`.

Next, we'll define a list of suggestions.
These will be shown underneath our input component.

```go
func main() {
    textInput := simpleinput.New[any]()

    suggestions := []input.Suggestion[any]{
        {Text: "banana", Description: "good with peanut butter"},
        {Text: "\"sugar apple\"", SuggestionText: "sugar apple", Description: "spherical...ish"},
        {Text: "jackfruit", Description: "the jack of all fruits"},
        {Text: "snozzberry", Description: "tastes like snozzberries"},
        {Text: "lychee", Description: "better than leeches"},
        {Text: "mangosteen", Description: "it's not a mango"},
        {Text: "durian", Description: "stinky"},
    }
}
```

The `Suggestion` struct defines each list entry that we show.
Here we're using three properties: `Text`, `Description`, and `SuggestionText`.

- `Text` represents the text that the user should type to choose this suggestion.
  It gets rendered on the left side of the suggestion.
- `Description` is an optional second property we can pass in to add some additional context to the suggestion.
  It gets rendered on the right side.
- `SuggestionText` is a special property that we can pass in to override the text that gets shown in the suggestion lists.

Here we're using the `SuggestionText` property for the second entry because it has two words, so we need to wrap it in quotes to treat it as a single token.
However, we don't want to show the quotes in the suggestion list because that would look odd.

Now, let's create a model. This will implement the `InputHandler` interface and hold our program state.
Additionally, we store a style struct from [lipgloss](https://github.com/charmbracelet/lipgloss) that we can use to add formatting to our output.

```go
type model struct {
    // list of suggestions that we'll display using the completer function
    suggestions []input.Suggestion[any]
    // Reference to our input component. We'll use this to read user input
    textInput   *simpleinput.Model[any]
    // Style struct for formatting the output
    outputStyle lipgloss.Style
    // Number of times the user enters some input
    numChoices  int64
}
```

Now we can create our model in our `main` function:

```go
func main() {
    // Initialize the input
    textInput := simpleinput.New[any]()

    // Define our suggestions
    suggestions := []input.Suggestion[any]{
        {Text: "banana", Description: "good with peanut butter"},
        {Text: "\"sugar apple\"", SuggestionText: "sugar apple", Description: "spherical...ish"},
        {Text: "jackfruit", Description: "the jack of all fruits"},
        {Text: "snozzberry", Description: "tastes like snozzberries"},
        {Text: "lychee", Description: "better than leeches"},
        {Text: "mangosteen", Description: "it's not a mango"},
        {Text: "durian", Description: "stinky"},
    }

    model := model{
        suggestions: suggestions,
        textInput:   textInput,
        // Add some coloring to the foreground of our output to make it look pretty
        outputStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("13")),
    }
```

## The Complete Method

In order to render our suggestions onto the screen, we need to define the `Complete` method.

```go
func (m model) Complete(promptModel prompt.Model[any]) ([]input.Suggestion[any], error) {
    // Our program only takes one token as input,
    // so don't return any suggestions if the user types more than one word
    if len(m.textInput.AllTokens()) > 1 {
        return nil, nil
    }

    // Filter suggestions based on the text before the cursor
    return completer.FilterHasPrefix(m.textInput.CurrentTokenBeforeCursor(), m.suggestions), nil
}
```

This method is responsible for returning a list of suggestions based on the user input.
Typically you'll have a predefined list of suggestions and you'll want to apply some kind of filtering function to replace the suggestions that aren't relevant to what the user typed.
Bubbleprompt provides a few predefined filtering functions in the `completer` package for convenience, but you're free to generate the list of suggestions however you want.

We use `simpleinput`'s `CurrentTokenBeforeCursor` method to get the text that the user typed before the cursor.
Since the list of suggestions always stays in sync with the cursor as it moves left or right,
it's expected that the completer function should only take into account what's before the cursor, rather than always checking the entire input.

## The Update Method

The `Update` method is part of the standard Bubbletea event loop.
It gets invoked whenever the program receives some kind of event.
See the [Bubbletea docs](https://github.com/charmbracelet/bubbletea/tree/master/tutorials/basics#the-update-method) for more information.

```go
func (m model) Update(msg tea.Msg) (prompt.InputHandler[any], tea.Cmd) {
    // Update the counter every time the user submits something
    if msg, ok := msg.(tea.KeyMsg); ok && msg.Type == tea.KeyEnter {
        m.numChoices++
    }
    return m, nil
}
```

Here we record every time the user presses enter so we can show this information later.

## The Execute Method

The executor method is invoked whenever the user presses enter.
It checks the user's input and returns a [tea.Model](https://github.com/charmbracelet/bubbletea/tree/master/tutorials/basics#the-model) that gets rendered to the output.
The returned model will take over the event loop until it finishes, and then we can start the process over.

```go
func (m model) Execute(input string, promptModel *prompt.Model[any]) (tea.Model, error) {
    // Get a list of all the tokens from the input
    tokens := m.textInput.TokenValues()
    if len(tokens) == 0 {
        // We didn't receive any input, which is invalid
        // Returning an error will output text will special error styling
        return nil, fmt.Errorf("No selection")
    }
    // The user entered a selection
    // Render their choice with styling applied
    return executor.NewStringModel(m.formatOutput(tokens[0])), nil
}

func (m model) formatOutput(choice string) string {
    return fmt.Sprintf("You picked: %s\nYou've entered %s submissions(s)\n\n",
        m.outputStyle.Render(choice),
        m.outputStyle.Render(strconv.FormatInt(m.numChoices, 10)))
}

```

Here we check if the user entered in any input and display their choice if they did.
The executor method requires that we return a `tea.Model`, but it would be rather annoying to have to
manually create a new model for simple cases like showing a line of text.
For these cases, the `executor` package supplies several prebuilt models for common situations.

## Putting It All Together

Now that we have all the building blocks, we can finish writing our `main` function.

```go
func main() {
    // Initialize the input
    textInput := simpleinput.New[any]()

    // Define our suggestions
    suggestions := []input.Suggestion[any]{
        {Text: "banana", Description: "good with peanut butter"},
        {Text: "\"sugar apple\"", SuggestionText: "sugar apple", Description: "spherical...ish"},
        {Text: "jackfruit", Description: "the jack of all fruits"},
        {Text: "snozzberry", Description: "tastes like snozzberries"},
        {Text: "lychee", Description: "better than leeches"},
        {Text: "mangosteen", Description: "it's not a mango"},
        {Text: "durian", Description: "stinky"},
    }

    // Combine everything into our model
    model := model{
        suggestions: suggestions,
        textInput:   textInput,
        // Add some coloring to the foreground of our output to make it look pretty
        outputStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("13")),
    }

    // Create the Bubbleprompt model
    // This struct fulfills the tea.Model interface so it can be passed directly to tea.NewProgram
    promptModel, err := prompt.New[any](modeltextInput)
    if err != nil {
        panic(err)
    }

    fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render("Pick a fruit!"))
    fmt.Println()

    if _, err := tea.NewProgram(promptModel, tea.WithFilter(prompt.MsgFilter)).Run(); err != nil {
        fmt.Printf("Could not start program\n%v\n", err)
        os.Exit(1)
    }
}
```

With that in place, everything should be functional.
This is what the whole program looks like:

## Complete Program

```go
package main

import (
    "fmt"
    "os"
    "strconv"

    prompt "github.com/aschey/bubbleprompt"
    "github.com/aschey/bubbleprompt/completer"
    "github.com/aschey/bubbleprompt/input"
    "github.com/aschey/bubbleprompt/input/simpleinput"
    "github.com/aschey/bubbleprompt/executor"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

type model struct {
    suggestions []input.Suggestion[any]
    textInput   *simpleinput.Model[any]
    outputStyle lipgloss.Style
    numChoices  int64
}

func (m model) Complete(promptModel prompt.Model[any]) ([]input.Suggestion[any], error) {
    if len(m.textInput.AllTokens()) > 1 {
        return nil, nil
    }

    return completer.FilterHasPrefix(m.textInput.CurrentTokenBeforeCursor(), m.suggestions), nil
}

func (m model) Execute(input string, promptModel *prompt.Model[any]) (tea.Model, error) {
    tokens := m.textInput.TokenValues()
    if len(tokens) == 0 {
        return nil, fmt.Errorf("No selection")
    }
    return executor.NewStringModel(m.formatOutput(tokens[0])), nil
}

func (m model) formatOutput(choice string) string {
    return fmt.Sprintf("You picked: %s\nYou've entered %s submissions(s)\n\n",
        m.outputStyle.Render(choice),
        m.outputStyle.Render(strconv.FormatInt(m.numChoices, 10)))
}

func (m model) Update(msg tea.Msg) (prompt.InputHandler[any], tea.Cmd) {
    if msg, ok := msg.(tea.KeyMsg); ok && msg.Type == tea.KeyEnter {
        m.numChoices++
    }
    return m, nil
}

func main() {
    textInput := simpleinput.New[any]()
    suggestions := []input.Suggestion[any]{
        {Text: "banana", Description: "good with peanut butter"},
        {Text: "\"sugar apple\"", SuggestionText: "sugar apple", Description: "spherical...ish"},
        {Text: "jackfruit", Description: "the jack of all fruits"},
        {Text: "snozzberry", Description: "tastes like snozzberries"},
        {Text: "lychee", Description: "better than leeches"},
        {Text: "mangosteen", Description: "it's not a mango"},
        {Text: "durian", Description: "stinky"},
    }

    model := model{
        suggestions: suggestions,
        textInput:   textInput,
        outputStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("13")),
    }

    promptModel, err := prompt.New[any](model,textInput)
    if err != nil {
        panic(err)
    }

    fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render("Pick a fruit!"))
    fmt.Println()

    if _, err := tea.NewProgram(promptModel, tea.WithFilter(prompt.MsgFilter)).Run(); err != nil {
        fmt.Printf("Could not start program\n%v\n", err)
        os.Exit(1)
    }
}
```
