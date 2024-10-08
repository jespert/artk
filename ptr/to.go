package ptr

// To returns a pointer to the specified value.
// This function never returns a nil pointer.
//
// It is a quality of life workaround for Go's inability to do things like:
//
//	x := &int(42)
func To[T any](value T) *T {
	return &value
}
