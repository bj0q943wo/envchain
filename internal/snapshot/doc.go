// Package snapshot provides point-in-time captures of resolved environment
// variable maps and utilities for diffing successive snapshots.
//
// Typical usage:
//
//	prev := snapshot.Take(firstLoad)
//	// ... time passes, files change ...
//	next := snapshot.Take(secondLoad)
//	diff := prev.Compare(next)
//	if diff.HasChanges() {
//		fmt.Println("changed keys:", diff.Keys())
//	}
//
// Snapshots are immutable: mutating the original map after calling Take
// does not affect the snapshot. Env returns a fresh copy each time.
package snapshot
