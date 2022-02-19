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
		model.Name.ForegroundColor = color
		return nil
	}
}

func OptionSelectedNameForegroundColor(color string) Option {
	return func(model *Model) error {
		model.Name.SelectedForegroundColor = color
		return nil
	}
}

func OptionNameBackgroundColor(color string) Option {
	return func(model *Model) error {
		model.Name.BackgroundColor = color
		return nil
	}
}

func OptionSelectedNameBackgroundColor(color string) Option {
	return func(model *Model) error {
		model.Name.SelectedBackgroundColor = color
		return nil
	}
}

func OptionNameFormatter(nameFormatter Formatter) Option {
	return func(model *Model) error {
		model.Name.Formatter = nameFormatter
		return nil
	}
}

func OptionDescriptionForegroundColor(color string) Option {
	return func(model *Model) error {
		model.Description.BackgroundColor = color
		return nil
	}
}

func OptionSelectedDescriptionForegroundColor(color string) Option {
	return func(model *Model) error {
		model.Description.SelectedBackgroundColor = color
		return nil
	}
}

func OptionDescriptionBackgroundColor(color string) Option {
	return func(model *Model) error {
		model.Description.BackgroundColor = color
		return nil
	}
}

func OptionSelectedDescriptionBackgroundColor(color string) Option {
	return func(model *Model) error {
		model.Description.SelectedBackgroundColor = color
		return nil
	}
}

func OptionDescriptionFormatter(descriptionFormatter Formatter) Option {
	return func(model *Model) error {
		model.Description.Formatter = descriptionFormatter
		return nil
	}
}

func OptionPlaceholderForegroundColor(color string) Option {
	return func(model *Model) error {
		model.Placeholder.ForegroundColor = color
		return nil
	}
}

func OptionPlaceholderBackgroundColor(color string) Option {
	return func(model *Model) error {
		model.Placeholder.BackgroundColor = color
		return nil
	}
}

func OptionPlaceholderFormatter(formatter func(text string) string) Option {
	return func(model *Model) error {
		model.Placeholder.Formatter = formatter
		return nil
	}
}

func OptionSelectedSuggestionForegroundColor(color string) Option {
	return func(model *Model) error {
		model.SelectedSuggestion.ForegroundColor = color
		return nil
	}
}

func OptionSelectedSuggestionBackgroundColor(color string) Option {
	return func(model *Model) error {
		model.SelectedSuggestion.BackgroundColor = color
		return nil
	}
}

func OptionSelectedSuggestionFormatter(formatter func(text string) string) Option {
	return func(model *Model) error {
		model.SelectedSuggestion.Formatter = formatter
		return nil
	}
}
