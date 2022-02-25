package main

func InArray[T any](needle *T, haystack []*T) bool {
	for _, e := range haystack {
		if e == needle {
			return true
		}
	}

	return false
}

func ArrayMap[T any](haystack []*T, f func(index int, element *T)) []*T {
	for i, e := range haystack {
		f(i, e)
	}

	return haystack
}
