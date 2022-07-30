package gocess

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type parallel[T any] struct {
	steps  []Step[T]
	merger merger[T]
}

type merger[T any] func(ctx context.Context, array []T) T

// Execute will create go routines for each parallel step to take.
// an error on one or more of the steps will be blocking
func (p parallel[T]) Execute(ctx context.Context, input T) (T, error) {
	rets := make([]T, 0, len(p.steps))

	valueChannel := make(chan T)
	doneChannel := make(chan bool)

	go func() {
		for v := range valueChannel {
			rets = append(rets, v)
		}
		doneChannel <- true
	}()

	grp, c := errgroup.WithContext(ctx)
	c, cancel := context.WithCancel(c)
	defer cancel()

	for _, step := range p.steps {
		step := step
		grp.Go(func() error {
			v, err := step.Execute(c, input)
			if err != nil {
				return err
			}

			valueChannel <- v
			return nil
		})
	}

	err := grp.Wait()
	close(valueChannel)

	<-doneChannel
	close(doneChannel)

	if err != nil {
		return input, err
	}

	return p.merger(ctx, rets), nil
}

func ParallelSteps[T any](m merger[T], array ...Step[T]) Step[T] {
	return &parallel[T]{steps: array, merger: m}
}
