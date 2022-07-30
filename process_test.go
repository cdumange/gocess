package gocess

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGoProcess_Execute(t *testing.T) {
	ctx := context.Background()
	t.Run("honeymoon", func(t *testing.T) {
		input := "honeymoon"
		processes := []Step[string]{
			&timedReturnedStep[string]{duration: time.Millisecond},
			&timedReturnedStep[string]{duration: time.Millisecond},
			&timedReturnedStep[string]{duration: time.Millisecond},
		}

		p := NewGoProcess[string](processes...)
		v, err := p.Execute(ctx, input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		assert.NotNil(t, v)
		assert.Equal(t, input, v)

		for _, p := range processes {
			assert.True(t, p.(*timedReturnedStep[string]).Done)
		}
	})

	t.Run("contextCancelled", func(t *testing.T) {
		input := "honeymoon"
		c, cancel := context.WithCancel(ctx)

		processes := []Step[string]{
			&timedReturnedStep[string]{duration: time.Millisecond * 200},
			&timedReturnedStep[string]{duration: time.Millisecond},
		}

		p := NewGoProcess[string](processes...)
		go func() {
			time.Sleep(time.Millisecond * 100)
			cancel()
		}()
		_, err := p.Execute(c, input)

		assert.ErrorIs(t, err, context.Canceled)

		assert.True(t, processes[0].(*timedReturnedStep[string]).Done)
		assert.False(t, processes[1].(*timedReturnedStep[string]).Done)
	})
}

type timedReturnedStep[T any] struct {
	duration time.Duration
	Done     bool
}

func (s *timedReturnedStep[T]) Execute(_ context.Context, input T) (T, error) {
	time.Sleep(s.duration)
	s.Done = true

	return input, nil
}
