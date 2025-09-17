package utils

import (
	"regexp"
	"strings"
)

func CleanFilename(s string) string {
	// Regex explanation:
	// \p{L}  : any Unicode letter (including accented chars like é, ñ, ü, etc.)
	// \p{N}  : any Unicode number
	// _ . - ! ? and space : explicitly allowed symbols
	re := regexp.MustCompile(`[\p{L}\p{N}_.\-!?¡¿ ]+`)

	// Find all allowed substrings and join them
	matches := re.FindAllString(s, -1)
	return strings.Join(matches, "")
}
