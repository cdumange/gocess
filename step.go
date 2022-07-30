package gocess

import "context"

type Step[T any] interface {
	Execute(ctx context.Context) (*T, error)
}
