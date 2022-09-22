package kvfs

import (
	"io/fs"
	"time"
)

var MemInfoDefaultMode fs.FileMode = 0644

type MemInfo struct {
	name     string
	size     int64
	mode     fs.FileMode
	modified time.Time
}

func (i MemInfo) Name() string       { return i.name }
func (i MemInfo) Size() int64        { return i.size }
func (i MemInfo) Mode() fs.FileMode  { return i.mode }
func (i MemInfo) ModTime() time.Time { return i.modified }
func (i MemInfo) IsDir() bool        { return false }
func (i MemInfo) Sys() any           { return nil }

type MemInfoBuilder func(modified time.Time) func(name string) func(size int64) MemInfo

func MemInfoBuilderNew(mode fs.FileMode) MemInfoBuilder {
	return func(modified time.Time) func(string) func(int64) MemInfo {
		return func(name string) func(int64) MemInfo {
			return func(size int64) MemInfo {
				return MemInfo{
					name,
					size,
					mode,
					modified,
				}
			}
		}
	}
}

var MemInfoBuilderDefault MemInfoBuilder = MemInfoBuilderNew(MemInfoDefaultMode)
