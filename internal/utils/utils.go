package utils

import (
	"strings"
)

func CleanFilename(s string) string {
	var result strings.Builder
	allowed := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-. "

	for _, c := range s {
		if strings.ContainsRune(allowed, c) {
			result.WriteRune(c)
		}
	}

	return result.String()
}
