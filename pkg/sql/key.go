package kvsql

import (
	ki "github.com/takanoriyanagitani/go-kvif"
)

type SqlKey struct {
	validTablename string
	key            string
}

func SqlKeyNew(validTablename string, key string) SqlKey {
	return SqlKey{
		validTablename,
		key,
	}
}

type SqlKeyBuilder func(key ki.Key) (SqlKey, error)
type SqlKeyConverter func(SqlKey) (ki.Key, error)
