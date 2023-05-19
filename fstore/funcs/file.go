package funcs

import (
	"github.com/rytsh/liz/fstore/generic"
	"github.com/rytsh/liz/loader/file"
)

func init() {
	generic.CallReg.AddFunction("file", new(File).init, "trust")
}

type File struct {
	trust bool
	api   *file.API
}

func (f *File) init(trust bool) *File {
	f.trust = trust
	f.api = file.New()

	return f
}

// Deprecated: Use Write instead.
func (f *File) Save(fileName string, data []byte) (bool, error) {
	return f.Write(fileName, data)
}

func (f *File) Write(fileName string, data []byte) (bool, error) {
	if !f.trust {
		return false, generic.ErrTrustRequired
	}

	if err := f.api.SetRaw(fileName, data); err != nil {
		return false, err
	}

	return true, nil
}

func (f *File) Read(fileName string) ([]byte, error) {
	return f.api.LoadRaw(fileName)
}
