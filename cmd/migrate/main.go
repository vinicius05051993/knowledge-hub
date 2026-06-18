package main

import (
	"fmt"
	"log"
	"os"

	"indexer/internal/config"
	"indexer/internal/database"

	"github.com/golang-migrate/migrate/v4"
	mssql "github.com/golang-migrate/migrate/v4/database/sqlserver"
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

	db, _, err := database.OpenDatabase(*cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	driver, err := mssql.WithInstance(db, &mssql.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		cfg.DBName,
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	switch os.Args[1] {

	case "up":

		err = m.Up()
		if err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}

		fmt.Println("migrations applied")

	case "down":

		err = m.Down()
		if err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}

		fmt.Println("migrations reverted")

	default:
		log.Fatal("usage: migrate [up|down]")
	}
}