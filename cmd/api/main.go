package main

import (
	"2019_2_Covenant/internal/app/apiserver"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "conf", "configs/apiserver.toml", "path to server config")
}

// @title Covenant API
// @version 1.0
// @description Covenant backend server
// @BasePath /api/v1
func main() {
	flag.Parse()

	config := apiserver.NewConfig()

	if _, err := toml.DecodeFile(configPath, config); err != nil {
		log.Fatal(err)
	}

	fmt.Println(config.Address, config.Port, config.Storage)

	server := apiserver.NewAPIServer(config)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
