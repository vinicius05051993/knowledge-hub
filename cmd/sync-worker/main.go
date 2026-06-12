package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"indexer/internal/config"
	"indexer/internal/database"
	"indexer/internal/documents"
	"indexer/internal/metrics"
	"indexer/internal/opensearch"
	"indexer/internal/server"
)

func main() {

	metrics.Register()

	go func() {

		log.Println(
			"worker metrics listening on :9090",
		)

		err := http.ListenAndServe(
			":9090",
			server.MetricsHandler(),
		)

		if err != nil {
			log.Fatal(err)
		}
	}()

	cfg, err := config.Load()

	if err != nil {
		log.Fatal(err)
	}

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
			30 * time.Second,
		)
	}
}