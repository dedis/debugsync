package channel

import (
	"context"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

var wasCalled bool
var stringFromCallback string

func testCallBack(s string) {
	wasCalled = true
	stringFromCallback = s
}

func TestNewWithTimeout(t *testing.T) {
	c := NewWithTimeout[bool](time.Millisecond, 1, nil, nil)
	require.NotNil(t, c)
}

func TestNewWithContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	c := NewWithContext[bool](ctx, 1, nil, nil)
	require.NotNil(t, c)
}

func TestPushWithTimeoutSuccess(t *testing.T) {
	wasCalled = false
	c := NewWithTimeout[bool](time.Millisecond, 1, testCallBack, nil)

	c.Push(false)
	require.False(t, wasCalled)
}

func TestPushWithContextSuccess(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	wasCalled = false
	c := NewWithContext[bool](ctx, 1, testCallBack, nil)

	c.Push(false)
	require.False(t, wasCalled)
}

func TestTimedChannelPopWithTimeoutSuccess(t *testing.T) {
	wasCalled = false
	c := NewWithTimeout[bool](time.Millisecond, 1, testCallBack, nil)
	c.Push(true)

	b, err := c.Pop()
	require.NoError(t, err)
	require.True(t, b)
	require.False(t, wasCalled)
}

func TestTimedChannelPopWithContextSuccess(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	wasCalled = false
	c := NewWithContext[bool](ctx, 1, testCallBack, nil)
	c.Push(true)

	b, err := c.Pop()
	require.NoError(t, err)
	require.True(t, b)
	require.False(t, wasCalled)
}

func TestTimedChannelPushWithTimeoutFail(t *testing.T) {
	wasCalled = false
	c := NewWithTimeout[bool](time.Millisecond, 1, testCallBack, nil)

	c.Push(false)
	c.Push(false)
	require.True(t, wasCalled)
	require.True(t, strings.Contains(stringFromCallback, "stack"))
}

func TestTimedChannelPushWithContextFail(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	wasCalled = false
	c := NewWithContext[int](ctx, 1, testCallBack, nil)

	c.Push(1)
	c.Push(1)
	require.True(t, wasCalled)
	require.True(t, strings.Contains(stringFromCallback, "stack"))
}

func TestTimedChannelPopWithTimeoutFail(t *testing.T) {
	wasCalled = false
	c := NewWithTimeout[bool](time.Millisecond, 1, nil, testCallBack)

	_, err := c.Pop()
	require.Error(t, err)
	require.True(t, wasCalled)
	require.True(t, strings.Contains(stringFromCallback, "stack"))
}

func TestTimedChannelPopWithContextFail(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	wasCalled = false
	c := NewWithContext[int](ctx, 1, nil, testCallBack)

	_, err := c.Pop()
	require.Error(t, err)
	require.True(t, wasCalled)
	require.True(t, strings.Contains(stringFromCallback, "stack"))
}

func TestChannel(t *testing.T) {
	c := NewWithTimeout[int](time.Millisecond, 1, nil, nil)
	channel := c.Channel()

	const data = 12345
	c.Push(data)

	value := <-*channel
	require.Equal(t, data, value)
}

func TestLen(t *testing.T) {
	c := NewWithTimeout[bool](time.Millisecond, 3, nil, nil)
	require.Equal(t, 0, c.Len())

	c.Push(true)
	c.Push(false)
	require.Equal(t, 2, c.Len())

	_, err := c.Pop()
	require.NoError(t, err)
	require.Equal(t, 1, c.Len())

	_, err = c.Pop()
	require.NoError(t, err)
	require.Equal(t, 0, c.Len())

	_, err = c.Pop()
	require.Error(t, err)
	require.Equal(t, 0, c.Len())
}
