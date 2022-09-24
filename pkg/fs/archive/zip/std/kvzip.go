package kvzip

import (
	"archive/zip"

	ki "github.com/takanoriyanagitani/go-kvif"
	kf "github.com/takanoriyanagitani/go-kvif/pkg/fs"
	ka "github.com/takanoriyanagitani/go-kvif/pkg/fs/archive"
)

func ras2rdr(ras kf.ReaderAtSized) (*zip.Reader, error) {
	return zip.NewReader(ras.ReaderAt(), ras.Size())
}

type zipKvBuilder func(*zip.Reader) (ka.ArcKv, error)

func zipKvBuilderNew(bld ka.ArcKvBuilder) func(name2BytesBuilder) func(ka.ArcBucket) zipKvBuilder {
	return func(n2b name2BytesBuilder) func(ka.ArcBucket) zipKvBuilder {
		return func(ab ka.ArcBucket) zipKvBuilder {
			return func(zr *zip.Reader) (k ka.ArcKv, e error) {
				var g ka.ArcGet = arcGetBuilderNew(n2b)(zr)
				var c ka.ArcCls = func() error { return nil } // nothing to close
				var l ka.ArcLst = lstBldNew(zr)
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
