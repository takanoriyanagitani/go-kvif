package kvif

type Val struct {
	raw []byte
}

func ValNew(raw []byte) Val {
	return Val{raw}
}

func (v Val) Raw() []byte { return v.raw }
