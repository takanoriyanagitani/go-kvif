package kvfs

import (
	"errors"
	"io/fs"

	ki "github.com/takanoriyanagitani/go-kvif"
)

func IsNotFound(e error) bool {
	return errors.Is(e, fs.ErrNotExist) || errors.Is(e, ki.ErrNotFound)
}
