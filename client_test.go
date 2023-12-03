package pokedex

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	client, err := New()
	require.NoError(t, err)
	require.NotNil(t, client)

	client, err = NewWithOptions(Options{
		BaseURL:          "https://example.com",
		Timeout:          10 * time.Second,
		CacheMaximumSize: 1,
		CacheTTL:         1 * time.Minute,
	})
	require.NoError(t, err)
	require.NotNil(t, client)
}

func TestNewClient_BadURL(t *testing.T) {
	client, err := NewWithOptions(Options{
		BaseURL: "https://\nexample.com",
	})
	require.ErrorContains(t, err, "parse")
	require.Nil(t, client)
}
