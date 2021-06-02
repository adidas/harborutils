package util

func ApiVersion(s string) string {
	if s != "" {
		return s + "/"
	}
	return s
}
