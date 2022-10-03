package sync

import (
	_sync "sync"
)

// Map is like a Go map[interface{}]interface{} but is safe for concurrent use
// by multiple goroutines without additional locking or coordination. Loads,
// stores, and deletes run in amortized constant time.
type Map struct {
	_sync.Map
}
