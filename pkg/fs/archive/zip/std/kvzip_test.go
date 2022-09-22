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
				})
			}
		}(ZipKvBuilderDefaultUnlimited))
	})
}
