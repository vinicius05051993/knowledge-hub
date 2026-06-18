package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"indexer/internal/apikeys"
	"indexer/internal/config"
	"indexer/internal/database"
	"indexer/internal/documentfilters"
	"indexer/internal/documents"
	"indexer/internal/metrics"
	"indexer/internal/opensearch"
	"indexer/internal/server"
)

func main() {

	metrics.Register()

	cfg, err := config.Load()

	if err != nil {
		log.Fatal(err)
	}

	db, err := database.New(*cfg)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	apiKeyRepository :=
		apikeys.NewRepository(db)

	apiKeyService :=
		apikeys.NewService(
			apiKeyRepository,
		)

	searchClient :=
		opensearch.NewClient(cfg)

	searchService :=
		opensearch.NewService(
			searchClient,
		)

	documentRepository :=
		documents.NewRepository(db)

	filterRepository :=
		documentfilters.NewRepository(
			db,
		)

	documentService :=
		documents.NewService(
			documentRepository,
			filterRepository,
			searchService,
		)

	documentHandler :=
		documents.NewHandler(
			documentService,
		)

	app := server.NewApp(
		db,
		apiKeyService,
		documentHandler,
	)

	address := ":" + cfg.AppPort

	srv := &http.Server{
		Addr:    address,
		Handler: app.Router(),
	}

	go func() {

		log.Printf(
			"server running on %s",
			address,
		)

		err := srv.ListenAndServe()

		if err != nil &&
			err != http.ErrServerClosed {

			log.Fatal(err)
		}
	}()

	stop := make(
		chan os.Signal,
		1,
	)

	signal.Notify(
		stop,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-stop

	log.Println(
		"shutdown signal received",
	)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		30*time.Second,
	)

	defer cancel()

	err = srv.Shutdown(
		ctx,
	)

	if err != nil {

		log.Printf(
			"server shutdown error: %v",
			err,
		)

		return
	}

	log.Println(
		"server stopped",
	)
}