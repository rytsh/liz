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
	defaultFileFlag   int         = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	defaultFilePerm   fs.FileMode = 0o644
	defaultFolderPerm fs.FileMode = 0o755
)

type API struct {
	Codec map[string]Codec
}

func New() *API {
	json := JSON{Indent: "  "}
	yaml := YAML{}
	toml := TOML{}
	raw := RAW{}

	return &API{
		Codec: map[string]Codec{
			"JSON":  json,
			".json": json,
			"YAML":  yaml,
			".yaml": yaml,
			".yml":  yaml,
			"TOML":  toml,
			".toml": toml,
			"RAW":   raw,
		},
	}
}

// Open opens a file with options.
func (a *API) OpenFile(path string, opts ...Option) (*os.File, error) {
	options, err := a.readOptions(opts)
	if err != nil {
		return nil, err
	}

	// open file
	f, err := a.openFileWrite(path, options)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (a *API) openFileRead(path string) (*os.File, error) {
	// open file
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", path, err)
	}

	return f, nil
}

func (a *API) openFileWrite(path string, opts options) (*os.File, error) {
	// create folder if not exist
	folder := filepath.Dir(path)
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		folderPerm := defaultFolderPerm
		if opts.folderPerm != nil {
			folderPerm = *opts.folderPerm
		}

		if err := os.MkdirAll(folder, folderPerm); err != nil {
			return nil, fmt.Errorf("failed to create folder %s: %w", folder, err)
		}
	}

	filePerm := defaultFilePerm
	if opts.filePerm != nil {
		filePerm = *opts.filePerm
	}

	fileFlag := defaultFileFlag
	if opts.fileFlag != nil {
		fileFlag = *opts.fileFlag
	}

	// open file
	f, err := os.OpenFile(path, fileFlag, filePerm)
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

// LoadWithReturn loads the file with the specified codec and returns the result.
func (a *API) LoadWithReturn(path string, codec Codec) (any, error) {
	var dst any
	codec, ok := a.selectCodec(path, codec)
	if ok {
		dst = map[string]any{}
	} else {
		dst = []byte{}
	}

	err := a.LoadWithCodec(path, &dst, codec)
	return dst, err
}

// LoadWithCodec loads the file with the specified codec.
//
// If the codec is nil, the codec is determined by the file extension.
func (a *API) LoadWithCodec(path string, dst interface{}, codec Codec) error {
	codec, _ = a.selectCodec(path, codec)

	// open file
	f, err := a.openFileRead(path)
	if err != nil {
		return err
	}

	defer f.Close()

	if err := codec.Decode(f, dst); err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}

	return nil
}

func (a *API) selectCodec(path string, codec Codec) (Codec, bool) {
	var ok = true
	if codec == nil {
		ext := filepath.Ext(path)

		codec, ok = a.Codec[ext]
		if !ok {
			codec = a.Codec["RAW"]
		}
	}

	return codec, ok
}

func (a *API) SetWithCodec(path string, data any, opts ...Option) error {
	options, err := a.readOptions(opts)
	if err != nil {
		return err
	}

	// open file
	f, err := a.openFileWrite(path, options)
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

func (a *API) SetRaw(path string, data []byte, opts ...Option) error {
	options, err := a.readOptions(opts)
	if err != nil {
		return err
	}

	// open file
	f, err := a.openFileWrite(path, options)
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

func (a *API) readOptions(opts []Option) (options, error) {
	var options options
	for _, opt := range opts {
		if err := opt(&options); err != nil {
			return options, err
		}
	}

	return options, nil
}
