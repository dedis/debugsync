package sync

import (
	_sync "sync"
)

// Cond implements a condition variable, a rendezvous point for goroutines
// waiting for or announcing the occurrence of an event.
type Cond struct {
	_sync.Cond
}
