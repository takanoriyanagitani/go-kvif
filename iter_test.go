package kvif

import (
	"testing"
)

func TestIter(t *testing.T) {
	t.Parallel()

	t.Run("All", func(t *testing.T) {
		t.Parallel()

		var ii Iter[int] = IterFromArr([]int{
			333,
			634,
		})
		var allNon0Positive bool = ii.All(func(i int) bool { return 0 < i })
		t.Run("non 0 positive", check(allNon0Positive, true))
	})

	t.Run("Map", func(t *testing.T) {
		t.Parallel()

		var ii Iter[int] = IterFromArr([]int{
			333,
			634,
		})
		var addInt = func(a, b int) int { return a + b }
		var add1 = Curry(addInt)(1)
		var i1 = ii.Map(add1)
		var tot int = i1.Reduce(0, addInt)
		t.Run("Must same", check(tot, 333+634+1+1))
	})
}
