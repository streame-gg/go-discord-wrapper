package testutil

func RandomItem[V any](items ...V) V {
	if len(items) == 0 {
		panic("no items provided to RandomItem")
	}
	if len(items) == 1 {
		return items[0]
	}

	randomItem := RandomIntInRange(0, len(items)-1)
	return items[randomItem]
}
