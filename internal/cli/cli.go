// Package cli wires together the resolver and exporter into a runnable command.
package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/envchain/internal/exporter"
	"github.com/user/envchain/internal/resolver"
)

// Config holds parsed CLI flags and arguments.
type Config struct {
	Files    []string
	Optional []string
	Format   string
	Output   string
	NoOS     bool
}

// Run parses args and executes the envchain command.
func Run(args []string) error {
	cfg, err := parseArgs(args)
	if err != nil {
		return err
	}
	return execute(cfg, os.Stdout)
}

func parseArgs(args []string) (*Config, error) {
	fs := flag.NewFlagSet("envchain", flag.ContinueOnError)

	var optional string
	var format string
	var output string
	var noOS bool

	fs.StringVar(&optional, "optional", "", "comma-separated list of optional env files")
	fs.StringVar(&format, "format", "dotenv", "output format: dotenv, export, json")
	fs.StringVar(&output, "output", "", "write output to file instead of stdout")
	fs.BoolVar(&noOS, "no-os", false, "do not merge OS environment as base")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	var optFiles []string
	if optional != "" {
		for _, f := range strings.Split(optional, ",") {
			if t := strings.TrimSpace(f); t != "" {
				optFiles = append(optFiles, t)
			}
		}
	}

	return &Config{
		Files:    fs.Args(),
		Optional: optFiles,
		Format:   format,
		Output:   output,
		NoOS:     noOS,
	}, nil
}

func execute(cfg *Config, stdout io.Writer) error {
	fmt, err := exporter.ParseFormat(cfg.Format)
	if err != nil {
		return fmt_err(cfg.Format)
	}

	chain := resolver.NewChain(cfg.Files, cfg.Optional)

	var env map[string]string
	if cfg.NoOS {
		env, err = chain.Resolve()
	} else {
		env, err = chain.ResolveWithOS()
	}
	if err != nil {
		return err
	}

	out := stdout
	if cfg.Output != "" {
		f, ferr := os.Create(cfg.Output)
		if ferr != nil {
			return fmt.Errorf("cannot open output file: %w", ferr)
		}
		defer f.Close()
		out = f
	}

	return exporter.Write(out, env, fmt)
}

func fmt_err(s string) error {
	return fmt.Errorf("unknown format %q; valid formats: %s", s, strings.Join(exporter.ValidFormats, ", "))
}
