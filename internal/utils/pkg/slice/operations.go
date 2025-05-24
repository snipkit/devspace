package slice

func ContainsInt(haystack []int, needle int) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

func ContainsString(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == "*" || s == needle {
			return true
		}
	}
	return false
}
