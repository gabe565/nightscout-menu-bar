package config

type Option func(conf *Config)

func WithVersion(version string) Option {
	return func(conf *Config) {
		conf.Version = version
	}
}

func WithData(data Data) Option {
	return func(conf *Config) {
		conf.data.Store(&data)
	}
}
