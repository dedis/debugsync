package channel

import (
	"context"
	"golang.org/x/xerrors"
	"time"
)

const defaultChannelTimeout = time.Second * 1

var FailedPush = xerrors.New("channel blocked on Push with call stack")
var FailedPop = xerrors.New("channel blocked on Push with call stack")

type Timed[T any] struct {
	c chan T
}

// WithExpiration creates a new channel of the given size, type and timeout
// Note: if the callbacks are not nil, they are called when the timeout expires
func WithExpiration[T any](bufSize int) Timed[T] {
	Logger = Logger.With().Int("size", bufSize).Logger()

	return Timed[T]{
		c: make(chan T, bufSize),
	}
}

// PushWithContext adds an element in the channel,
// or logs a warning if it fails after the given context
func (c *Timed[T]) PushWithContext(ctx context.Context, e T) {
	select {
	case c.c <- e:
		return
	case <-ctx.Done():
		Logger.Warn().Stack().Msg(FailedPush.Error())
		c.c <- e
	}
}

// PushWithTimeout adds an element in the channel,
// or logs a warning if it fails after the given timeout
func (c *Timed[T]) PushWithTimeout(t time.Duration, e T) {
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	c.PushWithContext(ctx, e)
}

// Push adds an element in the channel,
// or logs a warning if it fails after default timeout
func (c *Timed[T]) Push(e T) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultChannelTimeout)
	defer cancel()

	c.PushWithContext(ctx, e)
}

// PopWithContext removes an element from the channel
// or logs a warning if it fails after the given context timeout
func (c *Timed[T]) PopWithContext(ctx context.Context) T {
	var e T

	select {
	case e = <-c.c:
	case <-ctx.Done():
		Logger.Warn().Stack().Msg(FailedPop.Error())
		e = <-c.c
	}

	return e
}

// PopWithTimeout removes an element from the channel
// or logs a warning if it fails after the given timeout
func (c *Timed[T]) PopWithTimeout(t time.Duration) T {
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	return c.PopWithContext(ctx)
}

// Pop removes an element from the channel
// or logs a warning if it fails after the default timeout
func (c *Timed[T]) Pop() T {
	ctx, cancel := context.WithTimeout(context.Background(), defaultChannelTimeout)
	defer cancel()

	return c.PopWithContext(ctx)
}

// Len gives the current number of elements in the channel
func (c *Timed[T]) Len() int {
	return len(c.c)
}

// Channel returns the raw channel used
func (c *Timed[T]) Channel() *chan T {
	return &c.c
}
