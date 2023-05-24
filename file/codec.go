package file

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

type Codec interface {
	Encode(w io.Writer, v any) error
	Decode(r io.Reader, v interface{}) error
}

// YAML is a codec for YAML file.
type YAML struct{}

var _ Codec = YAML{}

func (YAML) Encode(w io.Writer, v any) error {
	encode := yaml.NewEncoder(w)

	if err := encode.Encode(v); err != nil {
		return fmt.Errorf("failed to encode YAML: %w", err)
	}

	return nil
}

func (YAML) Decode(r io.Reader, v interface{}) error {
	decode := yaml.NewDecoder(r)

	if err := decode.Decode(v); err != nil {
		return fmt.Errorf("failed to decode YAML: %w", err)
	}

	return nil
}

// JSON is a codec for JSON file.
type JSON struct {
	Indent string
}

var _ Codec = JSON{}

func (j JSON) Encode(w io.Writer, v any) error {
	encode := json.NewEncoder(w)

	if j.Indent != "" {
		encode.SetIndent("", j.Indent)
	}

	if err := encode.Encode(v); err != nil {
		return fmt.Errorf("json encode: %w", err)
	}

	return nil
}

func (JSON) Decode(r io.Reader, v interface{}) error {
	decode := json.NewDecoder(r)

	if err := decode.Decode(v); err != nil {
		return fmt.Errorf("json decode: %w", err)
	}

	return nil
}

// TOML is a codec for TOML file.
type TOML struct{}

var _ Codec = TOML{}

func (TOML) Encode(w io.Writer, v any) error {
	encode := toml.NewEncoder(w)

	if err := encode.Encode(v); err != nil {
		return fmt.Errorf("toml encode: %w", err)
	}

	return nil
}

func (TOML) Decode(r io.Reader, v interface{}) error {
	decode := toml.NewDecoder(r)

	if _, err := decode.Decode(v); err != nil {
		return fmt.Errorf("toml decode: %w", err)
	}

	return nil
}

type RAW struct{}

var _ Codec = RAW{}

func (RAW) Encode(w io.Writer, v any) error {
	vByte, ok := v.([]byte)
	if !ok {
		return fmt.Errorf("raw decode: expected []byte, got %T", v)
	}

	_, err := w.Write(vByte)
	return err
}

func (RAW) Decode(r io.Reader, v any) error {
	var vByte *[]byte

	switch vTyped := v.(type) {
	case *[]byte:
		vByte = vTyped
	case *interface{}:
		vByteInf, ok := (*vTyped).([]byte)
		if !ok {
			return fmt.Errorf("raw decode: expected *[]byte, got %T", *vTyped)
		}

		vByte = &vByteInf
	default:
		return fmt.Errorf("raw decode: expected *[]byte, got %T", v)
	}

	buffer := bytes.NewBuffer(*vByte)
	// read and write to vByte
	if _, err := buffer.ReadFrom(r); err != nil {
		return err
	}

	// update v
	if _, ok := v.(*interface{}); ok {
		*v.(*interface{}) = buffer.Bytes()
	}

	if _, ok := v.(*[]byte); ok {
		*v.(*[]byte) = buffer.Bytes()
	}

	return nil
}
