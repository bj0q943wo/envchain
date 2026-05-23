// Package cli wires together argument parsing, loading, resolving, validating,
// and exporting for the envchain command-line tool.
package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/yourorg/envchain/internal/exporter"
	"github.com/yourorg/envchain/internal/resolver"
	"github.com/yourorg/envchain/internal/validator"
)

// args holds parsed CLI configuration.
type args struct {
	files    []resolver.FileEntry
	format   exporter.Format
	output   string
	noOS     bool
	validate bool
}

// Run is the main entry point called by cmd/envchain/main.go.
func Run(argv []string, stdout io.Writer) int {
	a, err := parseArgs(argv)
	if err != nil {
		fmt_err(err)
		return 1
	}
	return execute(a, stdout)
}

func execute(a *args, stdout io.Writer) int {
	chain := resolver.NewChain(a.files...)

	var env map[string]string
	var err error
	if a.noOS {
		env, err = chain.Resolve()
	} else {
		env, err = chain.ResolveWithOS()
	}
	if err != nil {
		fmt_err(err)
		return 1
	}

	if a.validate {
		res := validator.Validate(env)
		if !res.OK() {
			fmt_err(fmt.Errorf("validation failed: %s", res.Error()))
			return 1
		}
	}

	out := stdout
	if a.output != "" {
		f, ferr := os.Create(a.output)
		if ferr != nil {
			fmt_err(ferr)
			return 1
		}
		defer f.Close()
		out = f
	}

	if werr := exporter.Write(out, env, a.format); werr != nil {
		fmt_err(werr)
		return 1
	}
	return 0
}

func parseArgs(argv []string) (*args, error) {
	fs := flag.NewFlagSet("envchain", flag.ContinueOnError)

	formatStr := fs.String("format", "dotenv", "output format: dotenv, export, json")
	output := fs.String("o", "", "write output to file instead of stdout")
	noOS := fs.Bool("no-os", false, "do not merge OS environment variables")
	validateFlag := fs.Bool("validate", false, "validate key names and values before output")

	if err := fs.Parse(argv); err != nil {
		return nil, err
	}

	fmt, err := exporter.ParseFormat(*formatStr)
	if err != nil {
		return nil, err
	}

	var entries []resolver.FileEntry
	for _, f := range fs.Args() {
		optional := len(f) > 0 && f[0] == '?'
		path := f
		if optional {
			path = f[1:]
		}
		entries = append(entries, resolver.FileEntry{Path: path, Optional: optional})
	}

	return &args{
		files:    entries,
		format:   fmt,
		output:   *output,
		noOS:     *noOS,
		validate: *validateFlag,
	}, nil
}

func fmt_err(err error) {
	fmt.Fprintf(os.Stderr, "envchain: %s\n", err)
}
