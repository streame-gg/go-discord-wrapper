package functions

func GetAs[T any](v any) T {
	return v.(T)
}
