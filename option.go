package prompt

type Option func(model *Model) error

func OptionPrompt(prompt string) Option {
	return func(model *Model) error {
		model.textInput.Prompt = prompt
		return nil
	}
}

func OptionInitialSuggestions(suggestions []Suggest) Option {
	return func(model *Model) error {
		model.suggestions = suggestions
		model.filteredSuggestions = suggestions
		return nil
	}
}
