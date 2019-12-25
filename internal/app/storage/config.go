package storage

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	DBUrl string `toml:"db_url"`
	Mod string
}

func NewConfig(mod string) *Config {
	return &Config{
		Mod: mod,
	}
}

func (c *Config) GetURL() string {
	var url string

	switch c.Mod {
	case "dev":
		url = c.DBUrl
	case "prod":
		url = getFromEnv()
	}

	return url
}

func getFromEnv() string {
	var (
		user = os.Getenv("DB_USER")
		host = os.Getenv("DB_HOST")
		port = os.Getenv("DB_PORT")
		dbName = os.Getenv("DB_NAME")
	)

	var url []string
	url = append(
		url,
		fmt.Sprintf("user=%s", user),
		fmt.Sprintf("host=%s", host),
		fmt.Sprintf("port=%s", port),
		fmt.Sprintf("dbname=%s", dbName),
		"sslmode=disable",
	)

	return strings.Join(url, " ")
}
