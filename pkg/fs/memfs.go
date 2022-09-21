package kvfs

import (
	"io/fs"

	ki "github.com/takanoriyanagitani/go-kvif"
)

func pathErrNew(Op string, Path string, Err error) *fs.PathError {
	return &fs.PathError{
		Op,
		Path,
		Err,
	}
}

func openErrNew(Err error, Path string) *fs.PathError { return pathErrNew("open", Path, Err) }

var openErr func(err error) func(path string) *fs.PathError = ki.Curry(openErrNew)

var invalidErr func(path string) *fs.PathError = openErr(fs.ErrInvalid)
var noentryErr func(path string) *fs.PathError = openErr(fs.ErrNotExist)

func getValidName(unchecked string) (validFilename string, e error) {
	return ki.Bool2Err(
		fs.ValidPath(validFilename),
		func() (string, error) { return unchecked, nil },
		func() error { return invalidErr(unchecked) },
	)
}

type MemFs struct {
	filemap map[string]fs.File
}

func (m MemFs) open(validFilename string) (fs.File, error) {
	f, ok := m.filemap[validFilename]
	return ki.Opt2Err(f, ok, func() error { return noentryErr(validFilename) })
}

func (m MemFs) Open(unchecked string) (fs.File, error) {
	var s2f func(string) (fs.File, error) = ki.ComposeErr(
		getValidName, // string -> (string,  error)
		m.open,       // string -> (fs.File, error)
	)
	return s2f(unchecked)
}

func (m MemFs) upsert(validFilename string, mf MemFile) {
	m.filemap[validFilename] = mf
}

func (m MemFs) Upsert(unchecked string, mf MemFile) error {
	valid, err := getValidName(unchecked)
	return ki.ErrorTryForEach(
		valid,
		err,
		func(validFilename string) error {
			m.upsert(validFilename, mf)
			return nil
		},
	)
}
