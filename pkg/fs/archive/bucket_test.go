package kvarc

import (
	"bytes"
	"io/fs"
	"testing"

	ki "github.com/takanoriyanagitani/go-kvif"
	kf "github.com/takanoriyanagitani/go-kvif/pkg/fs"
)

func TestBucket(t *testing.T) {
	t.Parallel()

	t.Run("ArcBucketBuilderDefault", func(t *testing.T) {
		t.Parallel()

		t.Run("builder got", func(bldr ArcBucketBuilder) func(*testing.T) {
			return func(t *testing.T) {
				t.Parallel()

				t.Run("empty", func(t *testing.T) {
					t.Parallel()

					_, e := bldr(ki.KeyNew("", nil))
					t.Run("Must fail(empty)", check(nil != e, true))
				})

				t.Run("valid path, invalid bucket", func(t *testing.T) {
					t.Parallel()

					_, e := bldr(ki.KeyNew("#invalid/path/to/archive.zip", nil))
					t.Run("Must fail(invalid bucket)", check(nil != e, true))
				})

				t.Run("valid bucket", func(t *testing.T) {
					t.Parallel()

					var s string = "path/to/archive.zip"
					b, e := bldr(ki.KeyNew(s, nil))
					t.Run("Must not fail(valid bucket)", check(nil == e, true))

					var mf fs.FS = kf.MemFsNew()

					_, e = b.Open(mf)
					t.Run("Must fail(empty filesystem)", check(nil != e, true))

					var isNoent bool = kf.IsNotFound(e)
					t.Run("Must be noent", check(isNoent, true))
				})

				t.Run("ToMemFile", func(t *testing.T) {
					t.Parallel()

					t.Run("noent", func(t *testing.T) {
						var s string = "path/to/archive.zip"
						b, e := bldr(ki.KeyNew(s, nil))
						t.Run("Must not fail(valid bucket)", check(nil == e, true))

						var mf fs.FS = kf.MemFsNew()
						var mb kf.MemFileBuilder = kf.MemFileBuilderDefault

						_, e = b.ToMemFile(mf, mb)
						t.Run("Must fail(no ent)", check(nil != e, true))
					})

					t.Run("empty file(invalid archive)", func(t *testing.T) {
						var s string = "path/to/archive.zip"
						b, e := bldr(ki.KeyNew(s, nil))
						t.Run("Must not fail(valid bucket)", check(nil == e, true))

						var mf kf.MemFs = kf.MemFsNew()
						var mb kf.MemFileBuilder = kf.MemFileBuilderDefault

						vfile, e := mb.
							WithName("archive.zip").
							WithReader(bytes.NewReader(nil)).
							Build()
						t.Run("Must not fail(create vfile)", check(nil == e, true))

						e = mf.Upsert(s, vfile)
						t.Run("Must not fail(upsert vfile)", check(nil == e, true))

						got, e := b.ToMemFile(mf, kf.MemFileBuilderDefault)
						t.Run("Must not fail(get invalid archive)", check(nil == e, true))
						defer got.Close() // optional(MemFile.Close is nop)
					})

				})
			}
		}(ArcBucketBuilderDefault))
	})
}
