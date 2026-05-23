package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/envchain/internal/cli"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRun_DotenvOutput(t *testing.T) {
	p := writeTemp(t, "APP=hello\nPORT=9000\n")
	var buf bytes.Buffer
	code := cli.Run([]string{"--no-os", p}, &buf)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	out := buf.String()
	if !strings.Contains(out, "APP=hello") {
		t.Errorf("missing APP=hello in output: %q", out)
	}
}

func TestRun_JSONOutput(t *testing.T) {
	p := writeTemp(t, "KEY=value\n")
	var buf bytes.Buffer
	code := cli.Run([]string{"--no-os", "--format=json", p}, &buf)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.Contains(buf.String(), `"KEY"`) {
		t.Errorf("expected JSON output, got: %q", buf.String())
	}
}

func TestRun_InvalidFormat(t *testing.T) {
	var buf bytes.Buffer
	code := cli.Run([]string{"--format=xml"}, &buf)
	if code == 0 {
		t.Fatal("expected non-zero exit for invalid format")
	}
}

func TestRun_OutputToFile(t *testing.T) {
	p := writeTemp(t, "X=1\n")
	out := filepath.Join(t.TempDir(), "result.env")
	var buf bytes.Buffer
	code := cli.Run([]string{"--no-os", "-o", out, p}, &buf)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "X=1") {
		t.Errorf("output file missing X=1: %q", string(data))
	}
}

func TestRun_ValidateFlag_InvalidKey(t *testing.T) {
	p := writeTemp(t, "1BAD=value\n")
	var buf bytes.Buffer
	code := cli.Run([]string{"--no-os", "--validate", p}, &buf)
	if code == 0 {
		t.Fatal("expected non-zero exit when validation fails")
	}
}

func TestRun_ValidateFlag_ValidEnv(t *testing.T) {
	p := writeTemp(t, "GOOD_KEY=value\n")
	var buf bytes.Buffer
	code := cli.Run([]string{"--no-os", "--validate", p}, &buf)
	if code != 0 {
		t.Fatalf("expected exit 0 for valid env, got %d", code)
	}
}
