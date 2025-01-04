package helpers

func FromPtr[T any](t *T) T {
	return *t
}

func ToPtr[T any](t T) *T {
	return &t
}
