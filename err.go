package kvif

import (
	"errors"
)

var ErrNotFound error = errors.New("key does not exist")

func ErrorFuncCreate[T, U any](f func(T) U) func(T) (U, error) {
	return func(t T) (U, error) {
		var u U = f(t)
		return u, nil
	}
}

func ErrorFromBool[T any](ok bool, okf func() (T, error), ngf func() error) (t T, err error) {
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

func ErrorOrElse(e error, ef func() error) error {
	if nil == e {
		return ef()
	}
	return e
}

func Error1st(ef []func() error) error {
	var ie Iter[func() error] = IterFromArr(ef)
	return IterReduce(ie, nil, ErrorOrElse)
}

func ErrorFlatMap[T, U any](t T, te error, f func(T) (U, error)) (u U, e error) {
	if nil != te {
		return u, te
	}
	return f(t)
}

func ErrorMap[T, U any](t T, te error, f func(T) U) (u U, e error) {
	return ErrorFlatMap(
		t,
		te,
		ErrorFuncCreate(f),
	)
}
