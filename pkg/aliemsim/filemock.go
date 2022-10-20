package aliemsim

import (
	"errors"
	"io"
	"os"
)

// osFS implements FileSystem using the local disk.
type osFS struct{}

// fakeFSErr implements FileSystem forcing open error
type fakeFSErr struct{}

// fakeFS implements FileSystem forcing with no errors
type fakeFS struct{}

var OsFS FileSystem = osFS{}

type MockedFile struct {
	os.FileInfo
}

func (f MockedFile) Close() error                                 { return nil }
func (f MockedFile) Read(p []byte) (int, error)                   { return 1, nil }
func (f MockedFile) ReadAt(p []byte, off int64) (int, error)      { return 0, nil }
func (f MockedFile) Seek(offset int64, whence int) (int64, error) { return 0, nil }
func (f MockedFile) Stat() (os.FileInfo, error)                   { return f, nil }

type FileSystem interface {
	Open(name string) (file, error)
	Stat(name string) (os.FileInfo, error)
}

type file interface {
	io.Closer
	io.Reader
	io.ReaderAt
	io.Seeker
	Stat() (os.FileInfo, error)
}

var FileOpenErrMock = errors.New("error opening file")

func (osFS) Open(name string) (file, error)        { return os.Open(name) }
func (osFS) Stat(name string) (os.FileInfo, error) { return os.Stat(name) }

func (fakeFSErr) Open(name string) (file, error) {
	f := MockedFile{}
	return f, FileOpenErrMock
}
func (fakeFSErr) Stat(name string) (os.FileInfo, error) {
	f := MockedFile{}
	return f.Stat()
}

func (fakeFS) Open(name string) (file, error) {
	f := MockedFile{}
	return f, nil
}
func (fakeFS) Stat(name string) (os.FileInfo, error) {
	f := MockedFile{}
	return f.Stat()
}
