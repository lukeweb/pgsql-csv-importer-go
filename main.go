package main

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	DbConnectionURI string `env:"DB_CONNECTION_URI" validate:"required"`
	Table           string `env:"TABLE_NAME" validate:"required"`
}

func main() {
	var csvFile string

	cfg := LoadConfig[Config]()

	log.Printf("Data import to the table %s started", cfg.Table)

	flag.StringVar(&csvFile, "file", "", "Path to CSV file")
	flag.Parse()

	if csvFile == "" {
		log.Fatal("Provide CSV file path using --file parameter")
	}
}

func LoadConfig[T any]() T {
	var config T

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	err = env.Parse(&config)
	if err != nil {
		panic(err)
	}

	return config
}
