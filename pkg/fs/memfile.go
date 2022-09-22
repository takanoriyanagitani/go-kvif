package kvfs

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"time"

	ki "github.com/takanoriyanagitani/go-kvif"
)

var TimeEpoch time.Time = time.Unix(0, 0)

type MemFile struct {
	data     *bytes.Reader
	fileinfo fs.FileInfo
}

func memFileNew(data *bytes.Reader, fileinfo fs.FileInfo) MemFile {
	return MemFile{
		data,
		fileinfo,
	}
}

func (m MemFile) Stat() (fs.FileInfo, error)              { return m.fileinfo, nil }
func (m MemFile) Read(b []byte) (int, error)              { return m.data.Read(b) }
func (m MemFile) Close() error                            { return nil }
func (m MemFile) Size() int64                             { return m.data.Size() }
func (m MemFile) ReadAt(p []byte, off int64) (int, error) { return m.data.ReadAt(p, off) }

func (m MemFile) ToSized() ReaderAtSized {
	return ReaderAtSized{
		ra: m,
		sz: m.Size(),
	}
}

type TimeProvider func() time.Time

func TimeProviderConstNew(ct time.Time) TimeProvider { return func() time.Time { return ct } }

var TimeProviderEpoch TimeProvider = TimeProviderConstNew(TimeEpoch)

type MemFileBuilder struct {
	Name     string
	Reader   io.Reader
	Modified time.Time
	ToBytes  ki.Reader2Bytes
	Mib      MemInfoBuilder
}

func (m MemFileBuilder) build(b []byte) MemFile {
	var sz int64 = int64(len(b))
	var mi MemInfo = m.Mib(m.Modified)(m.Name)(sz)
	return MemFile{
		data:     bytes.NewReader(b),
		fileinfo: mi,
	}
}

var MemFileBuilderDefault = MemFileBuilder{}.Default()

// Default creates MemFileBuilder with default settings.
//
// - ToBytes: UnlimitedRead2Bytes
// - Mib:     MemInfoBuilderDefault
func (m MemFileBuilder) Default() MemFileBuilder {
	m.ToBytes = ki.UnlimitedRead2Bytes
	m.Mib = MemInfoBuilderDefault
	return m
}

func (m MemFileBuilder) WithRead2Bytes(r2b ki.Reader2Bytes) MemFileBuilder {
	m.ToBytes = r2b
	return m
}

func (m MemFileBuilder) WithInfoBuilder(ib MemInfoBuilder) MemFileBuilder {
	m.Mib = ib
	return m
}

func (m MemFileBuilder) WithModified(t time.Time) MemFileBuilder {
	m.Modified = t
	return m
}

func (m MemFileBuilder) WithName(name string) MemFileBuilder {
	m.Name = name
	return m
}

func (m MemFileBuilder) WithReader(r io.Reader) MemFileBuilder {
	m.Reader = r
	return m
}

func (m MemFileBuilder) Build() (MemFile, error) {
	var valid bool = ki.IterFromArr([]bool{
		0 < len(m.Name),
		nil != m.Reader,
		nil != m.ToBytes,
		nil != m.Mib,
	}).All(ki.Identity[bool])

	return ki.ErrorFromBool(
		valid,
		func() (MemFile, error) {
			return ki.ComposeErr(
				m.ToBytes,
				ki.ErrorFuncCreate(m.build),
			)(m.Reader)
		},
		func() error { return fmt.Errorf("Invalid builder") },
	)
}
