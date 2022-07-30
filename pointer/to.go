package pointer

// To will return a pointer to passed value
func To[T any](i T) *T {
	return &i
}
