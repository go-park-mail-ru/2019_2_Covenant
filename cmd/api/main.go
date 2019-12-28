package main

import (
	_ "2019_2_Covenant/docs"
	"2019_2_Covenant/pkg/api"
	"2019_2_Covenant/pkg/auth"
	files "2019_2_Covenant/pkg/file_processor"
	files_repository "2019_2_Covenant/pkg/file_processor/repository"
	"database/sql"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
)

type Config struct {
	AuthAddress  string `required:"true" split_words:"true"`
	APIAddress   string `required:"true" split_words:"true"`
	FilesAddress string `required:"true" split_words:"true"`
	FilesDirPath string `required:"true" split_words:"true"`
	LogLevel     string `required:"true" split_words:"true"`
	PostgresURL  string `required:"true" split_words:"true"`
}

// @title Covenant API
// @version 1.0
// @description Covenant backend server
// @BasePath /api/v1
func main() {
	var config Config
	err := envconfig.Process("Covenant", &config)
	if err != nil {
		log.Fatal(err)
	}

	database, err := sql.Open("postgres", config.PostgresURL)
	if err != nil {
		log.Fatal(err)
	}

	defer database.Close()
	if err := database.Ping(); err != nil {
		log.Fatal(err)
	}
	storage := api.NewPGStorage(database)

	authServiceConnection, err := grpc.Dial(config.AuthAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer authServiceConnection.Close()
	auth := auth.NewAuthClient(authServiceConnection)


	filesServiceConnection, err := grpc.Dial(config.FilesAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer filesServiceConnection.Close()
	files := files.NewFilesClient(filesServiceConnection)
	file_repository := files_repository.NewFileRepository(files)

	server := api.NewAPIServer(storage, auth, file_repository)
	defer server.Stop()

	if err := server.Start(config.APIAddress, config.LogLevel, config.FilesDirPath); err != nil {
		log.Fatal(err)
	}
}
