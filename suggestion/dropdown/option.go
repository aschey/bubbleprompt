package dropdown

import "github.com/aschey/bubbleprompt/suggestion"

type Option[T any] func(model *Model[T])

func WithMaxSuggestions[T any](maxSuggestions int) Option[T] {
	return func(model *Model[T]) {
		model.SetMaxSuggestions(maxSuggestions)
	}
}

func WithSelectionIndicator[T any](indicator string) Option[T] {
	return func(model *Model[T]) {
		model.SetSelectionIndicator(indicator)
	}
}

func WithFormatters[T any](formatters suggestion.Formatters) Option[T] {
	return func(model *Model[T]) {
		model.SetFormatters(formatters)
	}
}
