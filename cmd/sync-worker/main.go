package main

import (
	"context"
	"log"
	"time"

	"indexer/internal/config"
	"indexer/internal/database"
	"indexer/internal/documents"
	"indexer/internal/opensearch"
)

func main() {

	cfg := config.Load()

	db, err := database.NewMySQL(
		cfg,
	)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	opensearchClient :=
		opensearch.NewClient(
			cfg,
		)

	opensearchService :=
		opensearch.NewService(
			opensearchClient,
		)

	documentRepository :=
		documents.NewRepository(
			db,
		)

	syncService :=
		documents.NewSyncService(
			documentRepository,
			opensearchService,
		)

	log.Println(
		"sync worker started",
	)

	for {

		ctx := context.Background()

		err := syncService.ProcessPendingUpserts(
			ctx,
			100,
		)

		if err != nil {
			log.Printf(
				"upsert sync error: %v",
				err,
			)
		}

		err = syncService.ProcessPendingDeletes(
			ctx,
			100,
		)

		if err != nil {
			log.Printf(
				"delete sync error: %v",
				err,
			)
		}

		time.Sleep(
			5 * time.Second,
		)
	}
}