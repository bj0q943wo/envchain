// Package snapshot provides functionality to capture and compare
// environment variable maps, enabling change detection between
// successive loads of .env file chains.
package snapshot

import (
	"maps"
	"slices"
)

// Snapshot holds an immutable copy of an environment map at a point in time.
type Snapshot struct {
	env map[string]string
}

// Take creates a new Snapshot from the provided environment map.
// The map is copied so later mutations do not affect the snapshot.
func Take(env map[string]string) *Snapshot {
	copy := make(map[string]string, len(env))
	maps.Copy(copy, env)
	return &Snapshot{env: copy}
}

// Diff describes the changes between two snapshots.
type Diff struct {
	Added   map[string]string // keys present in next but not prev
	Removed map[string]string // keys present in prev but not next
	Changed map[string]string // keys in both but with different values (new value stored)
}

// HasChanges reports whether the diff contains any additions, removals, or changes.
func (d *Diff) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0 || len(d.Changed) > 0
}

// Keys returns a sorted slice of all affected keys across all diff categories.
func (d *Diff) Keys() []string {
	seen := make(map[string]struct{})
	for k := range d.Added {
		seen[k] = struct{}{}
	}
	for k := range d.Removed {
		seen[k] = struct{}{}
	}
	for k := range d.Changed {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

// Compare returns a Diff between the receiver (previous) and next snapshot.
func (s *Snapshot) Compare(next *Snapshot) *Diff {
	d := &Diff{
		Added:   make(map[string]string),
		Removed: make(map[string]string),
		Changed: make(map[string]string),
	}
	for k, v := range next.env {
		if prev, ok := s.env[k]; !ok {
			d.Added[k] = v
		} else if prev != v {
			d.Changed[k] = v
		}
	}
	for k, v := range s.env {
		if _, ok := next.env[k]; !ok {
			d.Removed[k] = v
		}
	}
	return d
}

// Env returns a copy of the environment map held by this snapshot.
func (s *Snapshot) Env() map[string]string {
	copy := make(map[string]string, len(s.env))
	maps.Copy(copy, s.env)
	return copy
}
