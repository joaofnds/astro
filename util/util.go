package util

func Min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func Max(a, b int) int {
	if a < b {
		return b
	}

	return a
}
