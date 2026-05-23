// Package exporter provides functionality to export resolved env vars
// in various shell-compatible formats.
package exporter

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Format represents the output format for exported variables.
type Format string

const (
	// FormatExport emits `export KEY=VALUE` lines (bash/zsh compatible).
	FormatExport Format = "export"
	// FormatDotenv emits `KEY=VALUE` lines (plain .env style).
	FormatDotenv Format = "dotenv"
	// FormatJSON emits a JSON object.
	FormatJSON Format = "json"
)

// Write serializes the given env map to w in the requested format.
// Keys are sorted alphabetically for deterministic output.
func Write(w io.Writer, env map[string]string, format Format) error {
	keys := sortedKeys(env)

	switch format {
	case FormatExport:
		return writeExport(w, env, keys)
	case FormatDotenv:
		return writeDotenv(w, env, keys)
	case FormatJSON:
		return writeJSON(w, env, keys)
	default:
		return fmt.Errorf("exporter: unknown format %q", format)
	}
}

func writeExport(w io.Writer, env map[string]string, keys []string) error {
	for _, k := range keys {
		if _, err := fmt.Fprintf(w, "export %s=%s\n", k, shellQuote(env[k])); err != nil {
			return err
		}
	}
	return nil
}

func writeDotenv(w io.Writer, env map[string]string, keys []string) error {
	for _, k := range keys {
		if _, err := fmt.Fprintf(w, "%s=%s\n", k, env[k]); err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(w io.Writer, env map[string]string, keys []string) error {
	var sb strings.Builder
	sb.WriteString("{\n")
	for i, k := range keys {
		comma := ","
		if i == len(keys)-1 {
			comma = ""
		}
		sb.WriteString(fmt.Sprintf("  %q: %q%s\n", k, env[k], comma))
	}
	sb.WriteString("}\n")
	_, err := fmt.Fprint(w, sb.String())
	return err
}

// shellQuote wraps value in single quotes if it contains special characters.
func shellQuote(v string) string {
	if strings.ContainsAny(v, " \t\n\r'\"\\$`!#&;|<>(){}") {
		return "'" + strings.ReplaceAll(v, "'", "'\\''")+"'"
	}
	return v
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
