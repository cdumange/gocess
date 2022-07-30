package gocess

import (
	"context"
	"sync"

	"go.uber.org/multierr"
)

type parallel[T any] struct {
	steps []Step[T]
}

func (p parallel[T]) Execute(ctx context.Context, input *T) (*T, error) {
	var err, errs error

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

	for _, v := range p.steps {
		go func() {
			input, err = v.Execute(ctx, input)
			errChannel <- err
			wg.Done()
		}()
	}

	wg.Wait()
	close(errChannel)

	<-doneChannel
	close(doneChannel)

	return input, errs
}

func ParallelSteps[T any](array ...Step[T]) Step[T] {
	return &parallel[T]{steps: array}
}
