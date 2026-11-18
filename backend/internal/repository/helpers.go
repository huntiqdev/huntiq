package repository

import (
	"strings"
)

// escapeLike escapes special characters in LIKE patterns
func escapeLike(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "%", "\\%")
	s = strings.ReplaceAll(s, "_", "\\_")
	return s
}

// likeWrap wraps a string for LIKE pattern matching
func likeWrap(s string) string {
	return "%" + escapeLike(s) + "%"
}
