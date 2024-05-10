package scaffold

// MergeMaps merges multiple maps into a single map. If a key is present in
// multiple maps, the value from the last map will be used.
func MergeMaps[T any](maps ...map[string]T) map[string]T {
	out := map[string]T{}
	for _, m := range maps {
		for k, v := range m {
			out[k] = v
		}
	}
	return out
}

// CastMap casts a map[string]T to a man[string]any map.
func CastMap[T any](m map[string]T) map[string]any {
	out := map[string]any{}
	for k, v := range m {
		out[k] = v
	}
	return out
}
