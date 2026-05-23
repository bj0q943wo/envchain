package watcher_test

import (
	"os"
	"testing"
	"time"

	"github.com/yourorg/envchain/internal/watcher"
)

func TestMultiWatcher_FansOutToSubscribers(t *testing.T) {
	path := writeTempEnv(t, "A=1\n")

	mw := watcher.NewMulti(20 * time.Millisecond)
	if err := mw.Add(path); err != nil {
		t.Fatalf("Add: %v", err)
	}

	ch1 := mw.Subscribe()
	ch2 := mw.Subscribe()

	mw.Start()
	defer mw.Stop()

	time.Sleep(30 * time.Millisecond)
	if err := os.WriteFile(path, []byte("A=2\n"), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	timeout := time.After(300 * time.Millisecond)
	var got1, got2 bool
	for !got1 || !got2 {
		select {
		case ev, ok := <-ch1:
			if ok && ev.Path == path {
				got1 = true
			}
		case ev, ok := <-ch2:
			if ok && ev.Path == path {
				got2 = true
			}
		case <-timeout:
			t.Fatalf("timed out: got1=%v got2=%v", got1, got2)
		}
	}
}

func TestMultiWatcher_NoSubscribers_NoDeadlock(t *testing.T) {
	path := writeTempEnv(t, "B=1\n")
	mw := watcher.NewMulti(20 * time.Millisecond)
	if err := mw.Add(path); err != nil {
		t.Fatalf("Add: %v", err)
	}
	mw.Start()
	defer mw.Stop()

	time.Sleep(30 * time.Millisecond)
	_ = os.WriteFile(path, []byte("B=2\n"), 0o600)
	time.Sleep(60 * time.Millisecond)
	// no hang — test passes if we reach here
}
