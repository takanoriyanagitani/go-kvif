package kvif

func ArrReduce[T, U any](a []T, init U, reducer func(state U, item T) U) U {
	var state U = init
	for _, item := range a {
		state = reducer(state, item)
	}
	return state
}
