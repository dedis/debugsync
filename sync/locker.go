package sync

import (
	_sync "sync"
)

// A Locker represents an object that can be locked and unlocked.
type Locker interface {
	_sync.Locker
}
