package gocess

import "context"

type Process[T any] interface {
	Execute(ctx context.Context) (T, error)
}
