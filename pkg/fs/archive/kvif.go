package kvarc

import (
	"context"
	"fmt"
	"io/fs"

	ki "github.com/takanoriyanagitani/go-kvif"
	kf "github.com/takanoriyanagitani/go-kvif/pkg/fs"
)

type ArcGet func(ctx context.Context, key ArcKey) (ki.Val, error)

type ArcLst func(ctx context.Context) (keys ki.Iter[ArcKey], err error)

type ArcCls func() error

type ArcKv struct {
	get ArcGet
	bld ArcKeyBuilder
	cls ArcCls
	lst ArcLst
}

func (a ArcKv) Get(ctx context.Context, key ki.Key) (v ki.Val, e error) {
	var f func(context.Context) func(ArcKey) (ki.Val, error) = ki.CurryCtx(a.get)
	var g func(ArcKey) (ki.Val, error) = f(ctx)
	return ki.ComposeErr(
		a.bld,
		g,
	)(key)
}

func (a ArcKv) Lst(ctx context.Context) (keys ki.Iter[ki.Key], err error) {
	return nil, nil
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
			return fmt.Errorf("Invalid builder")
		},
	)
}

var ArcKvBuilderDefault ArcKvBuilder = ArcKvBuilder{}.Default()

type RasKvBuilder func(kf.ReaderAtSized) (ArcKv, error)

type ArcKvFactory struct {
	vfs fs.FS
	mfb kf.MemFileBuilder
}

func (f ArcKvFactory) NewKvBuilder(kb RasKvBuilder) func(ArcBucket) (ArcKv, error) {
	return ki.ComposeErr(
		func(ab ArcBucket) (kf.MemFile, error) { return ab.ToMemFile(f.vfs, f.mfb) },
		func(mf kf.MemFile) (ArcKv, error) { return kb(mf.ToSized()) },
	)
}
