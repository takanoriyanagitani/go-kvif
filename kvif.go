package kvif

import (
	"context"
	"io"
)

// Getter gets Val by Key.
// Error for non-existent key must be ErrNotFound.
type Getter interface {
	Get(ctx context.Context, key Key) (Val, error)
}

// Lister gets keys by bucket.
type Lister interface {
	Lst(ctx context.Context, bucket string) (keys Iter[Key], err error)
}

type Kv interface {
	Getter
	Lister
	io.Closer
}
