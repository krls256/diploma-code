package utils

func SimpleMapDeepCopy[K, V comparable](m map[K]V) map[K]V {
	nm := map[K]V{}

	for k, v := range m {
		nm[k] = v
	}

	return nm
}
