package selector

func equalIgnoreCase(a, b []rune) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] && a[i] != b[i]+32 && a[i] != b[i]-32 { // ASCII case-insensitive comparison
			return false
		}
	}
	return true
}
