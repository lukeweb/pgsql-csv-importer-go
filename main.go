package main

import (
	"database/sql"
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

	conn, err := sql.Open("postgres", cfg.DbConnectionURI)
	if err != nil {
		log.Fatal("Cannot connect to database")
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error while closing database connection: %v", err)
		}
	}()

	if err := conn.Ping(); err != nil {
		log.Fatalf("Cannot connect to database (Ping failed): %v", err)
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
