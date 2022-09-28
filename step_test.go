package gocess

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewStepFromFunction(t *testing.T) {
	expected := 5

	f := func(ctx context.Context, _ int) (int, error) {
		return expected, nil
	}

	p := NewGoProcess(NewStepFromFunction(f))

	value, err := p.Execute(context.Background(), 0)
	assert.NoError(t, err)
	assert.Equal(t, expected, value)

}
