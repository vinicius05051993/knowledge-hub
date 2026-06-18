package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string

	OpenSearchHost string
	OpenSearchPort string
}

func Load() (*Config, error) {

	err := godotenv.Load()

	if err != nil {
		log.Println(".env not found")
	}

	cfg := &Config{
		AppPort: os.Getenv("APP_PORT"),

		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBName:     os.Getenv("DB_NAME"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),

		OpenSearchHost: os.Getenv("OPENSEARCH_HOST"),
		OpenSearchPort: os.Getenv("OPENSEARCH_PORT"),
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {

	required := map[string]string{
		"APP_PORT":        c.AppPort,
		"DB_HOST":         c.DBHost,
		"DB_PORT":         c.DBPort,
		"DB_NAME":         c.DBName,
		"DB_USER":         c.DBUser,
		"DB_PASSWORD":     c.DBPassword,
		"OPENSEARCH_HOST": c.OpenSearchHost,
		"OPENSEARCH_PORT": c.OpenSearchPort,
	}

	for key, value := range required {

		if strings.TrimSpace(value) == "" {

			return fmt.Errorf(
				"%s is required",
				key,
			)
		}
	}

	return nil
}