package watcher

import "time"

// MultiWatcher wraps Watcher and fans out events to multiple subscribers.
type MultiWatcher struct {
	w           *Watcher
	mu          chanMu
	subscribers []chan Event
}

type chanMu struct{}

// NewMulti creates a MultiWatcher polling at the given interval.
func NewMulti(interval time.Duration) *MultiWatcher {
	mw := &MultiWatcher{
		w: New(interval),
	}
	return mw
}

// Add registers a path with the underlying Watcher.
func (m *MultiWatcher) Add(path string) error {
	return m.w.Add(path)
}

// Subscribe returns a channel that receives events for this subscriber.
func (m *MultiWatcher) Subscribe() <-chan Event {
	ch := make(chan Event, 8)
	m.subscribers = append(m.subscribers, ch)
	return ch
}

// Start begins polling and fanning out to all subscribers.
func (m *MultiWatcher) Start() {
	m.w.Start()
	go m.fanOut()
}

// Stop halts the underlying watcher and closes subscriber channels.
func (m *MultiWatcher) Stop() {
	m.w.Stop()
	for _, ch := range m.subscribers {
		close(ch)
	}
}

func (m *MultiWatcher) fanOut() {
	for ev := range m.w.Events {
		for _, ch := range m.subscribers {
			select {
			case ch <- ev:
			default:
			}
		}
	}
}
