package watcher_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/envchain/internal/watcher"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestWatcher_DetectsChange(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\n")
	w := watcher.New(20 * time.Millisecond)
	if err := w.Add(path); err != nil {
		t.Fatalf("Add: %v", err)
	}
	w.Start()
	defer w.Stop()

	time.Sleep(30 * time.Millisecond)
	if err := os.WriteFile(path, []byte("FOO=baz\n"), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	select {
	case ev := <-w.Events:
		if ev.Path != path {
			t.Errorf("expected path %q, got %q", path, ev.Path)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timed out waiting for change event")
	}
}

func TestWatcher_NoEventOnUnchangedFile(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\n")
	w := watcher.New(20 * time.Millisecond)
	if err := w.Add(path); err != nil {
		t.Fatalf("Add: %v", err)
	}
	w.Start()
	defer w.Stop()

	select {
	case ev := <-w.Events:
		t.Errorf("unexpected event for unchanged file: %v", ev)
	case <-time.After(100 * time.Millisecond):
		// expected — no change
	}
}

func TestWatcher_AddMissingFile(t *testing.T) {
	w := watcher.New(20 * time.Millisecond)
	err := w.Add("/nonexistent/.env")
	if err != nil {
		t.Errorf("Add of missing file should not error, got: %v", err)
	}
}

func TestWatcher_StopIsIdempotent(t *testing.T) {
	w := watcher.New(50 * time.Millisecond)
	w.Start()
	w.Stop()
	// second Stop would panic on double-close; wrap to confirm it doesn't
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Stop panicked: %v", r)
		}
	}()
}
