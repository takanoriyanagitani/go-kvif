package kvzip

import (
	"archive/zip"
	"context"
	"io"
	"io/fs"

	ki "github.com/takanoriyanagitani/go-kvif"
	ka "github.com/takanoriyanagitani/go-kvif/pkg/fs/archive"
)

func reader2file(r *zip.Reader) func(name string) (fs.File, error) {
	return r.Open
}

type Reader2Bytes func(r io.Reader) ([]byte, error)

var UnlimitedRead2Bytes Reader2Bytes = io.ReadAll

type Name2Bytes func(r *zip.Reader) func(name string) ([]byte, error)

func Name2BytesBuilderNew(r2b Reader2Bytes) Name2Bytes {
	return func(r *zip.Reader) func(name string) ([]byte, error) {
		return func(name string) ([]byte, error) {
			f, e := reader2file(r)(name)
			if nil != e {
				return nil, e
			}
			defer f.Close() // Reading file -> ignore close error

			return r2b(f)
		}
	}
}

var UnlimitedName2Bytes Name2Bytes = Name2BytesBuilderNew(UnlimitedRead2Bytes)

func ArcGetBuilderNew(n2b Name2Bytes) func(*zip.Reader) ka.ArcGet {
	return func(zr *zip.Reader) ka.ArcGet {
		return func(_ctx context.Context, key ka.ArcKey) (ki.Val, error) {
			var validFilename string = key.ToFilename()
			return ki.ComposeErr(
				n2b(zr),                       // string -> ([]byte, error)
				ki.ErrorFuncCreate(ki.ValNew), // []byte -> (ki.Val, error)
			)(validFilename)
		}
	}
}
