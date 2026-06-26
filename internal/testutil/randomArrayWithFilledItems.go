package testutil

func RandomArrayWithFilledItems[V any](arraySize int, fillFunction func(arrayToFill *[]V)) []V {
	arr := make([]V, 0, arraySize)
	for i := 0; i < arraySize; i++ {
		fillFunction(&arr)
	}
	return arr
}
