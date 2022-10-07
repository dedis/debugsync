package channel

import (
	"context"
	"errors"
	"runtime/debug"
	"time"
)

type TimedChannel[T any] struct {
	c      chan T
	t      time.Duration
	ctx    context.Context
	pushCb func(string)
	popCb  func(string)
}

// NewWithTimeout creates a new channel of the given size, type and timeout
// Note: if the callbacks are not nil, they are called when the timeout expires
func NewWithTimeout[T any](
	bufSize int,
	timeout time.Duration,
	pushCallback func(string),
	popCallback func(string)) TimedChannel[T] {
	Logger = Logger.With().Int("size", bufSize).Logger()

	return TimedChannel[T]{
		c:      make(chan T, bufSize),
		t:      timeout,
		ctx:    context.Background(),
		pushCb: pushCallback,
		popCb:  popCallback,
	}
}

// NewWithContext creates a new channel of the given size, type and context
// Note: if the callbacks are not nil, they are called when the context expires
func NewWithContext[T any](
	bufSize int,
	ctx context.Context,
	pushCallback func(string),
	popCallback func(string)) TimedChannel[T] {
	Logger = Logger.With().Int("size", bufSize).Logger()

	return TimedChannel[T]{
		c:      make(chan T, bufSize),
		t:      0,
		ctx:    ctx,
		pushCb: pushCallback,
		popCb:  popCallback,
	}
}

// Push adds an element in the channel,
// or calls the pushCb if it fails after the given timeout
func (c *TimedChannel[T]) Push(e T) {
	select {
	case c.c <- e:
		return
	case <-c.ctx.Done():
	case <-time.After(c.t):
	}

	// process timeout
	s := string(debug.Stack())
	if c.pushCb != nil {
		c.pushCb(s)
	} else {
		Logger.Warn().Str("stack", s)
	}
}

// Pop removes an element from the channel
func (c *TimedChannel[T]) Pop() (t T, err error) {
	select {
	case el := <-c.c:
		return el, nil
	case <-c.ctx.Done():
	case <-time.After(c.t):
	}

	// process timeout
	s := string(debug.Stack())
	if c.popCb != nil {
		c.popCb(s)
	} else {
		Logger.Warn().Str("stack", s)
	}

	return t, errors.New(s)
}

// Len gives the current number of elements in the channel
func (c *TimedChannel[T]) Len() int {
	return len(c.c)
}

// Channel returns the raw channel used
func (c *TimedChannel[T]) Channel() *chan T {
	return &c.c
}
