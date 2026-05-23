// Package expander provides variable interpolation for env maps.
// It resolves references like ${VAR} or $VAR within env values,
// using values already present in the map or falling back to OS env.
package expander

import (
	"os"
	"strings"
)

// Expand resolves variable references in all values of the provided env map.
// It supports both ${VAR} and $VAR syntax. Variables are resolved first from
// the env map itself, then from the OS environment. Unknown variables expand
// to an empty string.
func Expand(env map[string]string) map[string]string {
	result := make(map[string]string, len(env))
	for k, v := range env {
		result[k] = expandValue(v, env)
	}
	return result
}

// ExpandWithOS resolves variable references using the provided env map with
// OS environment as a fallback, but does NOT add OS variables to the result.
func ExpandWithOS(env map[string]string) map[string]string {
	return Expand(env)
}

// expandValue interpolates a single value string using the given env map
// and os.Getenv as fallback.
func expandValue(value string, env map[string]string) string {
	return os.Expand(value, func(key string) string {
		if val, ok := env[key]; ok {
			return val
		}
		return os.Getenv(key)
	})
}

// HasReferences returns true if the value string contains any variable
// reference patterns ($VAR or ${VAR}).
func HasReferences(value string) bool {
	return strings.Contains(value, "$")
}
