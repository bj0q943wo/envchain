package cli

import (
	"strings"
	"testing"
)

func TestParseArgs_OptionalFiles(t *testing.T) {
	cfg, err := parseArgs([]string{"--optional", "a.env, b.env", "main.env"})
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Optional) != 2 {
		t.Fatalf("expected 2 optional files, got %d: %v", len(cfg.Optional), cfg.Optional)
	}
	if cfg.Optional[0] != "a.env" {
		t.Errorf("expected a.env, got %s", cfg.Optional[0])
	}
}

func TestParseArgs_NoOS(t *testing.T) {
	cfg, err := parseArgs([]string{"--no-os", "main.env"})
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.NoOS {
		t.Error("expected NoOS to be true")
	}
}

func TestParseArgs_OutputFlag(t *testing.T) {
	cfg, err := parseArgs([]string{"--output", "out.env", "main.env"})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Output != "out.env" {
		t.Errorf("expected out.env, got %s", cfg.Output)
	}
}

func TestParseArgs_FormatFlag(t *testing.T) {
	for _, fmt := range []string{"dotenv", "export", "json"} {
		cfg, err := parseArgs([]string{"--format", fmt})
		if err != nil {
			t.Fatalf("unexpected error for format %s: %v", fmt, err)
		}
		if cfg.Format != fmt {
			t.Errorf("expected format %s, got %s", fmt, cfg.Format)
		}
	}
}

func TestParseArgs_UnknownFlag(t *testing.T) {
	_, err := parseArgs([]string{"--bogus-flag"})
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}
	if !strings.Contains(err.Error(), "bogus-flag") {
		t.Errorf("error should mention flag name, got: %v", err)
	}
}
