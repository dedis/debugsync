package sync

import (
	_sync "sync"
)

// A Pool is a set of temporary objects that may be individually saved
// and retrieved.
type Pool struct {
	_sync.Pool
}
