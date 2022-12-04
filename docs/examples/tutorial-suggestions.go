package tutorial

import (
	"github.com/aschey/bubbleprompt/input/simpleinput"
	"github.com/aschey/bubbleprompt/suggestion"
)

func main2() {
	textInput := simpleinput.New[any]()

	suggestions := []suggestion.Suggestion[any]{
		{Text: "banana", Description: "good with peanut butter"},
		{Text: "\"sugar apple\"", SuggestionText: "sugar apple", Description: "spherical...ish"},
		{Text: "jackfruit", Description: "the jack of all fruits"},
		{Text: "snozzberry", Description: "tastes like snozzberries"},
		{Text: "lychee", Description: "better than leeches"},
		{Text: "mangosteen", Description: "it's not a mango"},
		{Text: "durian", Description: "stinky"},
	}
	_ = textInput
	_ = suggestions
}
