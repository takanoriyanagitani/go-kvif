package kvif

type Val struct {
	raw []byte
}

func ValNew(raw []byte) Val {
	return Val{raw}
}
