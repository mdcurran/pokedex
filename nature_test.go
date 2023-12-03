package pokedex

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mdcurran/pokedex/internal/faker"
	"github.com/stretchr/testify/require"
)

func TestGetNature(t *testing.T) {
	ctx := context.Background()

	fixture, err := json.Marshal(faker.NewFaker().GenerateNature())
	require.NoError(t, err)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", fixture)
	}))
	t.Cleanup(srv.Close)

	sdk, err := NewWithOptions(Options{
		BaseURL:          srv.URL,
		Timeout:          5 * time.Second,
		CacheMaximumSize: 1 << 27,
		CacheTTL:         10 * time.Second,
	})
	require.NoError(t, err)
	t.Cleanup(sdk.Close)

	res, err := sdk.GetNature(ctx, GetRequest{ID: 1})
	require.NoError(t, err)

	serialised, err := json.Marshal(res.Nature)
	require.NoError(t, err)
	require.JSONEq(t, string(fixture), string(serialised))
}

func TestNature_NotFound(t *testing.T) {
	ctx := context.Background()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Not Found")
	}))
	t.Cleanup(srv.Close)

	sdk, err := NewWithOptions(Options{
		BaseURL:          srv.URL,
		Timeout:          5 * time.Second,
		CacheMaximumSize: 1 << 27,
		CacheTTL:         10 * time.Second,
	})
	require.NoError(t, err)
	t.Cleanup(sdk.Close)

	res, err := sdk.GetNature(ctx, GetRequest{ID: 999999})
	require.Nil(t, res)

	sdkErr, ok := err.(*SDKError)
	require.True(t, ok)

	require.Equal(t, "Not Found", sdkErr.Message)
	require.Equal(t, http.StatusNotFound, sdkErr.StatusCode)
}

func TestNature_UnexpectedError(t *testing.T) {
	ctx := context.Background()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprint(w, "Unexpected Error")
	}))
	t.Cleanup(srv.Close)

	sdk, err := NewWithOptions(Options{
		BaseURL:          srv.URL,
		Timeout:          5 * time.Second,
		CacheMaximumSize: 1 << 27,
		CacheTTL:         10 * time.Second,
	})
	require.NoError(t, err)
	t.Cleanup(sdk.Close)

	res, err := sdk.GetNature(ctx, GetRequest{ID: 999999})
	require.Nil(t, res)

	sdkErr, ok := err.(*SDKError)
	require.True(t, ok)

	require.Equal(t, "Unexpected Error", sdkErr.Message)
	require.Equal(t, http.StatusBadGateway, sdkErr.StatusCode)
}

func TestNature_ClientClosed(t *testing.T) {
	ctx := context.Background()

	sdk, err := NewWithOptions(Options{
		BaseURL:          "https://example.com",
		Timeout:          5 * time.Second,
		CacheMaximumSize: 1 << 27,
		CacheTTL:         10 * time.Second,
	})
	require.NoError(t, err)
	t.Cleanup(sdk.Close)

	sdk.Close()

	res, err := sdk.GetNature(ctx, GetRequest{ID: 1})
	require.Nil(t, res)

	sdkErr, ok := err.(*SDKError)
	require.True(t, ok)

	require.Equal(t, "sdk client closed", sdkErr.Message)
	require.Equal(t, CodeClientClosed, sdkErr.StatusCode)
}
