package validator_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/validator"
)

func TestValidate_ValidEnv(t *testing.T) {
	env := map[string]string{
		"APP_NAME":    "myapp",
		"PORT":        "8080",
		"_PRIVATE":    "secret",
		"DB_HOST_123": "localhost",
	}
	res := validator.Validate(env)
	if !res.OK() {
		t.Fatalf("expected no issues, got: %s", res.Error())
	}
}

func TestValidate_EmptyKey(t *testing.T) {
	env := map[string]string{"": "value"}
	res := validator.Validate(env)
	if res.OK() {
		t.Fatal("expected issue for empty key")
	}
}

func TestValidate_KeyStartsWithDigit(t *testing.T) {
	env := map[string]string{"1BAD": "val"}
	res := validator.Validate(env)
	if res.OK() {
		t.Fatal("expected issue for key starting with digit")
	}
	if len(res.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(res.Issues))
	}
}

func TestValidate_KeyWithHyphen(t *testing.T) {
	env := map[string]string{"BAD-KEY": "val"}
	res := validator.Validate(env)
	if res.OK() {
		t.Fatal("expected issue for key with hyphen")
	}
}

func TestValidate_ValueWithNewline(t *testing.T) {
	env := map[string]string{"GOOD_KEY": "line1\nline2"}
	res := validator.Validate(env)
	if res.OK() {
		t.Fatal("expected issue for value with newline")
	}
}

func TestValidate_MultipleIssues(t *testing.T) {
	env := map[string]string{
		"1BAD":    "ok",
		"GOOD":    "has\nnewline",
		"BAD-KEY": "fine",
	}
	res := validator.Validate(env)
	if res.OK() {
		t.Fatal("expected multiple issues")
	}
	if len(res.Issues) < 3 {
		t.Fatalf("expected at least 3 issues, got %d: %s", len(res.Issues), res.Error())
	}
}

func TestIssue_Error(t *testing.T) {
	issue := validator.Issue{Key: "MY_KEY", Message: "some problem"}
	want := `key "MY_KEY": some problem`
	if got := issue.Error(); got != want {
		t.Fatalf("Issue.Error() = %q, want %q", got, want)
	}
}
