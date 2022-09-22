package kvarc

import (
	"context"

	ki "github.com/takanoriyanagitani/go-kvif"
)

func MemArcGetBuilderNew(m map[ArcKey]ki.Val) ArcGet {
	return func(_ context.Context, key ArcKey) (v ki.Val, e error) {
		v, found := m[key]
		return ki.Opt2Err(
			v,
			found,
			func() error {
				return ki.ErrNotFound
			},
		)
	}
}
