package main

import (
	"fmt"
	"log"
	"os"

	"indexer/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlserver"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	if len(os.Args) < 2 {
		log.Fatal("usage: migrate [up|down]")
	}

	cfg, err := config.Load()

	if err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf(
		"sqlserver://%s:%s@%s:%s?database=%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	m, err := migrate.New(
		"file://migrations",
		dsn,
	)

	if err != nil {
		log.Fatal(err)
	}

	command := os.Args[1]

	switch command {

	case "up":

		err = m.Up()

		if err != nil &&
			err != migrate.ErrNoChange {

			log.Fatal(err)
		}

		fmt.Println("migrations applied")

	case "down":

		err = m.Down()

		if err != nil &&
			err != migrate.ErrNoChange {

			log.Fatal(err)
		}

		fmt.Println("migrations reverted")

	default:

		log.Fatal("usage: migrate [up|down]")
	}
}