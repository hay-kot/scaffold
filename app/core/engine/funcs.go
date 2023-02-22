package engine

// wraptmpl wraps a string in the template delimiters.
func wraptmpl(s string) string {
	return "{{ " + s + " }}"
}
