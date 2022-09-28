package gocess

import "context"

// Step is a single step in a process.
type Step[T any] interface {
	Execute(ctx context.Context, input T) (T, error)
}

type emptyStep[T any] struct {
	f func(ctx context.Context, input T) (T, error)
}

func (e *emptyStep[T]) Execute(ctx context.Context, input T) (T, error) {
	return e.f(ctx, input)
}

// NewStepFromFunction allows to create an empty step structure to allow execution of simple steps as Step.
// will execute passed f function when called.
func NewStepFromFunction[T any](f func(ctx context.Context, input T) (T, error)) Step[T] {
	return &emptyStep[T]{f}
}
