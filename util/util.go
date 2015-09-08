package util

func Keys(words map[string][]string) []string {
	var keys []string

	for key, _ := range words {
		keys = append(keys, key)
	}
	return keys
}
