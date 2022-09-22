package kvif

import (
	"io"
)

type Reader2Bytes func(r io.Reader) ([]byte, error)

var UnlimitedRead2Bytes Reader2Bytes = io.ReadAll
