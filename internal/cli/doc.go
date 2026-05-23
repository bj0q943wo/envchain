// Package cli provides the command-line interface for envchain.
//
// It parses flags and positional arguments, delegates environment
// resolution to the resolver package, and writes the merged result
// via the exporter package.
//
// Usage:
//
//	envchain [flags] [file...]
//
// Flags:
//
//	-format   string   output format: dotenv, export, json (default "dotenv")
//	-optional string   comma-separated list of optional env files
//	-output   string   write output to file instead of stdout
//	-no-os            do not merge OS environment as base
package cli
