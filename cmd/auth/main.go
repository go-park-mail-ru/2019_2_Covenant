package main

import (
	"2019_2_Covenant/pkg/auth"
	"database/sql"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"log"
)

type Config struct {
	AuthAddress string `required:"true" split_words:"true"`
	PostgresURL string `required:"true" split_words:"true"`
}

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

	server := auth.NewAuthServer(database)
	defer server.Stop()

	if err := server.Start(config.AuthAddress); err != nil {
		log.Fatal(err)
	}
}