package iterator

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPaginator(t *testing.T) {
	var (
		ctx  = context.Background()
		ints = ints(55)
	)

	it := NewPaginator(ctx, 10, func(ctx context.Context, start, end uint) ([]int, error) {
		if start >= uint(len(ints)) {
			return []int{}, nil
		}

		if end >= uint(len(ints)) {
			end = uint(len(ints))
		}

		return ints[start:end], nil
	})

	for i := 0; i < 5; i++ {
		ns, err := it.Next(ctx)
		require.NoError(t, err)
		require.Len(t, ns, 10)
	}

	// The last page of this iterator should return an abridged number of
	// integers as we've reached the final page.
	ns, err := it.Next(ctx)
	require.NoError(t, err)
	require.Equal(t, []int{51, 52, 53, 54, 55}, ns)

	// The next iteration should return an EndOfIteration error. Any subsequent
	// calls to Next() should be idempotent and return the same error.
	ns, err = it.Next(ctx)
	require.ErrorIs(t, err, EndOfIterator)
	require.Empty(t, ns)

	ns, err = it.Next(ctx)
	require.ErrorIs(t, err, EndOfIterator)
	require.Empty(t, ns)
}

func TestPaginatorError(t *testing.T) {
	var (
		ctx         = context.Background()
		ints        = ints(100)
		inducedErr  = errors.New("something broke")
		shouldError = true
	)

	it := NewPaginator(ctx, 10, func(ctx context.Context, start, end uint) ([]int, error) {
		// Induce a transient failure. This is akin to some sort of
		// 5xx response from an external API. After inducing the error we
		// can "fix" it.
		if shouldError {
			return nil, inducedErr
		}

		if start >= uint(len(ints)) {
			return []int{}, nil
		}

		if end >= uint(len(ints)) {
			end = uint(len(ints))
		}

		return ints[start:end], nil
	})

	ns, err := it.Next(ctx)
	require.ErrorIs(t, err, inducedErr)
	require.Empty(t, ns)

	// If the transient error is addressed, the iterator can continue.
	shouldError = false

	for i := 0; i < 10; i++ {
		ns, err := it.Next(ctx)
		require.NoError(t, err)
		require.Len(t, ns, 10)
	}

	// Get to the end of this iterator.
	ns, err = it.Next(ctx)
	require.ErrorIs(t, err, EndOfIterator)
	require.Empty(t, ns)
}

func ints(count int) []int {
	var ns []int
	for i := 1; i <= count; i++ {
		ns = append(ns, i)
	}
	return ns
}
