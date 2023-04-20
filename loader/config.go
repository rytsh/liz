package loader

type Configs []Config

type Config struct {
	// Name for export value, default is empty.
	Name       string
	Export     string
	FilePerm   string
	FolderPerm string
	Statics    []ConfigStatic
	Dynamics   []ConfigDynamic
}

type ConfigStatic struct {
	Consul  *ConfigConsul
	Vault   *ConfigVault
	File    *ConfigFile
	Content *ConfigContent
}

type ConfigDynamic struct {
	Consul *ConfigConsul
}

type ConfigConsul struct {
	// Name for export, default is empty.
	Name string
	// Path is the location in consul KV.
	Path string
	// PathPrefix default is empty.
	PathPrefix string
	// Raw to load as raw, don't mix with other loaders.
	Raw bool
	// Codec YAML,JSON,TOML default is YAML.
	Codec string
	// InnerPath is get the inner path from vault response, / separated as db/settings.
	// Cannot work with Raw.
	InnerPath string
	// Map is the wrapper map, / separated as db/settings.
	Map string
	// Template to run go template after the load.
	Template bool
	// base64 to decode the content.
	Base64 bool
}

type ConfigVault struct {
	// Name for export, default is empty.
	Name string
	Path string
	// PathPrefix default is empty, path_prefix is must!
	PathPrefix string
	// AppRoleBasePath default is auth/approle/login, not need to set.
	AppRoleBasePath string
	// InnerPath is get the inner path from vault response, / separated as db/settings.
	InnerPath string
	// Map is the wrapper map, / separated as db/settings.
	Map string
	// Template to run go template after the load.
	Template bool
	// base64 to decode the content.
	Base64 bool
}

type ConfigFile struct {
	// Name for export, default is empty.
	Name string
	// Path is the file location, [toml, yml, yaml, json] supported.
	Path string
	// Raw to load as raw, don't mix with other loaders.
	Raw bool
	// InnerPath is get the inner path from vault response, / separated as db/settings.
	// Cannot work with Raw.
	InnerPath string
	// Map is the wrapper map, / separated as db/settings.
	Map string
	// Template to run go template after the load.
	Template bool
	// base64 to decode the content.
	Base64 bool
}

type ConfigContent struct {
	// Name for export, default is empty.
	Name string
	// Codec YAML,JSON,TOML default is YAML.
	Codec   string
	Content string
	Raw     bool
	// InnerPath is get the inner path from vault response, / separated as db/settings.
	// Cannot work with Raw.
	InnerPath string
	// Map is the wrapper map, / separated as db/settings.
	Map string
	// Template to run go template after the load.
	Template bool
	// base64 to decode the content.
	Base64 bool
}
