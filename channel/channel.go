package channel

import (
	"context"
	"runtime/debug"
	"time"
)

const defaultChannelTimeout = time.Second * 1

type Timed[T any] struct {
	c chan T
}

type Error string

const (
	ErrFailedToSend    = Error("Could not send data on channel.")
	ErrFailedToReceive = Error("Could not receive data from channel")
)

func (e Error) Error() string {
	return string(e)
}

// WithExpiration creates a new channel of the given size and type
func WithExpiration[T any](bufSize int) Timed[T] {
	Logger = Logger.With().Int("size", bufSize).Logger()

	return Timed[T]{
		c: make(chan T, bufSize),
	}
}

// BlockingSendWithContext adds an element in the channel,
// or logs a warning if it fails after the given context
func (c *Timed[T]) BlockingSendWithContext(ctx context.Context, e T) {
	select {
	case c.c <- e:
		return
	case <-ctx.Done():
		Logger.Warn().Msgf("%s %X\n%s", ErrFailedToSend, c.c, string(debug.Stack()))
		c.c <- e
		Logger.Info().Msgf("unblocked channel %X on send", c.c)
	}
}

// BlockingSendWithTimeout adds an element in the channel,
// or logs a warning if it fails after the given timeout
func (c *Timed[T]) BlockingSendWithTimeout(t time.Duration, e T) {
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	c.BlockingSendWithContext(ctx, e)
}

// BlockingSend adds an element in the channel,
// or logs a warning if it fails after default timeout
func (c *Timed[T]) BlockingSend(e T) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultChannelTimeout)
	defer cancel()

	c.BlockingSendWithContext(ctx, e)
}

// BlockingReceiveWithContext removes an element from the channel
// or logs a warning if it fails after the given context
func (c *Timed[T]) BlockingReceiveWithContext(ctx context.Context) T {
	var e T

	select {
	case e = <-c.c:
	case <-ctx.Done():
		Logger.Warn().Msgf("%s %X\n%s", ErrFailedToReceive, c.c, string(debug.Stack()))
		c.c <- e
		Logger.Info().Msgf("unblocked channel %X on receiving", c.c)
	}

	return e
}

// BlockingReceiveWithTimeout removes an element from the channel
// or logs a warning if it fails after the given timeout
func (c *Timed[T]) BlockingReceiveWithTimeout(t time.Duration) T {
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	return c.BlockingReceiveWithContext(ctx)
}

// BlockingReceive removes an element from the channel
// or logs a warning if it fails after the default timeout
func (c *Timed[T]) BlockingReceive() T {
	ctx, cancel := context.WithTimeout(context.Background(), defaultChannelTimeout)
	defer cancel()

	return c.BlockingReceiveWithContext(ctx)
}

// NotBlockingSendWithContext adds an element in the channel,
// or returns an error if it fails in the given context
func (c *Timed[T]) NotBlockingSendWithContext(ctx context.Context, e T) error {
	select {
	case c.c <- e:
		return nil
	case <-ctx.Done():
		return ErrFailedToSend
	}
}

// NotBlockingSendWithTimeout adds an element in the channel,
// or returns an error if it fails after the given timeout
func (c *Timed[T]) NotBlockingSendWithTimeout(t time.Duration, e T) error {
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	return c.NotBlockingSendWithContext(ctx, e)
}

// NotBlockingSend adds an element in the channel,
// or returns an error if it fails after the default timeout
func (c *Timed[T]) NotBlockingSend(e T) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultChannelTimeout)
	defer cancel()

	return c.NotBlockingSendWithContext(ctx, e)
}

// NotBlockingReceiveWithContext removes an element from the channel
// or returns an error if it fails in the given context
func (c *Timed[T]) NotBlockingReceiveWithContext(ctx context.Context) (T, error) {
	var e T

	select {
	case e = <-c.c:
		return e, nil
	case <-ctx.Done():
		return e, ErrFailedToReceive
	}
}

// NotBlockingReceiveWithTimeout removes an element from the channel
// or returns an error if it fails after the given timeout
func (c *Timed[T]) NotBlockingReceiveWithTimeout(t time.Duration) (e T, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	return c.NotBlockingReceiveWithContext(ctx)
}

// NotBlockingReceive removes an element from the channel
// or returns an error if it fails after the default timeout
func (c *Timed[T]) NotBlockingReceive() (e T, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultChannelTimeout)
	defer cancel()

	return c.NotBlockingReceiveWithContext(ctx)
}

// Len gives the current number of elements in the channel
func (c *Timed[T]) Len() int {
	return len(c.c)
}

// Channel returns the raw channel used
func (c *Timed[T]) Channel() chan T {
	return c.c
}
