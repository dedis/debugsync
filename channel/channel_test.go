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

func TestCryChanPassing(t *testing.T) {
	c := NewCryChan[bool](10, time.Second*1)
	require.NotNil(t, c)

	c.Push(true)
	c.Push(false)
	c.Push(false)
	require.Equal(t, 3, c.Len())

	b, err := c.Pop(context.Background())
	require.NoError(t, err)
	require.True(t, b)
	require.Equal(t, 2, c.Len())

	b, err = c.Pop(context.Background())
	require.NoError(t, err)
	require.False(t, b)
	require.Equal(t, 1, c.Len())

	b, err = c.Pop(context.Background())
	require.NoError(t, err)
	require.False(t, b)
	require.Equal(t, 0, c.Len())
}

func TestCryChanTimedOut(t *testing.T) {
	var logBuffer bytes.Buffer

	oldLog := Logger
	defer func() {
		Logger = oldLog
	}()

	Logger = zerolog.New(&logBuffer)

	c := NewCryChan[bool](3, time.Millisecond*100)
	require.NotNil(t, c)

	c.Push(true)
	c.Push(false)
	c.Push(false)
	require.Equal(t, 3, c.Len())

	go func() {
		c.Push(false)
	}()

	time.Sleep(time.Millisecond * 200)
	require.True(t, strings.Contains(logBuffer.String(), "channel blocking"))
	require.False(t, strings.Contains(logBuffer.String(), "channel unblocked"))

	b, err := c.Pop(context.Background())
	require.NoError(t, err)
	require.True(t, b)

	time.Sleep(time.Millisecond * 10)
	require.True(t, strings.Contains(logBuffer.String(), "channel unblocked"))
}
