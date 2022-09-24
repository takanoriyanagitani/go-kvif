package kvzip

import (
	"archive/zip"
	"context"

	ki "github.com/takanoriyanagitani/go-kvif"
	ka "github.com/takanoriyanagitani/go-kvif/pkg/fs/archive"
)

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
