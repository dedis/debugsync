package channel

import (
	"context"
	"runtime/debug"
	"time"
)

const defaultChannelTimeout = time.Second * 1

var BlockedPush = string("Push blocked on channel: ")
var UnblockedPush = string("Push unblocked on channel: ")
var BlockedPop = string("Pop blocked on channel: ")
var UnblockedPop = string("Pop unblocked on channel: ")

type Timed[T any] struct {
	c chan T
}

// WithExpiration creates a new channel of the given size and type
func WithExpiration[T any](bufSize int) Timed[T] {
	Logger = Logger.With().Int("size", bufSize).Logger()

	return Timed[T]{
		c: make(chan T, bufSize),
	}
}

// PushWithContext adds an element in the channel,
// or logs a warning if it fails after the given context
func (c *Timed[T]) PushWithContext(ctx context.Context, e T) error {
	select {
	case c.c <- e:
		return nil
	case <-ctx.Done():
		Logger.Warn().Msgf("%s %X\n%s", BlockedPush, c.c, string(debug.Stack()))
		c.c <- e
		Logger.Info().Msgf("%s %X", UnblockedPush, c.c)
		return ctx.Err()
	}
}

// PushWithTimeout adds an element in the channel,
// or logs a warning if it fails after the given timeout
func (c *Timed[T]) PushWithTimeout(t time.Duration, e T) error {
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	return c.PushWithContext(ctx, e)
}

// Push adds an element in the channel,
// or logs a warning if it fails after default timeout
func (c *Timed[T]) Push(e T) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultChannelTimeout)
	defer cancel()

	return c.PushWithContext(ctx, e)
}

// PopWithContext removes an element from the channel
// or logs a warning if it fails after the given context
func (c *Timed[T]) PopWithContext(ctx context.Context) (T, error) {
	var e T

	select {
	case e = <-c.c:
	case <-ctx.Done():
		Logger.Warn().Msgf("%s %X\n%s", BlockedPop, c.c, string(debug.Stack()))
		c.c <- e
		Logger.Info().Msgf("%s %X", UnblockedPop, c.c)
		return e, ctx.Err()
	}
	return e, nil
}

// PopWithTimeout removes an element from the channel
// or logs a warning if it fails after the given timeout
func (c *Timed[T]) PopWithTimeout(t time.Duration) (T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	return c.PopWithContext(ctx)
}

// Pop removes an element from the channel
// or logs a warning if it fails after the default timeout
func (c *Timed[T]) Pop() (T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultChannelTimeout)
	defer cancel()

	return c.PopWithContext(ctx)
}

// Len gives the current number of elements in the channel
func (c *Timed[T]) Len() int {
	return len(c.c)
}

// Channel returns the raw channel used
func (c *Timed[T]) Channel() chan T {
	return c.c
}
