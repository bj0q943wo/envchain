// Package validator provides lightweight validation for environment variable
// maps produced by the loader and resolver packages.
//
// Validation checks include:
//
//   - Key naming rules: keys must follow POSIX conventions (letters, digits,
//     underscores; must not start with a digit).
//   - Value safety: values must not contain embedded newline characters, which
//     can cause unexpected behaviour when the environment is exported to a shell.
//
// Usage:
//
//	res := validator.Validate(env)
//	if !res.OK() {
//		log.Fatalf("env validation failed: %s", res.Error())
//	}
package validator
