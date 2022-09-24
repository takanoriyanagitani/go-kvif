package kvarc

import (
	"fmt"
	"io/fs"
	"regexp"
	"strings"

	ki "github.com/takanoriyanagitani/go-kvif"
)

var ArcKeySimplePattern *regexp.Regexp = regexp.MustCompile(strings.Join([]string{
	`^[a-zA-Z0-9_\.-]`,
	`[/a-zA-Z0-9_\.-]{0,}$`,
}, ""))

type ArcKey struct {
	validFilename string
}

type ArcKeyBuilder func(ki.Key) (ArcKey, error)

func (k ArcKey) ToFilename() string { return k.validFilename }

func (k ArcKey) ToKey(bucket ArcBucket) ki.Key {
	return ki.KeyNew(bucket.ToFilename(), []byte(k.validFilename))
}

func valid2akey(validFilename string) ArcKey {
	return ArcKey{validFilename}
}

func ArcKeyBuilderNew(conv ArcKeyConverter) ArcKeyBuilder {
	return ki.ComposeErr(
		conv,                           // ki.Key => (string, error)
		ki.ErrorFuncCreate(valid2akey), // string => (ArcKey, error)
	)
}

type ArcKeyConverter func(ki.Key) (validKey string, e error)

var arcKeyConverterNocheck ArcKeyConverter = func(k ki.Key) (validKey string, e error) {
	var raw []byte = k.Raw()
	return string(raw), nil
}

func ArcKeyConverterBuilderNew(validator ArcKeyValidator) ArcKeyConverter {
	return ki.ComposeErr(
		validator,
		arcKeyConverterNocheck,
	)
}

type ArcKeyValidator func(ki.Key) (ki.Key, error)

func (v ArcKeyValidator) Append(other ArcKeyValidator) ArcKeyValidator {
	return ki.ComposeErr(
		v,
		other,
	)
}

func ArcKeyValidatorRegexpBuilderNew(r *regexp.Regexp) ArcKeyValidator {
	return func(k ki.Key) (ki.Key, error) {
		var valid bool = r.Match(k.Raw())
		return ki.ErrorFromBool(
			valid,
			func() (ki.Key, error) { return k, nil },
			func() error { return fmt.Errorf("Invalid key") },
		)
	}
}

var ArcKeyValidatorRegexpSimple = ArcKeyValidatorRegexpBuilderNew(ArcKeySimplePattern)

var ArcKeyValidatorPath ArcKeyValidator = func(k ki.Key) (ki.Key, error) {
	var valid bool = fs.ValidPath(string(k.Raw()))
	return ki.ErrorFromBool(
		valid,
		func() (ki.Key, error) { return k, nil },
		func() error { return fmt.Errorf("Invalid key") },
	)
}

var ArcKeyValidatorDefault ArcKeyValidator = ArcKeyValidatorPath.Append(ArcKeyValidatorRegexpSimple)

var ArcKeyConverterDefault ArcKeyConverter = ArcKeyConverterBuilderNew(ArcKeyValidatorDefault)

var ArcKeyBuilderDefault ArcKeyBuilder = ArcKeyBuilderNew(ArcKeyConverterDefault)
