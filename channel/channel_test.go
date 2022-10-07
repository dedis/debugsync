package channel

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewWithContext(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*1)
	c := NewWithContext[bool](1, ctx, nil, nil)
	require.NotNil(t, c)
}

func TestNewWithTimeout(t *testing.T) {
	c := NewWithTimeout[bool](1, time.Second*2, nil, nil)
	require.NotNil(t, c)
}

/*
func TestChannel(t *testing.T) {
	c := NewWithTimeout[bool](1, time.Second*2, nil, nil)
	channel := c.Channel()
	require.IsType(t, *chan bool, (channel))
}
*/

func TestLen(t *testing.T) {
	c := NewWithTimeout[bool](123, time.Millisecond*1, nil, nil)
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

func TestPushWithTimeoutSuccess(t *testing.T) {

}

func TestPushWithContextSuccess(t *testing.T) {

}

func TestTimedChannelPopWithTimeoutSuccess(t *testing.T) {

}

func TestTimedChannelPopWithContextSuccess(t *testing.T) {

}

func TestTimedChannelPushWithTimeoutFail(t *testing.T) {

}

func TestTimedChannelPushWithContextFail(t *testing.T) {

}

func TestTimedChannelPopWithTimeoutFail(t *testing.T) {

}

func TestTimedChannelPopWithContextFail(t *testing.T) {

}

/*
func TestTimedChannelNormalUse(t *testing.T) {
	c := NewTimedChannel[bool](10, time.Second*1)
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

func TestTimedChannelTimeOut(t *testing.T) {
	var logBuffer bytes.Buffer

	oldLog := Logger
	defer func() {
		Logger = oldLog
	}()

	Logger = zerolog.New(&logBuffer)

	c := NewTimedChannel[bool](3, time.Millisecond*100)
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

	time.Sleep(time.Millisecond * 10) // the log needs some time be generated
	require.True(t, strings.Contains(logBuffer.String(), "channel unblocked"))
}


*/
