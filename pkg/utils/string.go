package utils

func SafeString(s *string, defaultValue string) string {
	if s == nil {
		return defaultValue
	}
	return *s
}
