package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"indexer/internal/config"
	"indexer/internal/database"
	"indexer/internal/documents"
	"indexer/internal/metrics"
	"indexer/internal/opensearch"
	"indexer/internal/server"
)

const (
	batchSize    = 100
	syncInterval = 30 * time.Second
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

	stop := make(
		chan os.Signal,
		1,
	)

	signal.Notify(
		stop,
		os.Interrupt,
		syscall.SIGTERM,
	)

	ticker := time.NewTicker(
		syncInterval,
	)

	defer ticker.Stop()

	for {

		select {

		case <-stop:

			log.Println(
				"shutdown signal received",
			)

			log.Println(
				"sync worker stopped",
			)

			return

		case <-ticker.C:

			runCycle(
				syncService,
			)
		}
	}
}

func runCycle(
	syncService *documents.SyncService,
) {

	defer func() {

		if r := recover(); r != nil {

			log.Printf(
				"worker panic recovered: %v",
				r,
			)
		}
	}()

	ctx := context.Background()

	err := syncService.ProcessPendingUpserts(
		ctx,
		batchSize,
	)

	if err != nil {

		log.Printf(
			"upsert sync error: %v",
			err,
		)
	}

	err = syncService.ProcessPendingDeletes(
		ctx,
		batchSize,
	)

	if err != nil {

		log.Printf(
			"delete sync error: %v",
			err,
		)
	}
}