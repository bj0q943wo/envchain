package resolver_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envchain/envchain/internal/resolver"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writeTemp: %v", err)
	}
	return path
}

func TestChain_Resolve_MergesInOrder(t *testing.T) {
	base := writeTemp(t, "FOO=base\nBAR=base\n")
	override := writeTemp(t, "FOO=override\nBAZ=new\n")

	c := resolver.NewChain(base, override)
	envs, err := c.Resolve(true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if envs["FOO"] != "override" {
		t.Errorf("FOO: got %q, want %q", envs["FOO"], "override")
	}
	if envs["BAR"] != "base" {
		t.Errorf("BAR: got %q, want %q", envs["BAR"], "base")
	}
	if envs["BAZ"] != "new" {
		t.Errorf("BAZ: got %q, want %q", envs["BAZ"], "new")
	}
}

func TestChain_Resolve_SkipsMissingOptional(t *testing.T) {
	existing := writeTemp(t, "KEY=value\n")
	c := resolver.NewChain(existing, "/nonexistent/.env.local")

	envs, err := c.Resolve(false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if envs["KEY"] != "value" {
		t.Errorf("KEY: got %q, want %q", envs["KEY"], "value")
	}
}

func TestChain_Resolve_ErrorOnMissingRequired(t *testing.T) {
	c := resolver.NewChain("/nonexistent/.env")
	_, err := c.Resolve(true)
	if err == nil {
		t.Fatal("expected error for missing required file, got nil")
	}
}

func TestChain_ResolveWithOS_ChainOverridesBase(t *testing.T) {
	override := writeTemp(t, "APP_ENV=production\n")
	base := map[string]string{
		"APP_ENV": "development",
		"HOME":    "/home/user",
	}

	c := resolver.NewChain(override)
	result, err := c.ResolveWithOS(base, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["APP_ENV"] != "production" {
		t.Errorf("APP_ENV: got %q, want %q", result["APP_ENV"], "production")
	}
	if result["HOME"] != "/home/user" {
		t.Errorf("HOME should be preserved from base")
	}
}
