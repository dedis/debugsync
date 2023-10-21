package sync

import (
	"time"
)

var Timeout = 10 * time.Second

func startLockTimer(msg string, stack []byte) chan struct{} {
	done := make(chan struct{})

	go func(s []byte) {
		select {
		case <-time.After(Timeout):
			Logger.Error().Msgf("%v : %v", msg, string(s))
			return
		case <-done:
			return
		}
	}(stack)

	return done
}
