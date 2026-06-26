package testutil

func RandomItem[V any](items ...V) V {
	if len(items) == 0 {
		panic("no items provided to RandomItem")
	}

	randomItem := RandomNumberInRange(0, len(items)-1)
	return items[randomItem]
}
