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

func TestGetPokemon(t *testing.T) {
	ctx := context.Background()

	fixture, err := json.Marshal(faker.NewFaker().GeneratePokemon())
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

	res, err := sdk.GetPokemon(ctx, GetRequest{ID: 1})
	require.NoError(t, err)

	serialised, err := json.Marshal(res.Pokemon)
	require.NoError(t, err)
	require.JSONEq(t, string(fixture), string(serialised))

	require.Len(t, res.Pokemon.Abilities, 1)
}
