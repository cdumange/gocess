package gocess

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/cdumange/gocess/pointer"
	"github.com/stretchr/testify/assert"
)

func simpleMerger[T any](ctx context.Context, array []*T) *T {
	if len(array) == 0 {
		return nil
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
			_, err := step.Execute(ctx, pointer.To("test"))
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
			_, err := step.Execute(ctx, pointer.To("test"))
			elapsed := time.Since(start)
			wg.Wait()
			assert.Error(t, err)
			assert.Less(t, elapsed, time.Millisecond*200)
		})
	})

	t.Run("merger", func(t *testing.T) {

		steps := []Step[string]{
			simpleValueErrorStep[string]{value: "1", err: nil},
			simpleValueErrorStep[string]{value: "2", err: nil},
			simpleValueErrorStep[string]{value: "3", err: nil},
		}

		e := ParallelSteps(simpleMerger[string], steps...)

		v, err := e.Execute(ctx, pointer.To("honeymoon"))
		assert.NoError(t, err)
		assert.NotNil(t, v)
		assert.Equal(t, "1", *v)
	})
}

type wgStep[T any] struct {
	wg  *sync.WaitGroup
	err error
}

func (m wgStep[T]) Execute(ctx context.Context, input *T) (*T, error) {
	time.Sleep(time.Millisecond * 100)
	m.wg.Done()
	return input, m.err
}

type simpleValueErrorStep[T any] struct {
	value T
	err   error
}

func (m simpleValueErrorStep[T]) Execute(ctx context.Context, input *T) (*T, error) {
	return &m.value, m.err
}
