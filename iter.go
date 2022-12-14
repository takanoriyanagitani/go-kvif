package kvif

type Iter[T any] func() (o T, hasValue bool)

func IterReduce[T, U any](i Iter[T], init U, reducer func(state U, item T) U) U {
	var state U = init
	for o, hasValue := i(); hasValue; o, hasValue = i() {
		var t T = o
		state = reducer(state, t)
	}
	return state
}

func IterFromArr[T any](a []T) Iter[T] {
	var ix int = 0
	return func() (o T, hasValue bool) {
		if ix < len(a) {
			var t T = a[ix]
			ix += 1
			return t, OptHasValue
		}
		return o, OptEmpty
	}
}

func (i Iter[T]) All(f func(T) bool) bool {
	return IterReduce(i, true, func(state bool, item T) bool {
		return state && f(item)
	})
}

func IterMap[T, U any](i Iter[T], f func(T) U) Iter[U] {
	return func() (u U, hasValue bool) {
		t, hasValue := i()
		return OptMap(t, hasValue, f)
	}
}

func (i Iter[T]) Reduce(init T, reducer func(state T, item T) T) T {
	return IterReduce(i, init, reducer)
}

func (i Iter[T]) Map(f func(T) T) Iter[T] {
	return IterMap(i, f)
}

func (i Iter[T]) ToArray() []T {
	return IterReduce(i, nil, func(state []T, t T) []T {
		return append(state, t)
	})
}

func IterCompose[T, U any](f func(T) U) func(Iter[T]) Iter[U] {
	return func(it Iter[T]) Iter[U] {
		return func() (u U, hasValue bool) {
			t, ok := it()
			return OptMap(t, ok, f)
		}
	}
}

func IterReduceErr[T, U any](i Iter[T], u1st U, f func(st U, t T) (U, error)) (u U, e error) {
	u = u1st
	for o, hasValue := i(); hasValue; o, hasValue = i() {
		var t T = o
		u, e = f(u, t)
		if nil != e {
			return u, e
		}
	}
	return u, nil
}

func IterComposeErr[T, U any](f func(T) (U, error)) func(Iter[T]) (Iter[U], error) {
	return ComposeErr(
		func(it Iter[T]) ([]U, error) {
			reducer := func(state []U, t T) (au []U, e error) {
				return ComposeErr(
					f,
					func(u U) ([]U, error) { return append(au, u), nil },
				)(t)
			}
			return IterReduceErr(it, nil, reducer)
		},
		ErrorFuncCreate(IterFromArr[U]),
	)
}
