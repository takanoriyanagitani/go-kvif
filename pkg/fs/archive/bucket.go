package kvarc

import (
	"fmt"
	"io/fs"
	"regexp"
	"strings"

	ki "github.com/takanoriyanagitani/go-kvif"
	kf "github.com/takanoriyanagitani/go-kvif/pkg/fs"
)

var ArcBucketSimplePattern *regexp.Regexp = regexp.MustCompile(strings.Join([]string{
	`^[a-zA-Z0-9_\.-]`,
	`[/a-zA-Z0-9_\.-]{0,}$`,
}, ""))

// ArcBucket is archive identifier.
type ArcBucket struct {
	validFilename string
}

type ArcBucketBuilder func(ki.Key) (ArcBucket, error)

type ArcBucketConverter func(ki.Key) (validFilename string, e error)

type ArcBucketValidator func(ki.Key) (ki.Key, error)

func valid2bucket(validFilename string) ArcBucket {
	return ArcBucket{validFilename}
}

func ArcBucketBuilderNew(conv ArcBucketConverter) ArcBucketBuilder {
	return ki.ComposeErr(
		conv,
		ki.ErrorFuncCreate(valid2bucket),
	)
}

func (v ArcBucketValidator) Append(other ArcBucketValidator) ArcBucketValidator {
	return ki.ComposeErr(
		v,
		other,
	)
}

var arcBucketConvNoCheck ArcBucketConverter = func(k ki.Key) (validFilename string, e error) {
	return k.Bucket(), nil
}

func ArcBucketConvBuilderNew(v ArcBucketValidator) ArcBucketConverter {
	return ki.ComposeErr(
		v,
		arcBucketConvNoCheck,
	)
}

func ArcBucketValidatorRegexpBuilderNew(r *regexp.Regexp) ArcBucketValidator {
	return func(k ki.Key) (ki.Key, error) {
		return ki.ErrorFromBool(
			r.MatchString(k.Bucket()),
			func() (ki.Key, error) { return k, nil },
			func() error { return fmt.Errorf("Invalid bucket: %s", k.Bucket()) },
		)
	}
}

var ArcBucketValidatorRegexpSimple = ArcBucketValidatorRegexpBuilderNew(ArcBucketSimplePattern)

var ArcBucketValidatorPath ArcBucketValidator = func(k ki.Key) (ki.Key, error) {
	return ki.ErrorFromBool(
		fs.ValidPath(k.Bucket()),
		func() (ki.Key, error) { return k, nil },
		func() error { return fmt.Errorf("Invalid bucket: %s", k.Bucket()) },
	)
}

var ArcBucketValidatorDefault ArcBucketValidator = ArcBucketValidatorPath.Append(
	ArcBucketValidatorRegexpSimple,
)

var ArcBucketConverterDefault ArcBucketConverter = ArcBucketConvBuilderNew(ArcBucketValidatorDefault)

var ArcBucketBuilderDefault ArcBucketBuilder = ArcBucketBuilderNew(ArcBucketConverterDefault)

// Open tries to open archive file from virtual FS.
func (a ArcBucket) Open(f fs.FS) (fs.File, error) {
	return f.Open(a.validFilename)
}

// ToMemFile converts archive file into in-mem file.
func (a ArcBucket) ToMemFile(f fs.FS, b kf.MemFileBuilder) (m kf.MemFile, e error) {
	file, e := a.Open(f)
	if nil != e {
		return m, e
	}
	defer file.Close()

	s, e := file.Stat()
	if nil != e {
		return m, e
	}

	return b.WithModified(s.ModTime()).
		WithName(s.Name()).
		WithReader(file).
		Build()
}
