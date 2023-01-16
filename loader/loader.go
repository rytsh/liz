package loader

import (
	"context"
	"fmt"
	"path"
	"strings"
	"sync"

	"github.com/rytsh/liz/loader/file"
)

// Load loads all configs to export location.
// If loads with dynamic config, cancel context to stop loading.
// Export is used to export data if ExportToValue enabled.
// Currently just static config supported to export data.
func (c Configs) Load(ctx context.Context, wg *sync.WaitGroup, cache *Cache, export *Data) error {
	if cache == nil {
		cache = &Cache{}
	}

	cache.set()

	for _, config := range c {
		if err := config.load(ctx, wg, cache, export); err != nil {
			return err
		}
	}

	return nil
}

func (c Config) load(ctx context.Context, wg *sync.WaitGroup, cache *Cache, export *Data) error {
	to := Data{}

	for _, static := range c.Statics {
		if err := static.load(ctx, &to, cache); err != nil {
			return err
		}
	}

	if len(c.Dynamics) > 0 {
		for _, dynamic := range c.Dynamics {
			waitCtx, err := dynamic.load(ctx, wg, &to, cache, c.Export)
			if err != nil {
				return err
			}

			// wait to first load
			if waitCtx != nil {
				<-waitCtx.Done()
			}
		}
	} else {
		if c.Export != "" {
			if to.Raw != nil {
				if err := cache.File.SetRaw(c.Export, to.Raw); err != nil {
					return err
				}
			} else {
				if err := cache.File.SetWithCodec(c.Export, to); err != nil {
					return err
				}
			}
		}
	}

	if c.ExportToValue && export != nil {
		*export = to
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

		if c.Consul.Raw {
			to.Raw = data
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

			to.Merge(vMap)
		}
	}

	if c.Vault != nil {
		cache.Vault.AppRoleBasePath = c.Vault.AppRoleBasePath
		// load additional secrets
		for _, additional := range c.Vault.AdditionalPaths {
			v, err := cache.Vault.LoadMap(ctx, c.Vault.PathPrefix, additional.Path)
			if err != nil {
				return err
			}

			if additional.Map != "" {
				maps := strings.Split(additional.Map, "/")

				mapDef := map[string]interface{}{}
				mapRange := mapDef
				for _, m := range maps {
					if m == maps[len(maps)-1] {
						mapRange[m] = v
						break
					}

					mapRange[m] = map[string]interface{}{}
					mapRange = mapRange[m].(map[string]interface{})
				}

				v = mapDef
			}

			to.Merge(v)
		}

		// load main secret
		v, err := cache.Vault.LoadMap(ctx, c.Vault.PathPrefix, c.Vault.Path)
		if err != nil {
			return err
		}

		to.Merge(v)
	}

	if c.File != nil {
		data, err := cache.File.LoadRaw(c.File.Path)
		if err != nil {
			return err
		}

		if c.File.Raw {
			to.Raw = data
		} else {
			// convert to map
			var vMap map[string]interface{}

			if err := cache.File.Load(c.File.Path, &vMap); err != nil {
				return err
			}

			to.Merge(vMap)
		}
	}

	if c.Content != nil {
		content := c.Content.Content
		if c.Content.Template {
			v, err := cache.Template.Execute(nil, content)
			if err != nil {
				return err
			}

			content = v
		}
		if c.Content.Raw {
			to.Raw = []byte(content)
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

			to.Merge(vMap)
		}
	}

	return nil
}

func (c ConfigDynamic) load(ctx context.Context, wg *sync.WaitGroup, to *Data, cache *Cache, filePath string) (context.Context, error) {
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
					if c.Consul.Raw {
						to.Raw = data
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
					}

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
			}
		}()
	}

	return waitContext, nil
}
