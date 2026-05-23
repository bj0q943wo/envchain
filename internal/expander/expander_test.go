package expander_test

import (
	"os"
	"testing"

	"github.com/user/envchain/internal/expander"
)

func TestExpand_SimpleSubstitution(t *testing.T) {
	env := map[string]string{
		"HOME": "/home/user",
		"CONFIG": "${HOME}/.config",
	}
	result := expander.Expand(env)
	if result["CONFIG"] != "/home/user/.config" {
		t.Errorf("expected /home/user/.config, got %s", result["CONFIG"])
	}
}

func TestExpand_DollarSyntax(t *testing.T) {
	env := map[string]string{
		"NAME": "world",
		"GREETING": "hello $NAME",
	}
	result := expander.Expand(env)
	if result["GREETING"] != "hello world" {
		t.Errorf("expected 'hello world', got %s", result["GREETING"])
	}
}

func TestExpand_UnknownVarExpandsEmpty(t *testing.T) {
	env := map[string]string{
		"VAL": "prefix_${UNDEFINED_XYZ}_suffix",
	}
	os.Unsetenv("UNDEFINED_XYZ")
	result := expander.Expand(env)
	if result["VAL"] != "prefix__suffix" {
		t.Errorf("expected 'prefix__suffix', got %s", result["VAL"])
	}
}

func TestExpand_FallsBackToOS(t *testing.T) {
	os.Setenv("OS_ONLY_VAR", "from-os")
	defer os.Unsetenv("OS_ONLY_VAR")

	env := map[string]string{
		"DERIVED": "${OS_ONLY_VAR}/path",
	}
	result := expander.Expand(env)
	if result["DERIVED"] != "from-os/path" {
		t.Errorf("expected 'from-os/path', got %s", result["DERIVED"])
	}
}

func TestExpand_NoReferencesUnchanged(t *testing.T) {
	env := map[string]string{
		"PLAIN": "just-a-value",
	}
	result := expander.Expand(env)
	if result["PLAIN"] != "just-a-value" {
		t.Errorf("expected 'just-a-value', got %s", result["PLAIN"])
	}
}

func TestHasReferences(t *testing.T) {
	if !expander.HasReferences("${FOO}") {
		t.Error("expected true for ${FOO}")
	}
	if !expander.HasReferences("$BAR") {
		t.Error("expected true for $BAR")
	}
	if expander.HasReferences("no-refs-here") {
		t.Error("expected false for plain string")
	}
}
