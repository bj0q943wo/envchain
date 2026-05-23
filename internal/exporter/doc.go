// Package exporter serializes a resolved environment variable map into
// one of several output formats suitable for shell consumption or
// interoperability with other tooling.
//
// Supported formats:
//
//	 FormatExport  — `export KEY=VALUE` lines, suitable for `eval $(envchain ...)`
//	 FormatDotenv  — plain `KEY=VALUE` lines, compatible with .env file consumers
//	 FormatJSON    — a JSON object mapping keys to string values
//
// Example usage:
//
//	 env := map[string]string{"FOO": "bar", "BAZ": "qux"}
//	 err := exporter.Write(os.Stdout, env, exporter.FormatExport)
//
package exporter
