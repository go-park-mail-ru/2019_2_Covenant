package apiserver

type Config struct {
	Address  string `toml:"server_addr"`
	Port     string `toml:"server_port"`
	LogLevel string `toml:"log_level"`
}

func NewConfig() *Config {
	return &Config{
		Address: "127.0.0.1",
		Port:    "3000",
	}
}
