package resolver

import (
	"fmt"
	"os"

	"github.com/envchain/envchain/internal/loader"
)

// Chain represents an ordered list of .env file paths to be resolved.
type Chain struct {
	Files []string
}

// NewChain creates a Chain from the given ordered list of file paths.
// Files later in the list take higher precedence.
func NewChain(files ...string) *Chain {
	return &Chain{Files: files}
}

// Resolve loads and merges all .env files in the chain, with later files
// overriding earlier ones. Missing files are skipped unless required.
func (c *Chain) Resolve(required bool) (map[string]string, error) {
	merged := make(map[string]string)

	for _, path := range c.Files {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if required {
				return nil, fmt.Errorf("required env file not found: %s", path)
			}
			continue
		}

		envs, err := loader.LoadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to load %s: %w", path, err)
		}

		loader.Merge(merged, envs)
	}

	return merged, nil
}

// ResolveWithOS merges the chain result on top of a copy of the provided
// base map (typically os.Environ parsed), so chain values override the base.
func (c *Chain) ResolveWithOS(base map[string]string, required bool) (map[string]string, error) {
	result := make(map[string]string, len(base))
	for k, v := range base {
		result[k] = v
	}

	chainEnvs, err := c.Resolve(required)
	if err != nil {
		return nil, err
	}

	loader.Merge(result, chainEnvs)
	return result, nil
}
