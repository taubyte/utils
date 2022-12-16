package slices

/* Returns whether or not a provided []string contains a provided string value */
func Contains(slice []string, value string) bool {
	for _, s := range slice {
		if value == s {
			return true
		}
	}
	return false
}

/* Returns the reverse of a provides []string */
func ReverseArray(arr []string) (reversed []string) {
	reversed = make([]string, len(arr))
	last := len(arr) - 1
	for i, val := range arr {
		reversed[last-i] = val
	}
	return
}

/* Returns the last element of a provided []string */
func Last(arr []string) string {
	if len(arr) == 0 {
		return ""
	}

	if len(arr) == 1 {
		return arr[0]
	}

	return arr[len(arr)-1]
}
