package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

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

	EmbeddingDimension int
}

func Load() (*Config, error) {

	err := godotenv.Load()

	if err != nil {
		log.Println(".env not found")
	}

	embeddingDimension := 384

	if value := strings.TrimSpace(
		os.Getenv("EMBEDDING_DIMENSION"),
	); value != "" {

		parsed, err := strconv.Atoi(value)

		if err != nil {

			return nil, fmt.Errorf(
				"invalid EMBEDDING_DIMENSION: %w",
				err,
			)
		}

		embeddingDimension = parsed
	}

	cfg := &Config{
		AppPort: os.Getenv("APP_PORT"),

		MySQLHost:     os.Getenv("MYSQL_HOST"),
		MySQLPort:     os.Getenv("MYSQL_PORT"),
		MySQLDatabase: os.Getenv("MYSQL_DATABASE"),
		MySQLUser:     os.Getenv("MYSQL_USER"),
		MySQLPassword: os.Getenv("MYSQL_PASSWORD"),

		OpenSearchHost: os.Getenv("OPENSEARCH_HOST"),
		OpenSearchPort: os.Getenv("OPENSEARCH_PORT"),

		EmbeddingDimension: embeddingDimension,
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {

	required := map[string]string{
		"APP_PORT":        c.AppPort,
		"MYSQL_HOST":      c.MySQLHost,
		"MYSQL_PORT":      c.MySQLPort,
		"MYSQL_DATABASE":  c.MySQLDatabase,
		"MYSQL_USER":      c.MySQLUser,
		"MYSQL_PASSWORD":  c.MySQLPassword,
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

	if c.EmbeddingDimension <= 0 {

		return fmt.Errorf(
			"EMBEDDING_DIMENSION must be greater than zero",
		)
	}

	return nil
}