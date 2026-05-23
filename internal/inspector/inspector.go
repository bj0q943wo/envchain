// Package inspector provides utilities for inspecting and summarizing
// the resolved environment chain, including source tracking for each key.
package inspector

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Entry represents a single resolved environment variable with metadata
// about which source file it originated from.
type Entry struct {
	Key    string
	Value  string
	Source string // file path or "os" for OS environment
}

// Report holds the full inspection result for a resolved environment chain.
type Report struct {
	Entries []Entry
}

// Build constructs a Report from a slice of (source, env) pairs in
// precedence order (later entries override earlier ones).
func Build(sources []SourceEnv) *Report {
	seen := make(map[string]Entry)
	order := []string{}

	for _, s := range sources {
		for k, v := range s.Env {
			if _, exists := seen[k]; !exists {
				order = append(order, k)
			}
			seen[k] = Entry{Key: k, Value: v, Source: s.Source}
		}
	}

	sort.Strings(order)

	entries := make([]Entry, 0, len(order))
	for _, k := range order {
		entries = append(entries, seen[k])
	}

	return &Report{Entries: entries}
}

// SourceEnv pairs a source label with its key-value environment map.
type SourceEnv struct {
	Source string
	Env    map[string]string
}

// Print writes a human-readable summary of the report to w.
// Each line shows: KEY=VALUE  (source: <origin>)
func Print(r *Report, w io.Writer) error {
	for _, e := range r.Entries {
		line := fmt.Sprintf("%-30s (source: %s)\n",
			fmt.Sprintf("%s=%s", e.Key, truncate(e.Value, 40)),
			e.Source,
		)
		if _, err := io.WriteString(w, line); err != nil {
			return err
		}
	}
	return nil
}

func truncate(s string, max int) string {
	s = strings.ReplaceAll(s, "\n", "\\n")
	if len(s) > max {
		return s[:max] + "..."
	}
	return s
}
