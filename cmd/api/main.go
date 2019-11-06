package main

import (
	_ "2019_2_Covenant/docs"
	"2019_2_Covenant/internal/app/apiserver"
	"2019_2_Covenant/internal/app/storage"
	"flag"
	"github.com/BurntSushi/toml"
	"log"
)

var (
	serverConfPath  string
	storageConfPath string
)

func init() {
	flag.StringVar(&serverConfPath, "server", "configs/server.toml", "path to server config")
	flag.StringVar(&storageConfPath, "storage", "configs/storage.toml", "path to storage config")
}

// @title Covenant API
// @version 1.0
// @description Covenant backend server
// @BasePath /api/v1
func main() {
	flag.Parse()

	serverConfig := apiserver.NewConfig()

	if _, err := toml.DecodeFile(serverConfPath, serverConfig); err != nil {
		log.Fatal(err)
	}

	storageConfig := storage.NewConfig("dev")

	if _, err := toml.DecodeFile(storageConfPath, storageConfig); err != nil {
		log.Fatal(err)
	}

	st := storage.NewPGStorage(storageConfig)
	server := apiserver.NewAPIServer(serverConfig, st)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
