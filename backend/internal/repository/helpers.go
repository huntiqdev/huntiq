package repository

import (
	"fmt"
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

// whereCond represents a single optional WHERE condition.
type whereCond struct {
	ok    bool
	value interface{}
	expr  string // static clause template, e.g. "event=$1"
}

// buildArgs appends values of active conditions to args and returns the updated slice.
// The caller is responsible for building the WHERE string using static clause strings.
func buildArgs(conds ...whereCond) []interface{} {
	var args []interface{}
	for _, c := range conds {
		if c.ok && c.value != nil {
			args = append(args, c.value)
		}
	}
	return args
}

// buildWhere builds a parameterized WHERE clause from a list of conditions.
// Each condition's value is appended to args; $? in expr is replaced with the
// correct $N index. Returns args and the WHERE string (empty if no conditions).
func buildWhere(conds ...whereCond) ([]interface{}, string) {
	var args []interface{}
	var clauses []string
	for _, c := range conds {
		if !c.ok || c.value == nil {
			continue
		}
		args = append(args, c.value)
		clauses = append(clauses, strings.ReplaceAll(c.expr, "$?", fmt.Sprintf("$%d", len(args)))) //nolint:gosec // G201: Building static SQL fragment with computed $N parameter placeholder
	}
	if len(clauses) == 0 {
		return args, ""
	}
	return args, " WHERE " + strings.Join(clauses, " AND ")
}

// paginate returns " ORDER BY ... LIMIT $limitIdx OFFSET $offsetIdx" using
// pure integer indices — no user data ever reaches this string.
func paginate(limitIdx, offsetIdx int) string {
	return fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", limitIdx, offsetIdx) //nolint:gosec // G201: Building static SQL with computed $N placeholders
}
