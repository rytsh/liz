package file

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

var (
	defaultFileFlag int         = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	defaultFilePerm fs.FileMode = 0o644
)

type API struct {
	Codec map[string]Codec

	FileFlag int
	Perm     fs.FileMode
}

func New() *API {
	json := JSON{Indent: "  "}
	yaml := YAML{}
	toml := TOML{}

	return &API{
		Codec: map[string]Codec{
			"JSON":  json,
			".json": json,
			"YAML":  yaml,
			".yaml": yaml,
			".yml":  yaml,
			"TOML":  toml,
			".toml": toml,
		},
	}
}

func (a *API) openFileRead(path string) (*os.File, error) {
	// open file
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", path, err)
	}

	return f, nil
}

func (a *API) openFileWrite(path string) (*os.File, error) {
	if a.FileFlag == 0 {
		a.FileFlag = defaultFileFlag
	}

	if a.Perm == 0 {
		a.Perm = defaultFilePerm
	}

	// open file
	f, err := os.OpenFile(path, a.FileFlag, a.Perm)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", path, err)
	}

	return f, nil
}

func (a *API) getCodec(path string) (Codec, error) {
	// get codec
	ext := filepath.Ext(path)
	codec, ok := a.Codec[ext]
	if !ok {
		return nil, fmt.Errorf("unsupported file extension %s", ext)
	}

	return codec, nil
}

func (a *API) LoadRaw(path string) ([]byte, error) {
	// open file
	f, err := a.openFileRead(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	// read file
	b, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	return b, nil
}

func (a *API) Load(path string, dst interface{}) error {
	// open file
	f, err := a.openFileRead(path)
	if err != nil {
		return err
	}

	defer f.Close()

	// check the path extension
	ext := filepath.Ext(path)
	codec := a.Codec[ext]
	if codec == nil {
		return fmt.Errorf("failed to find codec for extension %s", ext)
	}

	if err := codec.Decode(f, dst); err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}

	return nil
}

func (a *API) LoadContent(v []byte, dst interface{}, codec Codec) error {
	if codec == nil {
		return fmt.Errorf("failed codec is nil")
	}

	f := bytes.NewReader(v)

	if err := codec.Decode(f, dst); err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}

	return nil
}

func (a *API) SetWithCodec(path string, data any) error {
	// open file
	f, err := a.openFileWrite(path)
	if err != nil {
		return err
	}

	defer f.Close()

	codec, err := a.getCodec(path)
	if err != nil {
		return err
	}

	if err := codec.Encode(f, data); err != nil {
		return err
	}

	return nil
}

func (a *API) SetRaw(path string, data []byte) error {
	// open file
	f, err := a.openFileWrite(path)
	if err != nil {
		return err
	}

	defer f.Close()

	// write file
	if _, err := f.Write(data); err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}

	return nil
}
