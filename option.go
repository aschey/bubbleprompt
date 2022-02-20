package prompt

type Option func(model *Model) error

func OptionPrompt(prompt string) Option {
	return func(model *Model) error {
		model.SetPrompt(prompt)
		return nil
	}
}

func OptionNameForegroundColor(color string) Option {
	return func(model *Model) error {
		model.Formatters.Name.ForegroundColor = color
		return nil
	}
}

func OptionSelectedNameForegroundColor(color string) Option {
	return func(model *Model) error {
		model.Formatters.Name.SelectedForegroundColor = color
		return nil
	}
}

func OptionNameBackgroundColor(color string) Option {
	return func(model *Model) error {
		model.Formatters.Name.BackgroundColor = color
		return nil
	}
}

func OptionSelectedNameBackgroundColor(color string) Option {
	return func(model *Model) error {
		model.Formatters.Name.SelectedBackgroundColor = color
		return nil
	}
}

func OptionNameFormatter(nameFormatter Formatter) Option {
	return func(model *Model) error {
		model.Formatters.Name.Formatter = nameFormatter
		return nil
	}
}

func OptionDescriptionForegroundColor(color string) Option {
	return func(model *Model) error {
		model.Formatters.Description.BackgroundColor = color
		return nil
	}
}

func OptionSelectedDescriptionForegroundColor(color string) Option {
	return func(model *Model) error {
		model.Formatters.Description.SelectedBackgroundColor = color
		return nil
	}
}

func OptionDescriptionBackgroundColor(color string) Option {
	return func(model *Model) error {
		model.Formatters.Description.BackgroundColor = color
		return nil
	}
}

func OptionSelectedDescriptionBackgroundColor(color string) Option {
	return func(model *Model) error {
		model.Formatters.Description.SelectedBackgroundColor = color
		return nil
	}
}

func OptionDescriptionFormatter(descriptionFormatter Formatter) Option {
	return func(model *Model) error {
		model.Formatters.Description.Formatter = descriptionFormatter
		return nil
	}
}

func OptionPlaceholderForegroundColor(color string) Option {
	return func(model *Model) error {
		model.Formatters.Placeholder.ForegroundColor = color
		return nil
	}
}

func OptionPlaceholderBackgroundColor(color string) Option {
	return func(model *Model) error {
		model.Formatters.Placeholder.BackgroundColor = color
		return nil
	}
}

func OptionPlaceholderFormatter(formatter func(text string) string) Option {
	return func(model *Model) error {
		model.Formatters.Placeholder.Formatter = formatter
		return nil
	}
}

func OptionSelectedSuggestionForegroundColor(color string) Option {
	return func(model *Model) error {
		model.Formatters.SelectedSuggestion.ForegroundColor = color
		return nil
	}
}

func OptionSelectedSuggestionBackgroundColor(color string) Option {
	return func(model *Model) error {
		model.Formatters.SelectedSuggestion.BackgroundColor = color
		return nil
	}
}

func OptionSelectedSuggestionFormatter(formatter func(text string) string) Option {
	return func(model *Model) error {
		model.Formatters.SelectedSuggestion.Formatter = formatter
		return nil
	}
}
