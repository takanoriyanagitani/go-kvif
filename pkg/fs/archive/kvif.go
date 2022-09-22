package kvarc

import (
	"context"

	ki "github.com/takanoriyanagitani/go-kvif"
)

type ArcKv interface {
	Get(ctx context.Context, key ArcKey) (ki.Val, error)
}
