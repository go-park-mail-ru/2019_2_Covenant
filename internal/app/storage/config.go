package storage

type Config struct {
	DBUrl string `toml:"db_url"`
}

func NewConfig() *Config {
	return &Config{}
}
