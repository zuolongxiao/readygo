package utils

// StrInSlice find a string in the slice
func StrInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
