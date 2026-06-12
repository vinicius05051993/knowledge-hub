package main

import (
	"fmt"
	"log"
	"os"

	"indexer/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
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
		"mysql://%s:%s@tcp(%s:%s)/%s",
		cfg.MySQLUser,
		cfg.MySQLPassword,
		cfg.MySQLHost,
		cfg.MySQLPort,
		cfg.MySQLDatabase,
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