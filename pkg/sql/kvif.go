package kvsql

import (
	"context"

	ki "github.com/takanoriyanagitani/go-kvif"
)

type SqlGet func(context.Context, SqlKey) (ki.Val, error)

type SqlLst func(context.Context, SqlBucket) (keys ki.Iter[SqlKey], err error)

type SqlKv struct {
	get SqlGet
	lst SqlLst

	bb SqlBucketBuilder
	bk SqlKeyBuilder

	ck SqlKeyConverter
}

func (sk SqlKv) Get(ctx context.Context, key ki.Key) (ki.Val, error) {
	return ki.ComposeErr(
		sk.bk,
		func(k SqlKey) (ki.Val, error) { return sk.get(ctx, k) },
	)(key)
}

func (sk SqlKv) list(ctx context.Context, bucket SqlBucket) (keys ki.Iter[ki.Key], err error) {
	return ki.ComposeErr(
		func(b SqlBucket) (ki.Iter[SqlKey], error) { return sk.lst(ctx, b) },
		func(i ki.Iter[SqlKey]) (ki.Iter[ki.Key], error) {
			return ki.IterComposeErr(sk.ck)(i)
		},
	)(bucket)
}

func (sk SqlKv) Lst(ctx context.Context, bucket string) (keys ki.Iter[ki.Key], err error) {
	return ki.ComposeErr(
		func(b string) (SqlBucket, error) { return sk.bb(b) },
		func(b SqlBucket) (ki.Iter[ki.Key], error) { return sk.list(ctx, b) },
	)(bucket)
}
