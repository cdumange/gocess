package gocess

import (
	"context"
	"sync"

	"go.uber.org/multierr"
)

type parallel[T any] struct {
	steps  []Step[T]
	merger merger[T]
}

type merger[T any] func(ctx context.Context, array []*T) *T

func (p parallel[T]) Execute(ctx context.Context, input *T) (*T, error) {
	var err, errs error

	rets := make([]*T, len(p.steps))

	errChannel := make(chan error)
	doneChannel := make(chan bool)

	go func() {
		for err := range errChannel {
			if err != nil {
				errs = multierr.Append(errs, err)
			}
		}
		doneChannel <- true
	}()

	wg := sync.WaitGroup{}
	wg.Add(len(p.steps))

	for i, v := range p.steps {
		go func(index int, s Step[T]) {
			input, err = s.Execute(ctx, input)
			errChannel <- err
			rets[index] = input
			wg.Done()
		}(i, v)
	}

	wg.Wait()
	close(errChannel)

	<-doneChannel
	close(doneChannel)

	return p.merger(ctx, rets), errs
}

func ParallelSteps[T any](m merger[T], array ...Step[T]) Step[T] {
	return &parallel[T]{steps: array, merger: m}
}
