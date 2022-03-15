package commandinput

import "regexp"

type Option func(model *Model) error

func WithPrompt(prompt string) Option {
	return func(model *Model) error {
		model.SetPrompt(prompt)
		return nil
	}
}

func WithDelimiterRegex(delimiterRegex *regexp.Regexp) Option {
	return func(model *Model) error {
		model.SetDelimiterRegex(delimiterRegex)
		return nil
	}
}

func WithStringRegex(stringRegex *regexp.Regexp) Option {
	return func(model *Model) error {
		model.SetStringRegex(stringRegex)
		return nil
	}
}
