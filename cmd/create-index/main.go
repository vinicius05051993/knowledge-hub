package main

import (
	"context"
	"log"

	"indexer/internal/config"
	"indexer/internal/opensearch"
)

func main() {

	cfg, err := config.Load()

	if err != nil {
		log.Fatal(err)
	}

	client :=
		opensearch.NewClient(cfg)

	err = opensearch.CreateDocumentsIndex(
		context.Background(),
		client,
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(
		"documents index created",
	)
}