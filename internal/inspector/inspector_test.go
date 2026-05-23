package inspector_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/envchain/internal/inspector"
)

func TestBuild_SingleSource(t *testing.T) {
	sources := []inspector.SourceEnv{
		{Source: ".env", Env: map[string]string{"FOO": "bar", "BAZ": "qux"}},
	}
	r := inspector.Build(sources)
	if len(r.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(r.Entries))
	}
}

func TestBuild_LaterSourceOverrides(t *testing.T) {
	sources := []inspector.SourceEnv{
		{Source: ".env", Env: map[string]string{"FOO": "original"}},
		{Source: ".env.local", Env: map[string]string{"FOO": "overridden"}},
	}
	r := inspector.Build(sources)
	if len(r.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(r.Entries))
	}
	e := r.Entries[0]
	if e.Value != "overridden" {
		t.Errorf("expected overridden, got %q", e.Value)
	}
	if e.Source != ".env.local" {
		t.Errorf("expected source .env.local, got %q", e.Source)
	}
}

func TestBuild_SortedKeys(t *testing.T) {
	sources := []inspector.SourceEnv{
		{Source: ".env", Env: map[string]string{"ZEBRA": "1", "ALPHA": "2", "MANGO": "3"}},
	}
	r := inspector.Build(sources)
	keys := make([]string, len(r.Entries))
	for i, e := range r.Entries {
		keys[i] = e.Key
	}
	if keys[0] != "ALPHA" || keys[1] != "MANGO" || keys[2] != "ZEBRA" {
		t.Errorf("keys not sorted: %v", keys)
	}
}

func TestBuild_EmptySources(t *testing.T) {
	r := inspector.Build(nil)
	if len(r.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(r.Entries))
	}
}

func TestPrint_Output(t *testing.T) {
	sources := []inspector.SourceEnv{
		{Source: ".env", Env: map[string]string{"API_KEY": "secret"}},
	}
	r := inspector.Build(sources)
	var buf bytes.Buffer
	if err := inspector.Print(r, &buf); err != nil {
		t.Fatalf("Print error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "API_KEY=secret") {
		t.Errorf("expected API_KEY=secret in output, got: %s", out)
	}
	if !strings.Contains(out, "source: .env") {
		t.Errorf("expected source annotation in output, got: %s", out)
	}
}

func TestPrint_TruncatesLongValues(t *testing.T) {
	long := strings.Repeat("x", 60)
	sources := []inspector.SourceEnv{
		{Source: "os", Env: map[string]string{"LONG": long}},
	}
	r := inspector.Build(sources)
	var buf bytes.Buffer
	_ = inspector.Print(r, &buf)
	if strings.Contains(buf.String(), long) {
		t.Error("expected long value to be truncated in output")
	}
}
