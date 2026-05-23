// Package validator provides utilities for validating environment variable
// keys and values loaded from .env files.
package validator

import (
	"fmt"
	"strings"
	"unicode"
)

// Issue represents a single validation problem found in an env map.
type Issue struct {
	Key     string
	Message string
}

func (i Issue) Error() string {
	return fmt.Sprintf("key %q: %s", i.Key, i.Message)
}

// Result holds all issues found during validation.
type Result struct {
	Issues []Issue
}

// OK returns true when no issues were found.
func (r *Result) OK() bool { return len(r.Issues) == 0 }

// Error returns a combined error string of all issues, or empty string.
func (r *Result) Error() string {
	if r.OK() {
		return ""
	}
	var sb strings.Builder
	for i, iss := range r.Issues {
		if i > 0 {
			sb.WriteString("; ")
		}
		sb.WriteString(iss.Error())
	}
	return sb.String()
}

// Validate checks all keys and values in env for common problems.
// It returns a Result containing any issues found.
func Validate(env map[string]string) *Result {
	res := &Result{}
	for k, v := range env {
		if err := validateKey(k); err != nil {
			res.Issues = append(res.Issues, Issue{Key: k, Message: err.Error()})
		}
		if err := validateValue(k, v); err != nil {
			res.Issues = append(res.Issues, Issue{Key: k, Message: err.Error()})
		}
	}
	return res
}

// validateKey checks that a key is a valid POSIX environment variable name.
func validateKey(key string) error {
	if key == "" {
		return fmt.Errorf("empty key")
	}
	for i, ch := range key {
		switch {
		case ch == '_':
			// always allowed
		case i == 0 && unicode.IsDigit(ch):
			return fmt.Errorf("must not start with a digit")
		case !unicode.IsLetter(ch) && !unicode.IsDigit(ch):
			return fmt.Errorf("contains invalid character %q", ch)
		}
	}
	return nil
}

// validateValue checks for embedded newlines which can cause shell injection.
func validateValue(key, value string) error {
	if strings.ContainsAny(value, "\n\r") {
		return fmt.Errorf("value contains a newline character")
	}
	_ = key
	return nil
}
