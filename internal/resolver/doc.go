// Package resolver provides the Chain type for loading and merging multiple
// .env files with explicit override precedence.
//
// A Chain is an ordered sequence of .env file paths. Files appearing later
// in the chain take higher precedence, so values defined in a later file
// will override the same key defined in an earlier file.
//
// Basic usage:
//
//	c := resolver.NewChain(".env", ".env.local", ".env.production")
//	envs, err := c.Resolve(false) // false = skip missing files
//
// To layer chain values on top of an existing map (e.g. parsed OS environment):
//
//	result, err := c.ResolveWithOS(osEnvMap, false)
//
// Precedence order (lowest to highest):
//
//  1. OS environment variables (when using ResolveWithOS)
//  2. First file in the chain (e.g. ".env")
//  3. Subsequent files in the chain (e.g. ".env.local", ".env.production")
//
// Missing files are silently skipped when the strict argument is false.
// When strict is true, a missing file causes Resolve to return an error.
//
// This package depends on internal/loader for file parsing and merging.
package resolver
