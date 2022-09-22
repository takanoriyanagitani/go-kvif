package kvzip

import (
	"archive/zip"
	"context"
	"io/fs"

	ki "github.com/takanoriyanagitani/go-kvif"
	kf "github.com/takanoriyanagitani/go-kvif/pkg/fs"
	ka "github.com/takanoriyanagitani/go-kvif/pkg/fs/archive"
)

func ras2rdr(ras kf.ReaderAtSized) (*zip.Reader, error) {
	return zip.NewReader(ras.ReaderAt(), ras.Size())
}

func reader2file(r *zip.Reader) func(name string) (fs.File, error) {
	return r.Open
}

type name2Bytes func(r *zip.Reader) func(name string) ([]byte, error)

func name2BytesBuilderNew(r2b ki.Reader2Bytes) name2Bytes {
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

func arcGetBuilderNew(n2b name2Bytes) func(*zip.Reader) ka.ArcGet {
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

type zipKvBuilder func(*zip.Reader) (ka.ArcKv, error)

func zipKvBuilderNew(bld ka.ArcKvBuilder) func(name2Bytes) zipKvBuilder {
	return func(n2b name2Bytes) zipKvBuilder {
		return func(zr *zip.Reader) (k ka.ArcKv, e error) {
			var g ka.ArcGet = arcGetBuilderNew(n2b)(zr)
			var c ka.ArcCls = func() error { return nil } // nothing to close
			return bld.
				WithGet(g).
				WithClose(c).
				Build()
		}
	}
}

func ZipKvBuilderNew(bld ka.ArcKvBuilder) func(ki.Reader2Bytes) ka.RasKvBuilder {
	return func(r2b ki.Reader2Bytes) ka.RasKvBuilder {
		var n2b name2Bytes = name2BytesBuilderNew(r2b)
		var zkb zipKvBuilder = zipKvBuilderNew(bld)(n2b)
		return ki.ComposeErr(
			ras2rdr,
			zkb,
		)
	}
}

var ZipKvBuilderDefault func(ki.Reader2Bytes) ka.RasKvBuilder = ZipKvBuilderNew(
	ka.ArcKvBuilderDefault,
)

var ZipKvBuilderDefaultUnlimited ka.RasKvBuilder = ZipKvBuilderDefault(ki.UnlimitedRead2Bytes)
