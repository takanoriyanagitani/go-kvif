package kvif

const OptHasValue = true
const OptEmpty = false

func Opt2Err[T any](o T, hasValue bool, ng func() error) (t T, e error) {
	if !hasValue {
		return t, ng()
	}
	return o, nil
}

func OptMap[T, U any](o T, hasValue bool, f func(T) U) (u U, nonEmpty bool) {
	if !hasValue {
		return u, OptEmpty
	}
	u = f(o)
	return u, OptHasValue
}
