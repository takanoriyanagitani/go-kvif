package kvif

type Key struct {
	bucket string
	raw    []byte
}

func KeyNew(bucket string, id []byte) Key {
	return Key{
		bucket: bucket,
		raw:    id,
	}
}

func (k Key) Bucket() string { return k.bucket }
func (k Key) Raw() []byte    { return k.raw }
