package loader

type Configs []Config

type Config struct {
	Export        string
	ExportToValue bool
	Statics       []ConfigStatic
	Dynamics      []ConfigDynamic
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
	// Path is the location in consul KV
	Path string
	// PathPrefix default is empty
	PathPrefix string
	// Raw to load as raw, don't mix with other loaders
	Raw bool
	// Codec YAML,JSON,TOML default is YAML
	Codec string
}

type ConfigVault struct {
	Path string
	// PathPrefix default is empty, path_prefix is must!
	PathPrefix string
	// AppRoleBasePath default is auth/approle/login, not need to set
	AppRoleBasePath string
	// AdditionalPaths additional paths to get from extra content, default is none
	AdditionalPaths []ConfigVaultAdditional
}

type ConfigVaultAdditional struct {
	// Map is the where to add as trace/config -> ["trace"]["config"]
	Map string
	// Path show location in vault config
	Path string
}

type ConfigFile struct {
	// Path is the file location, [toml, yml, yaml, json] supported
	Path string
	// Raw to load as raw, don't mix with other loaders
	Raw bool
}

type ConfigContent struct {
	// Codec YAML,JSON,TOML default is YAML
	Codec    string
	Content  string
	Raw      bool
	Template bool
}
