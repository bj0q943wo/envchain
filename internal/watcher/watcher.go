// Package watcher provides file system watching for .env files,
// notifying callers when any watched file changes.
package watcher

import (
	"os"
	"sync"
	"time"
)

// Event represents a file change notification.
type Event struct {
	Path string
}

// Watcher monitors a set of file paths for modifications.
type Watcher struct {
	mu       sync.Mutex
	paths    []string
	mtimes   map[string]time.Time
	Events   chan Event
	Errors   chan error
	done     chan struct{}
	interval time.Duration
}

// New creates a Watcher that polls at the given interval.
func New(interval time.Duration) *Watcher {
	return &Watcher{
		mtimes:   make(map[string]time.Time),
		Events:   make(chan Event, 16),
		Errors:   make(chan error, 4),
		done:     make(chan struct{}),
		interval: interval,
	}
}

// Add registers a file path to be watched.
func (w *Watcher) Add(path string) error {
	info, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	w.paths = append(w.paths, path)
	if err == nil {
		w.mtimes[path] = info.ModTime()
	}
	return nil
}

// Start begins polling in a background goroutine.
func (w *Watcher) Start() {
	go w.poll()
}

// Stop halts the watcher.
func (w *Watcher) Stop() {
	close(w.done)
}

func (w *Watcher) poll() {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-w.done:
			return
		case <-ticker.C:
			w.check()
		}
	}
}

func (w *Watcher) check() {
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, path := range w.paths {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}
		prev, seen := w.mtimes[path]
		if !seen || info.ModTime().After(prev) {
			w.mtimes[path] = info.ModTime()
			if seen {
				w.Events <- Event{Path: path}
			}
		}
	}
}
