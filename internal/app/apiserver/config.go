package apiserver

import "2019_2_Covenant/internal/app/storage"

type Config struct {
	Address string `toml:"server_addr"`
	Port    string `toml:"server_port"`
	Storage *storage.Config
}

func NewConfig() *Config {
	return &Config{
		Address: "127.0.0.1",
		Port:    "3000",
		Storage: storage.NewConfig(),
	}
}
