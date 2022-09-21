package kvarc

type ArcKey struct {
	validFilename string
}

func (k ArcKey) ToFilename() string { return k.validFilename }
