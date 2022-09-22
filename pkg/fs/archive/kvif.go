package kvarc

import (
	"context"

	ki "github.com/takanoriyanagitani/go-kvif"
)

type ArcGet func(ctx context.Context, key ArcKey) (ki.Val, error)

type ArcCls func() error

type ArcKv struct {
	get ArcGet
	bld ArcKeyBuilder
	cls ArcCls
}

func (a ArcKv) Get(ctx context.Context, key ki.Key) (v ki.Val, e error) {
	var f func(context.Context) func(ArcKey) (ki.Val, error) = ki.CurryCtx(a.get)
	var g func(ArcKey) (ki.Val, error) = f(ctx)
	return ki.ComposeErr(
		a.bld,
		g,
	)(key)
}

func (a ArcKv) Close() error { return a.cls() }

type ArcKvBuilder struct {
	ArcGet
	ArcKeyBuilder
	ArcCls
}

func (b ArcKvBuilder) Default() ArcKvBuilder {
	b.ArcKeyBuilder = ArcKeyBuilderDefault
	return b
}

func (b ArcKvBuilder) WithGet(g ArcGet) ArcKvBuilder {
	b.ArcGet = g
	return b
}

func (b ArcKvBuilder) WithKeyBuilder(k ArcKeyBuilder) ArcKvBuilder {
	b.ArcKeyBuilder = k
	return b
}

func (b ArcKvBuilder) WithClose(c ArcCls) ArcKvBuilder {
	b.ArcCls = c
	return b
}

func (b ArcKvBuilder) Build() (a ArcKv, e error) {
	var valid bool = ki.IterFromArr([]bool{
		nil != b.ArcGet,
		nil != b.ArcKeyBuilder,
		nil != b.ArcCls,
	}).All(ki.Identity[bool])

	return ki.ErrorFromBool(
		valid,
		func() (ArcKv, error) {
			return ArcKv{
				get: b.ArcGet,
				bld: b.ArcKeyBuilder,
				cls: b.ArcCls,
			}, nil
		},
		func() error {
			return nil
		},
	)
}
