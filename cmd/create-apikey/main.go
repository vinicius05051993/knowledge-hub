package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"indexer/internal/apikeys"
	"indexer/internal/config"
	"indexer/internal/database"
)

func main() {

	name := flag.String(
		"name",
		"",
		"API key name",
	)

	namespace := flag.String(
		"namespace",
		"",
		"Namespace",
	)

	flag.Parse()

	if *name == "" {
		log.Fatal("name is required")
	}

	if *namespace == "" {
		log.Fatal("namespace is required")
	}

	cfg, err := config.Load()

	if err != nil {
		log.Fatal(err)
	}

	db, err := database.New(*cfg)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	repository := apikeys.NewRepository(
		db,
	)

	service := apikeys.NewService(
		repository,
	)

	apiKey, err := service.Create(
		context.Background(),
		*name,
		*namespace,
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()
	fmt.Println("API Key created")
	fmt.Println()

	fmt.Println("Name:", *name)
	fmt.Println("Namespace:", *namespace)

	fmt.Println()
	fmt.Println("API Key:")
	fmt.Println(apiKey)
	fmt.Println()
}