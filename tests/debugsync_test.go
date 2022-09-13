package tests

import (
	"github.com/dedis/debugsync"
	"github.com/stretchr/testify/require"
	"testing"
)

func criticalSection() {
	debugsync.Logger.Info().Msg("entering critical section")
	for i := 0; i < 10; i++ {
		debugsync.Logger.Debug().Msgf("counting %v", i)
	}
	debugsync.Logger.Info().Msg("leaving critical section")

}

func TestMutex(t *testing.T) {
	m := debugsync.Mutex{}

	m.Lock()
	criticalSection()
	m.Unlock()

	require.True(t, m.TryLock(), "Trylock() shouldn't fail here")
	criticalSection()
	m.Unlock()
}

func TestRwMutex(t *testing.T) {
	m := debugsync.RWMutex{}

	m.Lock()
	criticalSection()
	m.Unlock()

	require.True(t, m.TryLock(), "Trylock() shouldn't fail here")
	//require.Panics(t, m.RUnlock, "RUnlock() did not panic")
	criticalSection()
	m.Unlock()

	m.RLock()
	criticalSection()
	m.RUnlock()

	require.True(t, m.TryRLock(), "Trylock() shouldn't fail here")
	criticalSection()
	require.Panics(t, m.Unlock, "Unlock() did not panic")
	m.RUnlock()
}

func TestWaitGroup(t *testing.T) {
	wg := debugsync.WaitGroup{}

	wg.Add(1)

	go func() {
		wg.Done()

		//require.Panics(t, wg.Done, "Done failed to panic")
	}()

	wg.Wait()
}
