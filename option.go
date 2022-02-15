package prompt

type Option func(model *Model) error

func OptionPrompt(prompt string) Option {
	return func(model *Model) error {
		model.SetPrompt(prompt)
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

func OptionNameBackgroundColor(color string) Option {
	return func(model *Model) error {
		model.NameBackgroundColor = color
		return nil
	}
}

func OptionNameFormatter(nameFormatter func(name string, columnWidth int) string) Option {
	return func(model *Model) error {
		model.NameFormatter = nameFormatter
		return nil
	}
}

func OptionDescriptionBackgroundColor(color string) Option {
	return func(model *Model) error {
		model.DescriptionBackgroundColor = color
		return nil
	}
}

func OptionDescriptionFormatter(descriptionFormatter func(name string, columnWidth int) string) Option {
	return func(model *Model) error {
		model.DescriptionFormatter = descriptionFormatter
		return nil
	}
}
