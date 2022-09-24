package kvzip

import (
	"archive/zip"
	"context"
	"io/fs"

	ki "github.com/takanoriyanagitani/go-kvif"
	ka "github.com/takanoriyanagitani/go-kvif/pkg/fs/archive"
)

func reader2file(r *zip.Reader) func(name string) (fs.File, error) { return r.Open }

type name2BytesBuilder func(r *zip.Reader) func(name string) ([]byte, error)

func name2BytesBuilderNew(r2b ki.Reader2Bytes) name2BytesBuilder {
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

func arcGetBuilderNew(n2b name2BytesBuilder) func(*zip.Reader) ka.ArcGet {
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
