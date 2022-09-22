package kvarc

import (
	"context"

	ki "github.com/takanoriyanagitani/go-kvif"
)

type ArcGet func(ctx context.Context, key ArcKey) (ki.Val, error)

type ArcKv struct {
	get ArcGet
	bld ArcKeyBuilder
}

func (a ArcKv) Get(ctx context.Context, key ki.Key) (v ki.Val, e error) {
	var f func(context.Context) func(ArcKey) (ki.Val, error) = ki.CurryCtx(a.get)
	var g func(ArcKey) (ki.Val, error) = f(ctx)
	return ki.ComposeErr(
		a.bld,
		g,
	)(key)
}
