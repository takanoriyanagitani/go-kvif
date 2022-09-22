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

func TestAll(t *testing.T) {
	t.Parallel()

	t.Run("ZipKvBuilderDefaultUnlimited", func(t *testing.T) {
		t.Parallel()

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
			}
		}(ZipKvBuilderDefaultUnlimited))
	})
}
