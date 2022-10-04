package sync

import (
	"time"
)

const mutexTimeout = 30 * time.Minute

func startLockTimer(msg string, stack []byte) chan struct{} {
	done := make(chan struct{})

	go func(s []byte) {
		select {
		case <-time.After(mutexTimeout):
			Logger.Error().Msgf("%v : %v", msg, string(s))
			return
		case <-done:
			return
		}
	}(stack)

	return done
}
