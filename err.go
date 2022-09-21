package kvif

func ErrorFuncCreate[T, U any](f func(T) U) func(T) (U, error) {
	return func(t T) (U, error) {
		var u U = f(t)
		return u, nil
	}
}

func Bool2Err[T any](ok bool, okf func() (T, error), ngf func() error) (t T, err error) {
	if ok {
		return okf()
	}
	return t, ngf()
}

func ErrorTryForEach[T any](t T, e error, f func(T) error) error {
	if nil == e {
		return f(t)
	}
	return e
}
