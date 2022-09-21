package kvfs

import (
	"bytes"
	"io/fs"
)

type MemFile struct {
	data     *bytes.Reader
	fileinfo fs.FileInfo
}

func (m MemFile) Stat() (fs.FileInfo, error) { return m.fileinfo, nil }
func (m MemFile) Read(b []byte) (int, error) { return m.data.Read(b) }
func (m MemFile) Close() error               { return nil }
