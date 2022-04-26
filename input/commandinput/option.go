package commandinput

import "regexp"

type Option[T CmdMetadataAccessor] func(model *Model[T]) error

func WithPrompt[T CmdMetadataAccessor](prompt string) Option[T] {
	return func(model *Model[T]) error {
		model.SetPrompt(prompt)
		return nil
	}
}

func WithDelimiterRegex[T CmdMetadataAccessor](delimiterRegex *regexp.Regexp) Option[T] {
	return func(model *Model[T]) error {
		model.SetDelimiterRegex(delimiterRegex)
		return nil
	}
}

func WithDefaultDelimiter[T CmdMetadataAccessor](defaultDelimiter string) Option[T] {
	return func(model *Model[T]) error {
		model.defaultDelimiter = defaultDelimiter
		return nil
	}
}

func WithStringRegex[T CmdMetadataAccessor](stringRegex *regexp.Regexp) Option[T] {
	return func(model *Model[T]) error {
		model.SetStringRegex(stringRegex)
		return nil
	}
}
