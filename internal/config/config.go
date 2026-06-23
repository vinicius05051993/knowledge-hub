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

	UseManagedIdentity bool

	OpenSearchUrl string
}

func Load() (*Config, error) {

	if err := godotenv.Load(); err != nil {
		log.Println(".env not found")
	}

	cfg := &Config{
		AppPort: os.Getenv("APP_PORT"),

		DBHost: os.Getenv("DB_HOST"),
		DBPort: os.Getenv("DB_PORT"),
		DBName: os.Getenv("DB_NAME"),

		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),

		UseManagedIdentity: strings.EqualFold(
			os.Getenv("USE_MANAGED_IDENTITY"),
			"true",
		),

		OpenSearchUrl: os.Getenv("OPENSEARCH_URL"),
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
	}

	for key, value := range required {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("%s is required", key)
		}
	}

	if !c.UseManagedIdentity {

		if strings.TrimSpace(c.DBUser) == "" {
			return fmt.Errorf("DB_USER is required")
		}

		if strings.TrimSpace(c.DBPassword) == "" {
			return fmt.Errorf("DB_PASSWORD is required")
		}
	}

	return nil
}