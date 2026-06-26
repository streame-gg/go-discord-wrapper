package testutil

func RandomArray[V any](arraySize int, entries ...V) []V {
	if len(entries) == 0 {
		arr := make([]V, arraySize)

		var def V
		for i := 0; i < arraySize; i++ {
			arr = append(arr, def)
		}

		return arr
	}

	arr := make([]V, arraySize)
	for i := 0; i < arraySize; i++ {
		arr[i] = entries[RandomNumberInRange(0, len(entries)-1)]
	}

	return arr
}
