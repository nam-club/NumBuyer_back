package utils

func ContainsString(s []string, v string) bool {
	for _, a := range s {
		if a == v {
			return true
		}
	}
	return false
}
