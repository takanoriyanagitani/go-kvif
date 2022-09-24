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

type name2BytesBuilder func(r *zip.Reader) func(name string) ([]byte, error)

type lstBuilder func(r *zip.Reader) func(context.Context) (ki.Iter[ka.ArcKey], error)

var lstBldNew lstBuilder = func(r *zip.Reader) func(context.Context) (ki.Iter[ka.ArcKey], error) {
	return func(_ context.Context) (keys ki.Iter[ka.ArcKey], err error) {
		var ik ki.Iter[ka.ArcKey] = ki.Compose(
			reader2files,
			ki.IterCompose(file2key),
		)(r)
		return ik, nil
	}
}

type reader2keys func(*zip.Reader) ki.Iter[ka.ArcKey]

func reader2files(r *zip.Reader) ki.Iter[*zip.File] { return ki.IterFromArr(r.File) }

func file2hdr(f *zip.File) zip.FileHeader     { return f.FileHeader }
func hdr2name(h zip.FileHeader) (name string) { return h.Name }

var file2name func(f *zip.File) (name string) = ki.Compose(
	file2hdr,
	hdr2name,
)

var name2key func(validFilename string) ka.ArcKey = ka.ArcKeyNew

var file2key func(f *zip.File) ka.ArcKey = ki.Compose(
	file2name,
	name2key,
)

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

type zipKvBuilder func(*zip.Reader) (ka.ArcKv, error)

func zipKvBuilderNew(bld ka.ArcKvBuilder) func(name2BytesBuilder) func(ka.ArcBucket) zipKvBuilder {
	return func(n2b name2BytesBuilder) func(ka.ArcBucket) zipKvBuilder {
		return func(ab ka.ArcBucket) zipKvBuilder {
			return func(zr *zip.Reader) (k ka.ArcKv, e error) {
				var g ka.ArcGet = arcGetBuilderNew(n2b)(zr)
				var c ka.ArcCls = func() error { return nil } // nothing to close
				// TODO implement lst
				l := func(context.Context) (ki.Iter[ka.ArcKey], error) { return nil, nil }
				return bld.
					WithBucket(ab).
					WithLst(l).
					WithGet(g).
					WithClose(c).
					Build()
			}
		}
	}
}

func ZipKvBuilderNew(bld ka.ArcKvBuilder) func(ki.Reader2Bytes) func(ka.ArcBucket) ka.RasKvBuilder {
	return func(r2b ki.Reader2Bytes) func(ka.ArcBucket) ka.RasKvBuilder {
		return func(ab ka.ArcBucket) ka.RasKvBuilder {
			var n2b name2BytesBuilder = name2BytesBuilderNew(r2b)
			var zkb zipKvBuilder = zipKvBuilderNew(bld)(n2b)(ab)
			return ki.ComposeErr(
				ras2rdr,
				zkb,
			)
		}
	}
}

var ZipKvBuilderDefault func(ki.Reader2Bytes) func(ka.ArcBucket) ka.RasKvBuilder = ZipKvBuilderNew(
	ka.ArcKvBuilderDefault,
)

var ZipKvBuilderDefaultUnlimited func(ka.ArcBucket) ka.RasKvBuilder = ZipKvBuilderDefault(
	ki.UnlimitedRead2Bytes,
)
