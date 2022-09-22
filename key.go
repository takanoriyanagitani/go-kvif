package kvif

type Key struct {
	raw []byte
}

func (k Key) Raw() []byte { return k.raw }
