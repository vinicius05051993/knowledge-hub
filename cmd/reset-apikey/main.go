package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"indexer/internal/apikeys"
	"indexer/internal/config"
	"indexer/internal/database"
)

func main() {

	if len(os.Args) != 2 {
		log.Fatalf(
			"usage: %s <namespace>",
			os.Args[0],
		)
	}

	namespace := os.Args[1]

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

	ctx, cancel := context.WithTimeout(
		context.Background(),
		30*time.Second,
	)

	defer cancel()

	apiKey, err := service.ResetByNamespace(
		ctx,
		namespace,
	)

	if err != nil {
		log.Fatalf(
			"failed to reset api key: %v",
			err,
		)
	}

	fmt.Println("API Key reset successfully")
	fmt.Println()

	fmt.Printf(
		"Namespace: %s\n",
		namespace,
	)

	fmt.Printf(
		"API Key: %s\n",
		apiKey,
	)
}