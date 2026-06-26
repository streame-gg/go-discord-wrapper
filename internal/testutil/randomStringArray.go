package testutil

func RandomStringArray(arraySize, stringMinLength, stringMaxLength int) []string {
	arr := make([]string, 0, arraySize)

	for i := 0; i < arraySize; i++ {
		arr = append(arr, RandString(RandomNumberInRange(stringMinLength, stringMaxLength)))
	}

	return arr
}
