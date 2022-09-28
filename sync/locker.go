package sync

// Locker is only defined for compatibility reasons
type Locker interface {
	Lock()
	Unlock()
}
