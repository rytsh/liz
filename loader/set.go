package loader

import (
	"github.com/rytsh/liz/loader/consul"
	"github.com/rytsh/liz/loader/file"
	"github.com/rytsh/liz/loader/vault"
	"github.com/rytsh/liz/utils/mapx"
	"github.com/rytsh/liz/utils/templatex"
)

type Data struct {
	Map  map[string]interface{}
	Raw  []byte
	Hold map[string]interface{}
}

func (d *Data) AddHold(k string, v interface{}) {
	if k == "" {
		return
	}

	if d.Hold == nil {
		d.Hold = map[string]interface{}{}
	}

	d.Hold[k] = v
}

func (d *Data) Merge(v map[string]interface{}) {
	if d.Map == nil {
		d.Map = map[string]interface{}{}
	}

	mapx.Merge(v, d.Map)
}

type Cache struct {
	Consul   *consul.API
	Vault    *vault.API
	File     *file.API
	Template *templatex.Template
}

func (c *Cache) set() {
	if c.Consul == nil {
		c.Consul = &consul.API{}
	}

	if c.Vault == nil {
		c.Vault = &vault.API{}
	}

	if c.File == nil {
		c.File = file.New()
	}

	if c.Template == nil {
		c.Template = templatex.New()
	}
}

func copyMap(in map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	for k, v := range in {
		vm, ok := v.(map[string]interface{})
		if ok {
			out[k] = copyMap(vm)
		} else {
			out[k] = v
		}
	}

	return out
}
