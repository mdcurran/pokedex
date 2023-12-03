package iterator

import (
	"context"
	"errors"
)

var EndOfIterator = errors.New("iterator finished")

type Iterator[T any] interface {
	// Next runs the upcoming iteration of the iterator.
	// The worker function runs, and returns a collection of T.
	// Before returning, the iterator's offset is updated so it is ready for
	// the next iteration.
	Next(ctx context.Context) ([]T, any)
}

// Iterator is an abstraction over a paginated collection.
// It is not safe to used by multiple goroutines simultaneously.
type Paginator[T any] struct {
	done   bool
	limit  uint
	offset uint

	worker func(ctx context.Context, start, end uint) ([]T, error)
}

func NewPaginator[T any](ctx context.Context, limit uint, worker func(ctx context.Context, start, end uint) ([]T, error)) *Paginator[T] {
	return &Paginator[T]{
		done:   false,
		limit:  limit,
		offset: 0,
		worker: worker,
	}
}

func (it *Paginator[T]) Next(ctx context.Context) ([]T, error) {
	if it.done {
		return nil, EndOfIterator
	}

	start := it.offset

	result, err := it.worker(ctx, start, start+it.limit)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		it.done = true
		return nil, EndOfIterator
	}

	it.offset += uint(len(result))
	return result, nil
}
