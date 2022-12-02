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

func TestBlockingSendFail(t *testing.T) {
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

func TestNotBlockingSendWithContextSuccess(t *testing.T) {
	defer restoreLogger()

	c := WithExpiration[int](1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	err := c.NotBlockingSendWithContext(ctx, 0)
	require.NoError(t, err)
}

func TestNotBlockingSendWithTimeoutSuccess(t *testing.T) {
	defer restoreLogger()

	c := WithExpiration[int](1)

	err := c.NotBlockingSendWithTimeout(time.Millisecond, 0)
	require.NoError(t, err)
}

func TestNotBlockingSendSuccess(t *testing.T) {
	defer restoreLogger()

	c := WithExpiration[int](1)

	err := c.NotBlockingSend(0)
	require.NoError(t, err)
}

func TestNotBlockingReceiveWithContextSuccess(t *testing.T) {
	defer restoreLogger()

	c := WithExpiration[int](1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	expectedValue := -1
	err := c.NotBlockingSend(expectedValue)
	require.NoError(t, err)

	v, err := c.NotBlockingReceiveWithContext(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedValue, v)
}

func TestNotBlockingReceiveWithTimeoutSuccess(t *testing.T) {
	defer restoreLogger()

	c := WithExpiration[int](1)

	expectedValue := 2
	err := c.NotBlockingSend(expectedValue)
	require.NoError(t, err)

	v, err := c.NotBlockingReceiveWithTimeout(time.Millisecond)
	require.NoError(t, err)
	require.Equal(t, expectedValue, v)
}

func TestNotBlockingReceiveSuccess(t *testing.T) {
	defer restoreLogger()

	c := WithExpiration[int](1)

	expectedValue := -7
	err := c.NotBlockingSend(expectedValue)
	require.NoError(t, err)

	v, err := c.NotBlockingReceive()
	require.NoError(t, err)
	require.Equal(t, expectedValue, v)
}

func TestNotBlockingSendFail(t *testing.T) {
	defer restoreLogger()

	c := WithExpiration[int](1)

	err := c.NotBlockingSendWithTimeout(time.Millisecond, 0)
	require.NoError(t, err)

	go func() {
		err := c.NotBlockingSendWithTimeout(time.Millisecond, 0)
		require.Equal(t, err, ErrFailedToSend)
	}()

	// need a looong time on Windows to see the logs in the buffer
	time.Sleep(time.Millisecond * 100)
}

func TestNotBlockingReceiveFail(t *testing.T) {
	defer restoreLogger()

	c := WithExpiration[int](1)

	go func() {
		_, err := c.NotBlockingReceiveWithTimeout(time.Millisecond)
		require.Error(t, err, ErrFailedToReceive)
	}()

	// need a looong time on Windows to see the logs in the buffer
	time.Sleep(time.Millisecond * 100)
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
