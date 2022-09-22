package kvfs

import (
	"errors"
	"io/fs"
)

func IsNotFound(e error) bool {
	return errors.Is(e, fs.ErrNotExist)
}
