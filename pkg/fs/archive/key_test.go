package kvarc

import (
	"testing"

	ki "github.com/takanoriyanagitani/go-kvif"
)

func TestKey(t *testing.T) {
	t.Parallel()

	t.Run("ArcKeyBuilderDefault", func(t *testing.T) {
		t.Parallel()

		t.Run("ArcKeyBuilder got", func(akb ArcKeyBuilder) func(*testing.T) {
			return func(t *testing.T) {
				t.Parallel()

				t.Run("empty", func(t *testing.T) {
					_, e := akb(ki.KeyNew("", nil))
					t.Run("Must fail(empty)", check(nil != e, true))
				})

				t.Run("valid path, invalid key", func(t *testing.T) {
					_, e := akb(ki.KeyNew("", []byte("#/invalid/path/to/something.txt")))
					t.Run("Must fail(invalid key)", check(nil != e, true))
				})

				t.Run("valid key", func(t *testing.T) {
					var s string = "path/to/file/inside/archive.txt"
					a, e := akb(ki.KeyNew("", []byte(s)))
					t.Run("Must not fail(valid key)", check(nil == e, true))

					t.Run("Must be same", check(a.ToFilename(), s))
				})
			}
		}(ArcKeyBuilderDefault))
	})
}
