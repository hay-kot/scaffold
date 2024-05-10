package scaffold

func parseDefaultString(value ...any) string {
	for _, val := range value {
		if str, ok := val.(string); ok {
			return str
		}
	}

	return ""
}

func parseDefaultStrings(value ...any) []string {
	for _, val := range value {
		out := make([]string, 0)
		// likely to be []interface{} so we need to cast each value
		// to a string
		if arr, ok := val.([]interface{}); ok {
			for _, v := range arr {
				if str, ok := v.(string); ok {
					out = append(out, str)
				}
			}
		}

		if len(out) > 0 {
			return out
		}
	}

	return nil
}

func parseDefaultBool(value ...any) bool {
	for _, val := range value {
		if b, ok := val.(bool); ok {
			return b
		}
	}

	return false
}
