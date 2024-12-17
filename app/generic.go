package app

// look if a key in a map exists, generic variant
func Exists[K comparable, V any](m map[K]V, v K) bool {
	if _, ok := m[v]; ok {
		return true
	}

	return false
}
