package channel

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCryChan_Passing(t *testing.T) {
	c := NewCryChan[bool](10)
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

	/*
		b, err = c.Pop(context.Background())
		require.Error(t, err)
		require.Equal(t, 0, c.Len())
	*/
}
