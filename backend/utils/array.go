package utils

func Includes(a []string, t string) bool {
	for _, s := range a {
		if s == t {
			return true
		}
	}
	return false
}
