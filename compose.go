package kvif

func ComposeErr[T, U, V any](f func(T) (U, error), g func(U) (V, error)) func(T) (V, error) {
	return func(t T) (v V, e error) {
		u, e := f(t)
		if nil != e {
			return v, e
		}

		return g(u)
	}
}

func Compose[T, U, V any](f func(T) U, g func(U) V) func(T) V {
	h := ComposeErr(
		ErrorFuncCreate(f),
		ErrorFuncCreate(g),
	)
	return func(t T) V {
		v, _ := h(t)
		return v
	}
}
