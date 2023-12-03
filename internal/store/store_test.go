package store

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	c, err := NewCache(CacheOptions{
		MaximumSize: 1 << 27,
		TTL:         10 * time.Minute,
		debug:       true,
	})
	require.NoError(t, err)
	defer c.Close()

	c.Set("https://example.com/foo/bar", []byte("foobar"))
	b, ok := c.Get("https://example.com/foo/bar")
	require.True(t, ok)
	require.Equal(t, []byte("foobar"), b)

	// Accumulate a number of cache hits so we can assert this against the
	// metrics later.
	for i := 0; i < 10; i++ {
		_, ok = c.Get("https://example.com/foo/bar")
		require.True(t, ok)
	}

	// Try and fetch a record not in the cache.
	b, ok = c.Get("https://example.com/missing")
	require.False(t, ok)
	require.Nil(t, b)

	// Add a bunch of stuff to the cache then fetch it.
	for i := 1; i <= 5; i++ {
		u := fmt.Sprintf("https://example.com/foo/%d", i)

		c.Set(u, []byte("foobar"))

		b, ok = c.Get(u)
		require.True(t, ok)
		require.Equal(t, []byte("foobar"), b)
	}

	require.Equal(t, uint64(6), c.cache.Metrics.KeysAdded())
	require.Equal(t, uint64(16), c.cache.Metrics.Hits())
	require.Equal(t, uint64(1), c.cache.Metrics.Misses())
}
