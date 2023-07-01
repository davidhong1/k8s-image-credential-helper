package conf

type IConfigLoader interface {
	Load() (*Config, error)
}
