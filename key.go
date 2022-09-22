package kvif

type Key struct {
	bucket string
	raw    []byte
}

func (k Key) Bucket() string { return k.bucket }
func (k Key) Raw() []byte    { return k.raw }
