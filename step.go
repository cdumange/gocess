package gocess

import "context"

// Step is a single step in a process.
type Step[T any] interface {
	Execute(ctx context.Context, input T) (T, error)
}
