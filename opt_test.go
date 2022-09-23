package kvif

import (
	"fmt"
	"testing"
)

func TestOpt(t *testing.T) {
	t.Parallel()

	t.Run("Opt2Err", func(t *testing.T) {
		t.Parallel()

		t.Run("empty", func(t *testing.T) {
			_, e := Opt2Err("", OptEmpty, func() error { return fmt.Errorf("Must fail") })
			t.Run("Must fail(empty -> error)", check(nil != e, true))
		})

		t.Run("non empty", func(t *testing.T) {
			i, e := Opt2Err(42, OptHasValue, func() error { panic("Must not fail") })
			t.Run("Must not fail(non empty)", check(nil == e, true))
			t.Run("Must same", check(i, 42))
		})
	})
}
