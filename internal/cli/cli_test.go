package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestRun_DotenvOutput(t *testing.T) {
	file := writeTemp(t, "FOO=bar\nBAZ=qux\n")
	var buf bytes.Buffer
	cfg := &Config{Files: []string{file}, Format: "dotenv", NoOS: true}
	if err := execute(cfg, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "FOO=bar") {
		t.Errorf("expected FOO=bar in output, got: %s", out)
	}
}

func TestRun_JSONOutput(t *testing.T) {
	file := writeTemp(t, "KEY=value\n")
	var buf bytes.Buffer
	cfg := &Config{Files: []string{file}, Format: "json", NoOS: true}
	if err := execute(cfg, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), `"KEY"`) {
		t.Errorf("expected JSON key in output, got: %s", buf.String())
	}
}

func TestRun_InvalidFormat(t *testing.T) {
	file := writeTemp(t, "A=1\n")
	cfg := &Config{Files: []string{file}, Format: "xml", NoOS: true}
	err := execute(cfg, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for invalid format")
	}
	if !strings.Contains(err.Error(), "xml") {
		t.Errorf("error should mention format name, got: %v", err)
	}
}

func TestRun_OutputToFile(t *testing.T) {
	file := writeTemp(t, "OUT=written\n")
	out := filepath.Join(t.TempDir(), "result.env")
	cfg := &Config{Files: []string{file}, Format: "dotenv", Output: out, NoOS: true}
	if err := execute(cfg, &bytes.Buffer{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "OUT=written") {
		t.Errorf("expected OUT=written in file, got: %s", string(data))
	}
}

func TestParseArgs_Defaults(t *testing.T) {
	cfg, err := parseArgs([]string{"a.env", "b.env"})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Format != "dotenv" {
		t.Errorf("expected default format dotenv, got %s", cfg.Format)
	}
	if len(cfg.Files) != 2 {
		t.Errorf("expected 2 files, got %d", len(cfg.Files))
	}
	if cfg.NoOS {
		t.Error("expected NoOS to be false by default")
	}
}
