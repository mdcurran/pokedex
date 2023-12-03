package pokedex

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetRequest(t *testing.T) {
	r := GetRequest{ID: 1}
	resource, err := r.GetResource()
	require.NoError(t, err)
	require.Equal(t, "1", resource)

	r = GetRequest{Name: "foobar"}
	resource, err = r.GetResource()
	require.NoError(t, err)
	require.Equal(t, "foobar", resource)

	r = GetRequest{}
	resource, err = r.GetResource()
	require.Empty(t, resource)
	gErr := err.(*SDKError)
	require.Equal(t, ErrMissingResources.Error(), gErr.Message)
	require.Equal(t, CodeInvalidArgs, gErr.StatusCode)

	r = GetRequest{ID: 1, Name: "foobar"}
	resource, err = r.GetResource()
	require.Empty(t, resource)
	gErr = err.(*SDKError)
	require.Equal(t, ErrMultipleResources.Error(), gErr.Message)
	require.Equal(t, CodeInvalidArgs, gErr.StatusCode)
}
