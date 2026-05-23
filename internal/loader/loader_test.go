package loader

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return path
}

func TestLoadFile_Basic(t *testing.T) {
	path := writeTemp(t, "KEY=value\nFOO=bar\n")
	env, err := LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", env["KEY"])
	}
	if env["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", env["FOO"])
	}
}

func TestLoadFile_SkipsCommentsAndBlanks(t *testing.T) {
	path := writeTemp(t, "# comment\n\nKEY=hello\n")
	env, err := LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env) != 1 {
		t.Errorf("expected 1 entry, got %d", len(env))
	}
}

func TestLoadFile_QuotedValues(t *testing.T) {
	path := writeTemp(t, `A="double quoted"` + "\n" + `B='single quoted'` + "\n")
	env, err := LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["A"] != "double quoted" {
		t.Errorf("expected 'double quoted', got %q", env["A"])
	}
	if env["B"] != "single quoted" {
		t.Errorf("expected 'single quoted', got %q", env["B"])
	}
}

func TestLoadFile_InvalidLine(t *testing.T) {
	path := writeTemp(t, "BADLINE\n")
	_, err := LoadFile(path)
	if err == nil {
		t.Fatal("expected error for invalid line, got nil")
	}
}

func TestLoadFile_NotFound(t *testing.T) {
	_, err := LoadFile("/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestMerge_Precedence(t *testing.T) {
	base := Env{"KEY": "base", "SHARED": "from-base"}
	override := Env{"SHARED": "from-override", "EXTRA": "only-here"}

	result := Merge(base, override)

	if result["KEY"] != "base" {
		t.Errorf("expected KEY=base, got %q", result["KEY"])
	}
	if result["SHARED"] != "from-override" {
		t.Errorf("expected SHARED=from-override, got %q", result["SHARED"])
	}
	if result["EXTRA"] != "only-here" {
		t.Errorf("expected EXTRA=only-here, got %q", result["EXTRA"])
	}
}
