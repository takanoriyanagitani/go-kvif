package kvfs

import (
	"io"
)

type ReaderAtSized struct {
	ra io.ReaderAt
	sz int64
}

func (r ReaderAtSized) ReaderAt() io.ReaderAt { return r.ra }
func (r ReaderAtSized) Size() int64           { return r.sz }
