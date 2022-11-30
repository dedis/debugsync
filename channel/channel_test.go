package channel

import (
	"bytes"
	"context"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

var originalLogger = Logger

// setupLogger is a helper function to use a testable logger
func setupLogger() *bytes.Buffer {
	b := new(bytes.Buffer)
	Logger = zerolog.New(b)

	return b
}

// restoreLogger is a helper function to restore the original logger
func restoreLogger() {
	Logger = originalLogger
}

func TestWithExpiration(t *testing.T) {
	c := WithExpiration[bool](1)
	require.NotNil(t, c)
}

func TestBlockingSendWithContextSuccess(t *testing.T) {
	l := setupLogger()
	defer restoreLogger()

	c := WithExpiration[int](1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	c.BlockingSendWithContext(ctx, 0)
	require.False(t, strings.Contains(l.String(), ErrFailedToSend.Error()))
}

func TestBlockingSendWithTimeoutSuccess(t *testing.T) {
	l := setupLogger()
	defer restoreLogger()

	c := WithExpiration[int](1)

	c.BlockingSendWithTimeout(time.Millisecond, 0)
	require.False(t, strings.Contains(l.String(), ErrFailedToSend.Error()))
}

func TestBlockingSendSuccess(t *testing.T) {
	l := setupLogger()
	defer restoreLogger()

	c := WithExpiration[int](1)

	c.BlockingSend(0)
	require.False(t, strings.Contains(l.String(), ErrFailedToSend.Error()))
}

func TestBlockingReceiveWithContextSuccess(t *testing.T) {
	l := setupLogger()
	defer restoreLogger()

	c := WithExpiration[int](1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	expectedValue := -1
	c.BlockingSend(expectedValue)
	v := c.BlockingReceiveWithContext(ctx)
	require.False(t, strings.Contains(l.String(), ErrFailedToReceive.Error()))
	require.Equal(t, expectedValue, v)
}

func TestBlockingReceiveWithTimeoutSuccess(t *testing.T) {
	l := setupLogger()
	defer restoreLogger()

	c := WithExpiration[int](1)

	expectedValue := 2
	c.BlockingSend(expectedValue)
	v := c.BlockingReceiveWithTimeout(time.Millisecond)
	require.False(t, strings.Contains(l.String(), ErrFailedToReceive.Error()))
	require.Equal(t, expectedValue, v)
}

func TestBlockingReceiveSuccess(t *testing.T) {
	l := setupLogger()
	defer restoreLogger()

	c := WithExpiration[int](1)

	expectedValue := -7
	c.BlockingSend(expectedValue)
	v := c.BlockingReceive()
	require.False(t, strings.Contains(l.String(), ErrFailedToReceive.Error()))
	require.Equal(t, expectedValue, v)
}

func TestPushFail(t *testing.T) {
	l := setupLogger()
	defer restoreLogger()

	c := WithExpiration[int](1)

	c.BlockingSendWithTimeout(time.Millisecond, 0)

	go func() {
		c.BlockingSendWithTimeout(time.Millisecond, 0)
	}()

	// need a looong time on Windows to see the logs in the buffer
	time.Sleep(time.Millisecond * 100)
	require.True(t, strings.Contains(l.String(), ErrFailedToSend.Error()))
}

func TestBlockingReceiveFail(t *testing.T) {
	l := setupLogger()
	defer restoreLogger()

	c := WithExpiration[int](1)

	go func() {
		c.BlockingReceiveWithTimeout(time.Millisecond)
	}()

	// need a looong time on Windows to see the logs in the buffer
	time.Sleep(time.Millisecond * 100)
	require.True(t, strings.Contains(l.String(), ErrFailedToReceive.Error()))
}

func TestChannel(t *testing.T) {
	c := WithExpiration[int](1)
	channel := c.Channel()

	const data = 12345
	c.BlockingSend(data)

	value := <-channel
	require.Equal(t, data, value)
}

func TestLen(t *testing.T) {
	c := WithExpiration[bool](3)
	require.Equal(t, 0, c.Len())

	c.BlockingSend(true)
	c.BlockingSend(false)
	require.Equal(t, 2, c.Len())

	c.BlockingReceive()
	require.Equal(t, 1, c.Len())

	c.BlockingReceive()
	require.Equal(t, 0, c.Len())
}
