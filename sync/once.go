package sync

import (
	_sync "sync"
)

// Once is an object that will perform exactly one action.
type Once struct {
	_sync.Once
}
