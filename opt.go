package kvif

func Opt2Err[T any](o T, ok bool, ng func() error) (t T, e error) {
	if ok {
		return o, nil
	}
	return t, ng()
}
