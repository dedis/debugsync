package channel

import (
	"context"
	"github.com/rs/zerolog"
	"runtime/debug"
	"time"
)

const defaultChannelTimeout = time.Second * 1

type Timed[T any] struct {
	c   chan T
	log zerolog.Logger
}

type Error string

const (
	ErrFailedToSend    = Error("Could not send data on channel.")
	ErrFailedToReceive = Error("Could not receive data from channel.")
)

func (e Error) Error() string {
	return string(e)
}

// WithExpiration creates a new channel of the given size and type
func WithExpiration[T any](bufSize int) Timed[T] {
	return Timed[T]{
		c:   make(chan T, bufSize),
		log: Logger.With().Int("size", bufSize).Logger(),
	}
}

// SendWithContext adds an element in the channel,
// or logs a warning if it fails in the given context.
// Note: this is a blocking call as it waits on a channel.
func (c *Timed[T]) SendWithContext(ctx context.Context, e T) {
	select {
	case c.c <- e:
		return
	case <-ctx.Done():
		c.log.Warn().Msgf("%s %X\n%s", ErrFailedToSend, c.c, string(debug.Stack()))
		c.c <- e
		c.log.Info().Msgf("unblocked channel %X on send", c.c)
	}
}

// SendWithTimeout adds an element in the channel,
// or logs a warning if it fails after the given timeout.
// Note: this is a blocking call as it waits on a channel.
func (c *Timed[T]) SendWithTimeout(t time.Duration, e T) {
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	c.SendWithContext(ctx, e)
}

// Send adds an element in the channel,
// or logs a warning if it fails after default timeout.
// Note: this is a blocking call as it waits on a channel.
func (c *Timed[T]) Send(e T) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultChannelTimeout)
	defer cancel()

	c.SendWithContext(ctx, e)
}

// ReceiveWithContext removes an element from the channel
// or logs a warning if it fails in the given context.
// Note: this is a blocking call as it waits on a channel.
func (c *Timed[T]) ReceiveWithContext(ctx context.Context) T {
	var e T

	select {
	case e = <-c.c:
	case <-ctx.Done():
		c.log.Warn().Msgf("%s %X\n%s", ErrFailedToReceive, c.c, string(debug.Stack()))
		c.c <- e
		c.log.Info().Msgf("unblocked channel %X on receiving", c.c)
	}

	return e
}

// ReceiveWithTimeout removes an element from the channel
// or logs a warning if it fails after the given timeout.
// Note: this is a blocking call as it waits on a channel.
func (c *Timed[T]) ReceiveWithTimeout(t time.Duration) T {
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	return c.ReceiveWithContext(ctx)
}

// Receive removes an element from the channel
// or logs a warning if it fails after the default timeout.
// Note: this is a blocking call as it waits on a channel.
func (c *Timed[T]) Receive() T {
	ctx, cancel := context.WithTimeout(context.Background(), defaultChannelTimeout)
	defer cancel()

	return c.ReceiveWithContext(ctx)
}

// NonBlockingSendWithContext adds an element in the channel,
// or returns an error if it fails in the given context.
func (c *Timed[T]) NonBlockingSendWithContext(ctx context.Context, e T) error {
	select {
	case c.c <- e:
		return nil
	case <-ctx.Done():
		return ErrFailedToSend
	}
}

// NonBlockingSendWithTimeout adds an element in the channel,
// or returns an error if it fails after the given timeout.
func (c *Timed[T]) NonBlockingSendWithTimeout(t time.Duration, e T) error {
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	return c.NonBlockingSendWithContext(ctx, e)
}

// NonBlockingSend adds an element in the channel,
// or returns an error if it fails after the default timeout.
func (c *Timed[T]) NonBlockingSend(e T) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultChannelTimeout)
	defer cancel()

	return c.NonBlockingSendWithContext(ctx, e)
}

// NonBlockingReceiveWithContext removes an element from the channel
// or returns an error if it fails in the given context.
func (c *Timed[T]) NonBlockingReceiveWithContext(ctx context.Context) (T, error) {
	var e T

	select {
	case e = <-c.c:
		return e, nil
	case <-ctx.Done():
		return e, ErrFailedToReceive
	}
}

// NonBlockingReceiveWithTimeout removes an element from the channel
// or returns an error if it fails after the given timeout
func (c *Timed[T]) NonBlockingReceiveWithTimeout(t time.Duration) (e T, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	return c.NonBlockingReceiveWithContext(ctx)
}

// NonBlockingReceive removes an element from the channel
// or returns an error if it fails after the default timeout
func (c *Timed[T]) NonBlockingReceive() (e T, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultChannelTimeout)
	defer cancel()

	return c.NonBlockingReceiveWithContext(ctx)
}

// Len gives the current number of elements in the channel
func (c *Timed[T]) Len() int {
	return len(c.c)
}

// Channel returns the raw channel used
func (c *Timed[T]) Channel() chan T {
	return c.c
}
