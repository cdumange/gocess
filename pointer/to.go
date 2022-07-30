package pointer

func To[T any](i T) *T {
	return &i
}
