package gocess

import "context"

// Process is a collection of steps.
type Process[T any] interface {
	Execute(ctx context.Context, input T) (T, error)
}

// GoProcess is a bunch of steps that are executed in sequence.
type GoProcess[T any] struct {
	Steps []Step[T]
}

// NewGoProcess creates a new GoProcess.
func NewGoProcess[T any](steps ...Step[T]) *GoProcess[T] {
	return &GoProcess[T]{
		Steps: steps,
	}
}

// Execute will execute the steps in the process.
func (p *GoProcess[T]) Execute(ctx context.Context, input T) (T, error) {
	var err error
	for _, step := range p.Steps {
		if ctx.Err() != nil {
			return input, ctx.Err()
		}
		if input, err = step.Execute(ctx, input); err != nil {
			return input, err
		}
	}
	return input, nil
}

// AddStep adds a step to a defined GoProcess
func (p *GoProcess[T]) AddStep(step Step[T]) {
	p.Steps = append(p.Steps, step)
}
