package util

func AsPointer[T any](s T) *T {
	return &s
}
