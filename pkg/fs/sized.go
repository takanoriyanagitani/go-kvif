package kvfs

import (
	"bytes"
	"io"
)

type ReaderAtSized struct {
	ra io.ReaderAt
	sz int64
}

func ReaderAtSizedNew(ra io.ReaderAt, sz int64) ReaderAtSized {
	return ReaderAtSized{
		ra,
		sz,
	}
}

func ReaderAtSizedFromBytes(b []byte) ReaderAtSized {
	var rdr *bytes.Reader = bytes.NewReader(b)
	return ReaderAtSizedNew(rdr, rdr.Size())
}

func (r ReaderAtSized) ReaderAt() io.ReaderAt { return r.ra }
func (r ReaderAtSized) Size() int64           { return r.sz }
