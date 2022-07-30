package gocess

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockStep a mock struct for Step[T]
type MockStep[T any] struct {
	mock.Mock
}

// Execute - duh ?
func (m *MockStep[T]) Execute(ctx context.Context, input T) (T, error) {
	called := m.Called(ctx, input)
	return called.Get(0).(T), called.Error(1)
}
