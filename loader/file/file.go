package file

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

var (
	defaultFileFlag = os.O_RDWR | os.O_CREATE | os.O_TRUNC
	defaultFilePerm = 0644
)

type API struct {
	Codec map[string]codec

	FileFlag int
	Perm     int
}

func (a *API) LoadRaw(path string) ([]byte, error) {
	// open file
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", path, err)
	}

	// read file
	b, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	return b, nil
}

func (a *API) LoadMap(path string) (map[string]interface{}, error) {
	// load raw
	b, err := a.LoadRaw(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load raw: %w", err)
	}

	// unmarshal
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	return m, nil
}

func (a *API) Set(path string, data []byte) error {
	// open file
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", path, err)
	}

	// write file
	if _, err := f.Write(data); err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}

	return nil
}
