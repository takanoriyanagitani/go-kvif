package kvif

import (
	"context"
)

type Kv interface {
	// Get tries to get val by key.
	// If val not exists, error must be ErrNotFound
	Get(ctx context.Context, key Key) (Val, error)

	// Lst tries to get keys in a bucket.
	Lst(ctx context.Context, bucket string) (keys Iter[Key], err error)

	// Close closes something(optional).
	// Caller must call close.
	// Some implementation may use idempotent close.
	Close() error
}
