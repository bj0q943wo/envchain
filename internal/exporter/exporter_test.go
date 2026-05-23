package exporter_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envchain/internal/exporter"
)

var sampleEnv = map[string]string{
	"APP_ENV":  "production",
	"DB_URL":   "postgres://localhost/mydb",
	"SECRET":   "p@ss w0rd!",
	"PLAIN":    "simple",
}

func TestWrite_DotenvFormat(t *testing.T) {
	var buf strings.Builder
	if err := exporter.Write(&buf, sampleEnv, exporter.FormatDotenv); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	assertContains(t, out, "APP_ENV=production")
	assertContains(t, out, "DB_URL=postgres://localhost/mydb")
	assertContains(t, out, "PLAIN=simple")
}

func TestWrite_ExportFormat(t *testing.T) {
	var buf strings.Builder
	if err := exporter.Write(&buf, sampleEnv, exporter.FormatExport); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	assertContains(t, out, "export APP_ENV=production")
	assertContains(t, out, "export PLAIN=simple")
	// value with spaces should be quoted
	assertContains(t, out, "export SECRET=")
}

func TestWrite_JSONFormat(t *testing.T) {
	var buf strings.Builder
	if err := exporter.Write(&buf, sampleEnv, exporter.FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	assertContains(t, out, `"APP_ENV": "production"`)
	assertContains(t, out, `"PLAIN": "simple"`)
	if !strings.HasPrefix(out, "{") || !strings.HasSuffix(strings.TrimSpace(out), "}") {
		t.Errorf("expected JSON object, got: %s", out)
	}
}

func TestWrite_SortedOutput(t *testing.T) {
	env := map[string]string{"Z": "last", "A": "first", "M": "middle"}
	var buf strings.Builder
	_ = exporter.Write(&buf, env, exporter.FormatDotenv)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 || !strings.HasPrefix(lines[0], "A=") || !strings.HasPrefix(lines[2], "Z=") {
		t.Errorf("expected sorted output, got: %v", lines)
	}
}

func TestWrite_UnknownFormat(t *testing.T) {
	var buf strings.Builder
	err := exporter.Write(&buf, sampleEnv, exporter.Format("xml"))
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func assertContains(t *testing.T, haystack, needle string) {
	t.Helper()
	if !strings.Contains(haystack, needle) {
		t.Errorf("expected output to contain %q\ngot:\n%s", needle, haystack)
	}
}
