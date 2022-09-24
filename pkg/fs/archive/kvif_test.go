package kvarc

import (
	"bytes"
	"context"
	"testing"

	ki "github.com/takanoriyanagitani/go-kvif"
	kf "github.com/takanoriyanagitani/go-kvif/pkg/fs"
)

func checkBuilder[T any](comp func(a, b T) (same bool)) func(got, expected T) func(*testing.T) {
	return func(got, expected T) func(*testing.T) {
		return func(t *testing.T) {
			var same bool = comp(got, expected)
			if !same {
				t.Errorf("Unexpected value got.\n")
				t.Errorf("Expected: %v\n", expected)
				t.Fatalf("Got:      %v\n", got)
			}
		}
	}
}

func check[T comparable](got, expected T) func(*testing.T) {
	return checkBuilder(
		func(a, b T) (same bool) { return a == b },
	)(got, expected)
}

var checkBytes func(got, expected []byte) func(*testing.T) = checkBuilder(
	func(a, b []byte) (same bool) { return 0 == bytes.Compare(a, b) },
)

func TestAll(t *testing.T) {
	t.Parallel()

	t.Run("ArcKvBuilderDefault", func(t *testing.T) {
		t.Parallel()

		t.Run("builder got", func(akb ArcKvBuilder) func(*testing.T) {
			return func(t *testing.T) {
				t.Parallel()

				t.Run("MemArcGetBuilderNew", func(t *testing.T) {
					t.Parallel()

					t.Run("empty", func(t *testing.T) {
						t.Parallel()

						var g ArcGet = MemArcGetBuilderNew(make(map[ArcKey]ki.Val))
						var l ArcLst = func(context.Context) (ki.Iter[ArcKey], error) {
							return nil, nil
						}
						b, e := ArcBucketBuilderDefault(ki.KeyNew("archive.zip", []byte("hw")))
						t.Run("bucket built", check(nil == e, true))

						ak, e := akb.WithGet(g).
							WithLst(l).
							WithBucket(b).
							WithClose(func() error { return nil }).
							Build()
						t.Run("Must not fail(empty)", check(nil == e, true))
						defer ak.Close()

						_, e = ak.Get(context.Background(), ki.KeyNew("", []byte("hw")))
						t.Run("Must fail(empty)", check(nil != e, true))

						var isNoEnt bool = kf.IsNotFound(e)
						t.Run("Must be noent", check(isNoEnt, true))
					})

					t.Run("getter missing", func(t *testing.T) {
						t.Parallel()

						_, e := akb.WithClose(func() error { return nil }).
							WithKeyBuilder(nil).
							Build()
						t.Run("Must fail(getter missing)", check(nil != e, true))
					})
				})
			}
		}(ArcKvBuilderDefault))
	})
}
