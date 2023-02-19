package loader

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"strings"
	"sync"

	"github.com/rytsh/liz/loader/file"
	"github.com/rytsh/liz/utils/templatex"
)

type Call func(context.Context, string, map[string]interface{})

// Load loads all configs to export location.
// If loads with dynamic config, cancel context to stop loading.
func (c Configs) Load(ctx context.Context, wg *sync.WaitGroup, cache *Cache, call Call) error {
	if cache == nil {
		cache = &Cache{}
	}

	cache.set()

	for _, config := range c {
		if err := config.load(ctx, wg, cache, call); err != nil {
			return err
		}
	}

	return nil
}

func (c Config) load(ctx context.Context, wg *sync.WaitGroup, cache *Cache, call Call) error {
	to := Data{}

	for _, static := range c.Statics {
		if err := static.load(ctx, &to, cache); err != nil {
			return err
		}
	}

	if len(c.Dynamics) > 0 {
		for _, dynamic := range c.Dynamics {
			waitCtx, err := dynamic.load(ctx, wg, &to, cache, c.Export, c.Name, call)
			if err != nil {
				return err
			}

			// wait to first load
			if waitCtx != nil {
				<-waitCtx.Done()
			}
		}
	} else {
		if to.Raw != nil {
			to.AddHold(c.Name, to.Raw)
		} else {
			to.AddHold(c.Name, to.Map)
		}

		if c.Export != "" {
			if to.Raw != nil {
				if err := cache.File.SetRaw(c.Export, to.Raw); err != nil {
					return err
				}
			} else {
				if err := cache.File.SetWithCodec(c.Export, to.Map); err != nil {
					return err
				}
			}
		}

		if call != nil {
			call(ctx, c.Name, to.Hold)
		}
	}

	return nil
}

func (c ConfigStatic) load(ctx context.Context, to *Data, cache *Cache) error {
	if c.Consul != nil {
		contentPath := path.Join(c.Consul.PathPrefix, c.Consul.Path)

		data, err := cache.Consul.LoadRaw(ctx, contentPath)
		if err != nil {
			return err
		}

		if c.Consul.Template {
			v, err := cache.Template.ExecuteBuffer(templatex.WithContent(string(data)), templatex.WithData(to.Hold))
			if err != nil {
				return err
			}

			data = v
		}

		var dataProcessed interface{}

		if c.Consul.Raw {
			if c.Consul.Map != "" {
				vMap := MapPath(c.Consul.Map, data).(map[string]interface{})
				to.Merge(vMap)
				dataProcessed = vMap
			} else {
				to.Raw = data
				dataProcessed = data
			}
		} else {
			// convert to map
			var vMap map[string]interface{}
			codecTxt := strings.ToUpper(c.Consul.Codec)
			if codecTxt == "" {
				codecTxt = "YAML"
			}

			codec := cache.File.Codec[codecTxt]
			if codec == nil {
				return fmt.Errorf("codec %s not found", codecTxt)
			}

			if err := cache.File.LoadContent(data, &vMap, codec); err != nil {
				return err
			}

			vMap, ok := MapPath(c.Consul.Map, InnerPath(c.Consul.InnerPath, vMap)).(map[string]interface{})
			if !ok {
				return fmt.Errorf("consul mapping error")
			}

			to.Merge(vMap)
			dataProcessed = vMap
		}

		to.AddHold(c.Consul.Name, dataProcessed)
	}

	if c.Vault != nil {
		cache.Vault.AppRoleBasePath = c.Vault.AppRoleBasePath

		// load main secret
		vMap, err := cache.Vault.LoadMap(ctx, c.Vault.PathPrefix, c.Vault.Path)
		if err != nil {
			return err
		}

		if c.Vault.Template {
			data, err := json.Marshal(vMap)
			if err != nil {
				return err
			}

			vRendered, err := cache.Template.ExecuteBuffer(templatex.WithContent(string(data)), templatex.WithData(to.Hold))
			if err != nil {
				return err
			}

			var vX map[string]interface{}
			if err := json.Unmarshal(vRendered, &vX); err != nil {
				return err
			}

			vMap = vX
		}

		vMap, ok := MapPath(c.Vault.Map, InnerPath(c.Vault.InnerPath, vMap)).(map[string]interface{})
		if !ok {
			return fmt.Errorf("vault mapping error")
		}

		to.Merge(vMap)
		to.AddHold(c.Vault.Name, vMap)
	}

	if c.File != nil {
		data, err := cache.File.LoadRaw(c.File.Path)
		if err != nil {
			return err
		}

		if c.File.Template {
			v, err := cache.Template.ExecuteBuffer(templatex.WithContent(string(data)), templatex.WithData(to.Hold))
			if err != nil {
				return err
			}

			data = v
		}

		var dataProcessed interface{}

		if c.File.Raw {
			if c.File.Map != "" {
				vMap := MapPath(c.File.Map, data).(map[string]interface{})
				to.Merge(vMap)
				dataProcessed = vMap
			} else {
				to.Raw = data
				dataProcessed = data
			}
		} else {
			// convert to map
			var vMap map[string]interface{}

			if err := cache.File.Load(c.File.Path, &vMap); err != nil {
				return err
			}

			vMap, ok := MapPath(c.File.Map, InnerPath(c.File.InnerPath, vMap)).(map[string]interface{})
			if !ok {
				return fmt.Errorf("file mapping error")
			}

			to.Merge(vMap)
			dataProcessed = vMap
		}

		to.AddHold(c.File.Name, dataProcessed)
	}

	if c.Content != nil {
		content := c.Content.Content
		if c.Content.Template {
			v, err := cache.Template.ExecuteBuffer(templatex.WithContent(content), templatex.WithData(to.Hold))
			if err != nil {
				return err
			}

			content = string(v)
		}

		var dataProcessed interface{}

		if c.Content.Raw {
			if c.Content.Map != "" {
				vMap := MapPath(c.Content.Map, []byte(content)).(map[string]interface{})
				to.Merge(vMap)
				dataProcessed = vMap
			} else {
				to.Raw = []byte(content)
				dataProcessed = []byte(content)
			}
		} else {
			// convert to map
			var vMap map[string]interface{}

			codecTxt := strings.ToUpper(c.Content.Codec)
			if codecTxt == "" {
				codecTxt = "YAML"
			}

			codec := cache.File.Codec[codecTxt]
			if codec == nil {
				return fmt.Errorf("codec %s not found", codecTxt)
			}

			if err := cache.File.LoadContent([]byte(content), &vMap, codec); err != nil {
				return err
			}

			vMap, ok := MapPath(c.Content.Map, InnerPath(c.Content.InnerPath, vMap)).(map[string]interface{})
			if !ok {
				return fmt.Errorf("content mapping error")
			}

			to.Merge(vMap)
			dataProcessed = vMap
		}

		to.AddHold(c.Content.Name, dataProcessed)
	}

	return nil
}

func (c ConfigDynamic) load(ctx context.Context, wg *sync.WaitGroup, to *Data, cache *Cache, filePath, holdName string, call Call) (context.Context, error) {
	if wg == nil {
		wg = &sync.WaitGroup{}
	}

	var waitContext context.Context

	if c.Consul != nil {
		var codec file.Codec
		if !c.Consul.Raw {
			codecTxt := strings.ToUpper(c.Consul.Codec)
			if codecTxt == "" {
				codecTxt = "YAML"
			}

			codec = cache.File.Codec[codecTxt]
			if codec == nil {
				return nil, fmt.Errorf("codec %s not found", codecTxt)
			}
		}

		contentPath := path.Join(c.Consul.PathPrefix, c.Consul.Path)

		ch, cancel, err := cache.Consul.DynamicValue(ctx, wg, contentPath)
		if err != nil {
			return nil, err
		}

		var waitCancel context.CancelFunc
		waitContext, waitCancel = context.WithCancel(ctx)

		recordToMap := copyMap(to.Map)

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer cancel()

			i := 0

			for {
				if i == 0 {
					i++
				} else if i == 1 {
					waitCancel()
					i++
				}

				select {
				case <-ctx.Done():
					return
				case data := <-ch:
					if c.Consul.Template {
						v, err := cache.Template.ExecuteBuffer(templatex.WithContent(string(data)), templatex.WithData(to.Hold))
						if err != nil {
							logFromCtx(ctx).Warn("failed to execute consul template", "err", err.Error())
							continue
						}

						data = v
					}

					if c.Consul.Raw {
						to.Raw = data
						to.AddHold(c.Consul.Name, data)
					} else {
						// convert to map
						var vMap map[string]interface{}

						if err := cache.File.LoadContent(data, &vMap, codec); err != nil {
							logFromCtx(ctx).Warn("failed to load consul data", "err", err.Error())
							continue
						}

						// get back old map
						to.Map = copyMap(recordToMap)
						to.Merge(vMap)
						to.AddHold(c.Consul.Name, vMap)
					}

					if to.Raw != nil {
						to.AddHold(holdName, to.Raw)
					} else {
						to.AddHold(holdName, to.Map)
					}

					if filePath != "" {
						if to.Raw != nil {
							if err := cache.File.SetRaw(filePath, to.Raw); err != nil {
								logFromCtx(ctx).Warn("failed to save dynamic consul data to file", "filePath", filePath, "err", err.Error())
							}
						} else {
							if err := cache.File.SetWithCodec(filePath, to.Map); err != nil {
								logFromCtx(ctx).Warn("failed to save dynamic consul data to file", "filePath", filePath, "err", err.Error())
							}
						}
					}

					if call != nil {
						call(ctx, holdName, to.Hold)
					}
				}
			}
		}()
	}

	return waitContext, nil
}
