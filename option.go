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
		model.NameForegroundColor = color
		return nil
	}
}

func OptionNameBackgroundColor(color string) Option {
	return func(model *Model) error {
		model.NameBackgroundColor = color
		return nil
	}
}

func OptionNameFormatter(nameFormatter Formatter) Option {
	return func(model *Model) error {
		model.NameFormatter = nameFormatter
		return nil
	}
}

func OptionDescriptionForegroundColor(color string) Option {
	return func(model *Model) error {
		model.DescriptionBackgroundColor = color
		return nil
	}
}

func OptionDescriptionBackgroundColor(color string) Option {
	return func(model *Model) error {
		model.DescriptionBackgroundColor = color
		return nil
	}
}

func OptionDescriptionFormatter(descriptionFormatter Formatter) Option {
	return func(model *Model) error {
		model.DescriptionFormatter = descriptionFormatter
		return nil
	}
}
