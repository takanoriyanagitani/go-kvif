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
	bkt ArcBucket
}

func (a ArcKv) Get(ctx context.Context, key ki.Key) (v ki.Val, e error) {
	var f func(context.Context) func(ArcKey) (ki.Val, error) = ki.CurryCtx(a.get)
	var g func(ArcKey) (ki.Val, error) = f(ctx)
	return ki.ComposeErr(
		a.bld,
		g,
	)(key)
}

func (a ArcKv) ArchiveName() string { return a.bkt.ToFilename() }

func (a ArcKv) convertKey(ak ArcKey) ki.Key { return ak.ToKey(a.bkt) }

func (a ArcKv) sameBucket(bucket string) bool {
	return a.bkt.Equals(bucket)
}

func (a ArcKv) checkBucket(bucket string) (ArcKv, error) {
	return ki.ErrorFromBool(
		a.sameBucket(bucket),
		func() (ArcKv, error) { return a, nil },
		func() error { return fmt.Errorf("Invalid bucket: %s", bucket) },
	)
}

func (a ArcKv) lstChecked(ctx context.Context) (keys ki.Iter[ki.Key], err error) {
	return ki.ComposeErr(
		func(kv ArcKv) (ki.Iter[ArcKey], error) { return kv.lst(ctx) },
		ki.ErrorFuncCreate(ki.IterCompose(a.convertKey)),
	)(a)
}

func (a ArcKv) Lst(ctx context.Context, unchecked string) (keys ki.Iter[ki.Key], err error) {
	return ki.ComposeErr(
		a.checkBucket,
		func(kv ArcKv) (ki.Iter[ki.Key], error) { return kv.lstChecked(ctx) },
	)(unchecked)
}

func (a ArcKv) Close() error { return a.cls() }

type ArcKvBuilder struct {
	ArcGet
	ArcKeyBuilder
	ArcCls
	ArcLst
	ArcBucket
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

func (b ArcKvBuilder) WithLst(l ArcLst) ArcKvBuilder {
	b.ArcLst = l
	return b
}

func (b ArcKvBuilder) WithBucket(ab ArcBucket) ArcKvBuilder {
	b.ArcBucket = ab
	return b
}

func (b ArcKvBuilder) Build() (a ArcKv, e error) {
	var valid bool = ki.IterFromArr([]bool{
		nil != b.ArcGet,
		nil != b.ArcKeyBuilder,
		nil != b.ArcCls,
		nil != b.ArcLst,
		b.ArcBucket.hasValue(),
	}).All(ki.Identity[bool])

	return ki.ErrorFromBool(
		valid,
		func() (ArcKv, error) {
			return ArcKv{
				get: b.ArcGet,
				bld: b.ArcKeyBuilder,
				cls: b.ArcCls,
				lst: b.ArcLst,
				bkt: b.ArcBucket,
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
