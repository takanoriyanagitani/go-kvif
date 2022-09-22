package kvif

import (
	"context"
)

type Kv interface {
	// Get tries to get val by key.
	// If val not exists, error must be ErrNotFound
	Get(ctx context.Context, key Key) (Val, error)

	// Close closes something(optional).
	// Caller must call close.
	// Some implementation may use idempotent close.
	Close() error
}
