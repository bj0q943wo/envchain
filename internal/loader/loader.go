package loader

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Env holds a map of key-value pairs parsed from a .env file.
type Env map[string]string

// LoadFile reads a .env file and returns an Env map.
// Lines starting with '#' are treated as comments and skipped.
// Empty lines are also skipped.
// Values may optionally be quoted with single or double quotes.
func LoadFile(path string) (Env, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("loader: open %q: %w", path, err)
	}
	defer f.Close()

	env := make(Env)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("loader: %q line %d: %w", path, lineNum, err)
		}

		env[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("loader: scan %q: %w", path, err)
	}

	return env, nil
}

// Merge combines multiple Env maps in order, with later maps taking precedence.
func Merge(envs ...Env) Env {
	result := make(Env)
	for _, e := range envs {
		for k, v := range e {
			result[k] = v
		}
	}
	return result
}

func parseLine(line string) (string, string, error) {
	idx := strings.IndexByte(line, '=')
	if idx < 0 {
		return "", "", fmt.Errorf("invalid format %q: missing '='" , line)
	}

	key := strings.TrimSpace(line[:idx])
	if key == "" {
		return "", "", fmt.Errorf("empty key in line %q", line)
	}

	value := strings.TrimSpace(line[idx+1:])
	value = stripQuotes(value)

	return key, value, nil
}

func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
