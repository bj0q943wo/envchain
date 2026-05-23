package exporter

import "strings"

// ParseFormat converts a string to a Format constant, case-insensitively.
// Returns an error string (second value) if the format is unrecognised.
func ParseFormat(s string) (Format, bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "export", "sh", "shell":
		return FormatExport, true
	case "dotenv", "env", ".env":
		return FormatDotenv, true
	case "json":
		return FormatJSON, true
	default:
		return "", false
	}
}

// String returns the canonical string representation of a Format.
func (f Format) String() string {
	return string(f)
}

// ValidFormats returns all supported format identifiers for use in help text.
func ValidFormats() []string {
	return []string{
		string(FormatExport),
		string(FormatDotenv),
		string(FormatJSON),
	}
}
