package kvfs

import (
	"bytes"
	"io"
	"io/fs"
	"testing"
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

func TestMemFs(t *testing.T) {
	t.Parallel()

	t.Run("MemFs", func(t *testing.T) {
		t.Parallel()

		t.Run("empty filename", func(t *testing.T) {
			t.Parallel()

			var mf fs.FS = MemFsNew()
			_, e := mf.Open("")
			t.Run("Must fail(empty filename)", check(nil != e, true))
		})

		t.Run("empty filesystem", func(t *testing.T) {
			t.Parallel()

			var mf fs.FS = MemFsNew()
			_, e := mf.Open("test.txt")
			t.Run("Must fail(empty filesystem)", check(nil != e, true))

			var isNoent bool = IsNotFound(e)
			t.Run("Must be noent", check(isNoent, true))
		})

		t.Run("single empty file", func(t *testing.T) {
			t.Parallel()

			var mf MemFs = MemFsNew()
			var mb MemFileBuilder = MemFileBuilderDefault
			f, e := mb.
				WithName("file.txt").
				WithReader(bytes.NewReader(nil)).
				Build()
			t.Run("mem file built", check(nil == e, true))

			e = mf.Upsert("path/to/file.txt", f)
			t.Run("mem file upserted", check(nil == e, true))

			ff, e := mf.Open("path/to/file.txt")
			t.Run("fs file got", check(nil == e, true))
			defer ff.Close()

			b, e := io.ReadAll(ff)
			t.Run("fs file read", check(nil == e, true))

			t.Run("fs file content empty", check(len(b), 0))
		})
	})
}
