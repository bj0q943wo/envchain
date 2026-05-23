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
// This package depends on internal/loader for file parsing and merging.
package resolver
