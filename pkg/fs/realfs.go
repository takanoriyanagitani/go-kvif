package kvfs

import (
	"io/fs"
	"os"

	ki "github.com/takanoriyanagitani/go-kvif"
)

func os2fs(o *os.File) fs.File { return o }

var name2file func(name string) (fs.File, error) = ki.ComposeErr(
	os.Open,                   // string -> *os.File, error
	ki.ErrorFuncCreate(os2fs), // *os.File -> fs.File, error
)

// RealFs tries to open real file.
// Its caller's responsibility to check the filename.
type RealFs func(filename string) (fs.File, error)

var OsFs RealFs = name2file
