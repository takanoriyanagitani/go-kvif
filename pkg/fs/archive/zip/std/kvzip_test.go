package kvzip

import (
	"archive/zip"
	"bytes"
	"context"
	"testing"

	ki "github.com/takanoriyanagitani/go-kvif"
	kf "github.com/takanoriyanagitani/go-kvif/pkg/fs"
	ka "github.com/takanoriyanagitani/go-kvif/pkg/fs/archive"
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

	t.Run("ZipKvBuilderDefaultUnlimited", func(t *testing.T) {
		t.Parallel()

		ab, e := ka.ArcBucketBuilderDefault(ki.KeyNew("archive.zip", []byte("empty.txt")))
		t.Run("bucket got", check(nil == e, true))

		t.Run("ras2kv got", func(rkb ka.RasKvBuilder) func(*testing.T) {
			return func(t *testing.T) {
				t.Parallel()

				t.Run("invalid archive file", func(t *testing.T) {
					t.Parallel()

					var ras kf.ReaderAtSized = kf.ReaderAtSizedFromBytes(nil)

					_, e := rkb(ras)
					t.Run("Must fail(invalid zip archive)", check(nil != e, true))
				})

				t.Run("empty archive file", func(t *testing.T) {
					t.Parallel()

					var zipbytes bytes.Buffer
					var zw *zip.Writer = zip.NewWriter(&zipbytes)

					var e error = zw.Close()
					t.Run("empty zip created", check(nil == e, true))

					var ras kf.ReaderAtSized = kf.ReaderAtSizedFromBytes(zipbytes.Bytes())

					a, e := rkb(ras)
					t.Run("Must not fail(valid empty zip)", check(nil == e, true))
					defer a.Close()

					_, e = a.Get(context.Background(), ki.KeyNew("", []byte("hw.txt")))
					t.Run("Must fail(empty zip)", check(nil != e, true))

					var isNoent bool = kf.IsNotFound(e)
					t.Run("Must be noent", check(isNoent, true))
				})

				t.Run("single empty archive item", func(t *testing.T) {
					t.Parallel()

					var zipbytes bytes.Buffer
					var zw *zip.Writer = zip.NewWriter(&zipbytes)

					zh := zip.FileHeader{
						Name:   "empty.txt",
						Method: zip.Store,
					}

					_, e := zw.CreateHeader(&zh)
					t.Run("zip item created", check(nil == e, true))

					e = zw.Close()
					t.Run("zip created", check(nil == e, true))

					var ras kf.ReaderAtSized = kf.ReaderAtSizedFromBytes(zipbytes.Bytes())

					a, e := rkb(ras)
					t.Run("Must not fail(valid empty zip)", check(nil == e, true))
					defer a.Close()

					_, e = a.Get(context.Background(), ki.KeyNew("", []byte("hw.txt")))
					t.Run("Must fail(no such item)", check(nil != e, true))

					var isNoent bool = kf.IsNotFound(e)
					t.Run("Must be noent", check(isNoent, true))

					v, e := a.Get(context.Background(), ki.KeyNew("", []byte("empty.txt")))
					t.Run("Must not fail", check(nil == e, true))

					t.Run("Must be empty", check(0, len(v.Raw())))
				})

				t.Run("many non empty items", func(t *testing.T) {
					t.Parallel()

					var zipbytes bytes.Buffer
					var zw *zip.Writer = zip.NewWriter(&zipbytes)

					createItem := func(name string, content []byte) {
						zh := zip.FileHeader{
							Name:   name,
							Method: zip.Store,
						}

						w, e := zw.CreateHeader(&zh)
						t.Run("zip item created", check(nil == e, true))
						_, e = w.Write(content)
						t.Run("zip item wrote", check(nil == e, true))
					}

					createItem("test1.txt", []byte("hw"))
					createItem("test2.txt", []byte("hh"))

					e := zw.Close()
					t.Run("zip created", check(nil == e, true))

					var ras kf.ReaderAtSized = kf.ReaderAtSizedFromBytes(zipbytes.Bytes())

					akv, e := rkb(ras)
					t.Run("Must not fail(valid empty zip)", check(nil == e, true))
					defer akv.Close()

					_, e = akv.Get(context.Background(), ki.KeyNew("", []byte("hw.txt")))
					t.Run("Must fail(no such item)", check(nil != e, true))

					var isNoent bool = kf.IsNotFound(e)
					t.Run("Must be noent", check(isNoent, true))

					chk := func(name string, expected []byte) func(*testing.T) {
						return func(t *testing.T) {
							v, e := akv.Get(context.Background(), ki.KeyNew("", []byte(name)))
							t.Run("Must not fail", check(nil == e, true))

							t.Run("Must be same", checkBytes(v.Raw(), expected))
						}
					}

					t.Run("test1", chk("test1.txt", []byte("hw")))
					t.Run("test2", chk("test2.txt", []byte("hh")))

					t.Run("archive name", check(akv.ArchiveName(), "archive.zip"))

					t.Run("Lst", func(t *testing.T) {
						keys, e := akv.Lst(context.Background())
						t.Run("Must not fail(get keys)", check(nil == e, true))

						var ka []ki.Key = keys.ToArray()
						t.Run("2 files", check(len(ka), 2))

						var k1 ki.Key = ka[0]
						var k2 ki.Key = ka[1]
						t.Run("bucket", check(k1.Bucket(), "archive.zip"))
						t.Run("bucket", check(k2.Bucket(), "archive.zip"))

						t.Run("test1", checkBytes(k1.Raw(), []byte("test1.txt")))
						t.Run("test2", checkBytes(k2.Raw(), []byte("test2.txt")))
					})
				})
			}
		}(ZipKvBuilderDefaultUnlimited(ab)))
	})
}
