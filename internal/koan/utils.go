package koan

import (
	"strings"
)

// replaceBlank replaces the blank placeholder (__) with the answer
func replaceBlank(incomplete, answer string) string {
	return strings.Replace(incomplete, "__", answer, 1)
}

// normalizeAnswer trims spaces and converts to lowercase for comparison
func normalizeAnswer(answer string) string {
	return strings.TrimSpace(strings.ToLower(answer))
}
