package snapshot_test

import (
	"testing"

	"github.com/user/envchain/internal/snapshot"
)

func TestTake_CopiesMap(t *testing.T) {
	original := map[string]string{"A": "1", "B": "2"}
	snap := snapshot.Take(original)
	original["A"] = "mutated"
	env := snap.Env()
	if env["A"] != "1" {
		t.Errorf("expected snapshot to be immutable, got %q", env["A"])
	}
}

func TestCompare_Added(t *testing.T) {
	prev := snapshot.Take(map[string]string{"A": "1"})
	next := snapshot.Take(map[string]string{"A": "1", "B": "2"})
	diff := prev.Compare(next)
	if !diff.HasChanges() {
		t.Fatal("expected changes")
	}
	if diff.Added["B"] != "2" {
		t.Errorf("expected B=2 in Added, got %v", diff.Added)
	}
	if len(diff.Removed) != 0 || len(diff.Changed) != 0 {
		t.Errorf("unexpected removals or changes")
	}
}

func TestCompare_Removed(t *testing.T) {
	prev := snapshot.Take(map[string]string{"A": "1", "B": "2"})
	next := snapshot.Take(map[string]string{"A": "1"})
	diff := prev.Compare(next)
	if diff.Removed["B"] != "2" {
		t.Errorf("expected B in Removed, got %v", diff.Removed)
	}
}

func TestCompare_Changed(t *testing.T) {
	prev := snapshot.Take(map[string]string{"HOST": "localhost"})
	next := snapshot.Take(map[string]string{"HOST": "remotehost"})
	diff := prev.Compare(next)
	if diff.Changed["HOST"] != "remotehost" {
		t.Errorf("expected HOST in Changed with new value, got %v", diff.Changed)
	}
}

func TestCompare_NoChanges(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	prev := snapshot.Take(env)
	next := snapshot.Take(env)
	diff := prev.Compare(next)
	if diff.HasChanges() {
		t.Errorf("expected no changes, got added=%v removed=%v changed=%v",
			diff.Added, diff.Removed, diff.Changed)
	}
}

func TestDiff_Keys_Sorted(t *testing.T) {
	prev := snapshot.Take(map[string]string{"Z": "old", "M": "same"})
	next := snapshot.Take(map[string]string{"A": "new", "M": "same", "Z": "changed"})
	diff := prev.Compare(next)
	keys := diff.Keys()
	if len(keys) != 2 {
		t.Fatalf("expected 2 changed keys, got %v", keys)
	}
	if keys[0] != "A" || keys[1] != "Z" {
		t.Errorf("expected sorted keys [A Z], got %v", keys)
	}
}

func TestCompare_EmptySnapshots(t *testing.T) {
	prev := snapshot.Take(map[string]string{})
	next := snapshot.Take(map[string]string{})
	diff := prev.Compare(next)
	if diff.HasChanges() {
		t.Error("expected no changes between two empty snapshots")
	}
}
