# Usage

```go
func (m model) Complete(promptModel prompt.Model[any]) ([]input.Suggestion[any], error) {
	if len(m.textInput.AllTokens()) > 1 {
		return nil, nil
	}

	return completer.FilterHasPrefix(m.textInput.CurrentTokenBeforeCursor(), m.suggestions), nil
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
        textInput: textInput,
        suggestions: suggestions,
    }

    promptModel, err := prompt.New[any](model, textInput)
}

```
