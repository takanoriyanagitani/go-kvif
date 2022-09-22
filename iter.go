package kvif

type Iter[T any] func() (o T, empty bool)

func IterReduce[T, U any](i Iter[T], init U, reducer func(state U, item T) U) U {
	var state U = init
	for o, empty := i(); !empty; o, empty = i() {
		var t T = o
		state = reducer(state, t)
	}
	return state
}

func IterFromArr[T any](a []T) Iter[T] {
	var ix int = 0
	return func() (o T, empty bool) {
		if ix < len(a) {
			var t T = a[ix]
			ix += 1
			return t, false
		}
		return o, true
	}
}

func (i Iter[T]) All(f func(T) bool) bool {
	return IterReduce(i, true, func(state bool, item T) bool {
		return state && f(item)
	})
}
