package gocess

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func simpleMerger[T any](_ context.Context, array []T) T {
	if len(array) == 0 {
		return *new(T)
	}
	return array[0]
}

func TestParallel_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("test parallelism", func(t *testing.T) {
		t.Run("honeymoon", func(t *testing.T) {
			nb := 3
			wg := sync.WaitGroup{}
			steps := make([]Step[string], 0, nb)

			for n := 0; n < nb; n++ {
				wg.Add(1)
				steps = append(steps, wgStep[string]{wg: &wg})
			}

			step := ParallelSteps(simpleMerger[string], steps...)

			start := time.Now()
			_, err := step.Execute(ctx, "test")
			elapsed := time.Since(start)
			wg.Wait()
			assert.NoError(t, err)
			assert.Less(t, elapsed, time.Millisecond*200)
		})

		t.Run("err", func(t *testing.T) {
			nb := 3
			wg := sync.WaitGroup{}
			steps := make([]Step[string], 0, nb)

			for n := 0; n < nb; n++ {
				wg.Add(1)
				steps = append(steps, wgStep[string]{wg: &wg, err: errors.New("an error")})
			}

			step := ParallelSteps(simpleMerger[string], steps...)

			start := time.Now()
			_, err := step.Execute(ctx, "test")
			elapsed := time.Since(start)
			wg.Wait()
			assert.Error(t, err)
			assert.Less(t, elapsed, time.Millisecond*200)
		})
	})

	t.Run("merger", func(t *testing.T) {

		values := []string{"1", "2", "3"}
		steps := []Step[string]{
			simpleValueErrorStep[string]{value: values[0], err: nil},
			simpleValueErrorStep[string]{value: values[1], err: nil},
			simpleValueErrorStep[string]{value: values[2], err: nil},
		}

		e := ParallelSteps(simpleMerger[string], steps...)

		v, err := e.Execute(ctx, "honeymoon")
		assert.NoError(t, err)
		assert.NotNil(t, v)

		assert.Contains(t, values, v)
	})

	t.Run("errs", func(t *testing.T) {
		m := new(MockStep[string])
		m.On("Execute", mock.Anything, mock.Anything).
			Return("", errors.New("an error")).
			Times(3)

		s := Step[string](m)
		e := ParallelSteps(simpleMerger[string], s, s, s)
		_, err := e.Execute(ctx, "an input")

		require.Error(t, err)

		m.AssertExpectations(t)
	})
}

type wgStep[T any] struct {
	wg  *sync.WaitGroup
	err error
}

func (m wgStep[T]) Execute(_ context.Context, input T) (T, error) {
	time.Sleep(time.Millisecond * 100)
	m.wg.Done()
	return input, m.err
}

type simpleValueErrorStep[T any] struct {
	value T
	err   error
}

func (m simpleValueErrorStep[T]) Execute(_ context.Context, _ T) (T, error) {
	return m.value, m.err
}
