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

func TestSendWithContextSuccess(t *testing.T) {
	l := setupLogger()
	defer restoreLogger()

	c := WithExpiration[int](1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	c.SendWithContext(ctx, 0)
	require.False(t, strings.Contains(l.String(), ErrFailedToSend.Error()))
}

func TestSendWithTimeoutSuccess(t *testing.T) {
	l := setupLogger()
	defer restoreLogger()

	c := WithExpiration[int](1)

	c.SendWithTimeout(time.Millisecond, 0)
	require.False(t, strings.Contains(l.String(), ErrFailedToSend.Error()))
}

func TestSendSuccess(t *testing.T) {
	l := setupLogger()
	defer restoreLogger()

	c := WithExpiration[int](1)

	c.Send(0)
	require.False(t, strings.Contains(l.String(), ErrFailedToSend.Error()))
}

func TestReceiveWithContextSuccess(t *testing.T) {
	l := setupLogger()
	defer restoreLogger()

	c := WithExpiration[int](1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	expectedValue := -1
	c.Send(expectedValue)
	v := c.ReceiveWithContext(ctx)
	require.False(t, strings.Contains(l.String(), ErrFailedToReceive.Error()))
	require.Equal(t, expectedValue, v)
}

func TestReceiveWithTimeoutSuccess(t *testing.T) {
	l := setupLogger()
	defer restoreLogger()

	c := WithExpiration[int](1)

	expectedValue := 2
	c.Send(expectedValue)
	v := c.ReceiveWithTimeout(time.Millisecond)
	require.False(t, strings.Contains(l.String(), ErrFailedToReceive.Error()))
	require.Equal(t, expectedValue, v)
}

func TestReceiveSuccess(t *testing.T) {
	l := setupLogger()
	defer restoreLogger()

	c := WithExpiration[int](1)

	expectedValue := -7
	c.Send(expectedValue)
	v := c.Receive()
	require.False(t, strings.Contains(l.String(), ErrFailedToReceive.Error()))
	require.Equal(t, expectedValue, v)
}

func TestSendFail(t *testing.T) {
	l := setupLogger()
	defer restoreLogger()

	c := WithExpiration[int](1)

	c.SendWithTimeout(time.Millisecond, 0)

	go func() {
		c.SendWithTimeout(time.Millisecond, 0)
	}()

	// need a looong time on Windows to see the logs in the buffer
	time.Sleep(time.Millisecond * 100)
	require.True(t, strings.Contains(l.String(), ErrFailedToSend.Error()))
}

func TestReceiveFail(t *testing.T) {
	l := setupLogger()
	defer restoreLogger()

	c := WithExpiration[int](1)

	go func() {
		c.ReceiveWithTimeout(time.Millisecond)
	}()

	// need a looong time on Windows to see the logs in the buffer
	time.Sleep(time.Millisecond * 100)
	require.True(t, strings.Contains(l.String(), ErrFailedToReceive.Error()))
}

func TestNonBlockingSendWithContextSuccess(t *testing.T) {
	defer restoreLogger()

	c := WithExpiration[int](1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	err := c.NonBlockingSendWithContext(ctx, 0)
	require.NoError(t, err)
}

func TestNonBlockingSendWithTimeoutSuccess(t *testing.T) {
	defer restoreLogger()

	c := WithExpiration[int](1)

	err := c.NonBlockingSendWithTimeout(time.Millisecond, 0)
	require.NoError(t, err)
}

func TestNonBlockingSendSuccess(t *testing.T) {
	defer restoreLogger()

	c := WithExpiration[int](1)

	err := c.NonBlockingSend(0)
	require.NoError(t, err)
}

func TestNonBlockingReceiveWithContextSuccess(t *testing.T) {
	defer restoreLogger()

	c := WithExpiration[int](1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	expectedValue := -1
	err := c.NonBlockingSend(expectedValue)
	require.NoError(t, err)

	v, err := c.NonBlockingReceiveWithContext(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedValue, v)
}

func TestNonBlockingReceiveWithTimeoutSuccess(t *testing.T) {
	defer restoreLogger()

	c := WithExpiration[int](1)

	expectedValue := 2
	err := c.NonBlockingSend(expectedValue)
	require.NoError(t, err)

	v, err := c.NonBlockingReceiveWithTimeout(time.Millisecond)
	require.NoError(t, err)
	require.Equal(t, expectedValue, v)
}

func TestNonBlockingReceiveSuccess(t *testing.T) {
	defer restoreLogger()

	c := WithExpiration[int](1)

	expectedValue := -7
	err := c.NonBlockingSend(expectedValue)
	require.NoError(t, err)

	v, err := c.NonBlockingReceive()
	require.NoError(t, err)
	require.Equal(t, expectedValue, v)
}

func TestNonBlockingSendFail(t *testing.T) {
	defer restoreLogger()

	c := WithExpiration[int](1)

	err := c.NonBlockingSendWithTimeout(time.Millisecond, 0)
	require.NoError(t, err)

	go func() {
		err := c.NonBlockingSendWithTimeout(time.Millisecond, 0)
		require.Equal(t, err, ErrFailedToSend)
	}()

	// need a looong time on Windows to see the logs in the buffer
	time.Sleep(time.Millisecond * 100)
}

func TestNonBlockingReceiveFail(t *testing.T) {
	defer restoreLogger()

	c := WithExpiration[int](1)

	go func() {
		_, err := c.NonBlockingReceiveWithTimeout(time.Millisecond)
		require.Error(t, err, ErrFailedToReceive)
	}()

	// need a looong time on Windows to see the logs in the buffer
	time.Sleep(time.Millisecond * 100)
}

func TestChannel(t *testing.T) {
	c := WithExpiration[int](1)
	channel := c.Channel()

	const data = 12345
	c.Send(data)

	value := <-channel
	require.Equal(t, data, value)
}

func TestLen(t *testing.T) {
	c := WithExpiration[bool](3)
	require.Equal(t, 0, c.Len())

	c.Send(true)
	c.Send(false)
	require.Equal(t, 2, c.Len())

	c.Receive()
	require.Equal(t, 1, c.Len())

	c.Receive()
	require.Equal(t, 0, c.Len())
}
