package prompt

import (
	"github.com/aschey/bubbleprompt/formatter"
	"github.com/aschey/bubbleprompt/renderer"
	"github.com/aschey/bubbleprompt/suggestion"
)

type Option[T any] func(model *Model[T])

func WithFormatters[T any](formatters formatter.Formatters) Option[T] {
	return func(model *Model[T]) {
		model.formatters = formatters
	}
}

func WithUnmanagedRenderer[T any](opts ...renderer.Option) Option[T] {
	return func(model *Model[T]) {
		model.renderer = renderer.NewUnmanagedRenderer(opts...)
	}
}

func WithViewportRenderer[T any](opts ...renderer.Option) Option[T] {
	return func(model *Model[T]) {
		model.renderer = renderer.NewViewportRenderer(opts...)
	}
}

func WithRenderer[T any](r renderer.Renderer) Option[T] {
	return func(model *Model[T]) {
		model.renderer = r
	}
}

func WithSuggestionManager[T any](manager suggestion.Manager[T]) Option[T] {
	return func(model *Model[T]) {
		model.suggestionManager = manager
	}
}
