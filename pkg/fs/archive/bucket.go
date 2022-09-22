package kvarc

import (
	"io/fs"

	kf "github.com/takanoriyanagitani/go-kvif/pkg/fs"
)

type ArcBucket struct {
	validFilename string
}

func (a ArcBucket) Open(f fs.FS) (fs.File, error) {
	return f.Open(a.validFilename)
}

func (a ArcBucket) ToMemFile(f fs.FS, b kf.MemFileBuilder) (m kf.MemFile, e error) {
	file, e := a.Open(f)
	if nil != e {
		return m, e
	}
	defer file.Close()

	s, e := file.Stat()
	if nil != e {
		return m, e
	}

	return b.WithModified(s.ModTime()).
		WithName(s.Name()).
		WithReader(file).
		Build()
}
