package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	MySQLHost     string
	MySQLPort     string
	MySQLDatabase string
	MySQLUser     string
	MySQLPassword string

	OpenSearchHost string
	OpenSearchPort string
}

func Load() *Config {

	err := godotenv.Load()

	if err != nil {
		log.Println(".env not found")
	}

	return &Config{
		AppPort: os.Getenv("APP_PORT"),

		MySQLHost:     os.Getenv("MYSQL_HOST"),
		MySQLPort:     os.Getenv("MYSQL_PORT"),
		MySQLDatabase: os.Getenv("MYSQL_DATABASE"),
		MySQLUser:     os.Getenv("MYSQL_USER"),
		MySQLPassword: os.Getenv("MYSQL_PASSWORD"),
		OpenSearchHost: os.Getenv("OPENSEARCH_HOST"),
		OpenSearchPort: os.Getenv("OPENSEARCH_PORT"),
	}
}