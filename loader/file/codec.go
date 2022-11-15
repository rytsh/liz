package file

import (
	"encoding/json"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

type YAML struct{}

func (y *YAML) Marshal(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

func (y *YAML) Unmarshal(b []byte, v interface{}) error {
	return yaml.Unmarshal(b, v)
}

type JSON struct{}

func (j *JSON) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (j *JSON) Unmarshal(b []byte, v interface{}) error {
	return json.Unmarshal(b, v)
}

type TOML struct{}

func (t *TOML) Marshal(v interface{}) ([]byte, error) {

	toml.NewEncoder()

	return nil, nil
}

func (t *TOML) Unmarshal(b []byte, v interface{}) error {
	return nil
}
