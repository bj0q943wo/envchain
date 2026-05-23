// Package watcher implements lightweight polling-based file watching
// for .env files used in envchain workflows.
//
// It is intentionally simple: rather than relying on OS-level inotify
// or FSEvents (which require third-party dependencies), it polls
// modification times at a configurable interval.
//
// Typical usage:
//
//	w := watcher.New(500 * time.Millisecond)
//	_ = w.Add(".env")
//	_ = w.Add(".env.local")
//	w.Start()
//	defer w.Stop()
//
//	for ev := range w.Events {
//		fmt.Println("changed:", ev.Path)
//	}
package watcher
