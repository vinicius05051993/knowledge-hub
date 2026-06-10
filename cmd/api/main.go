package main

import (
	"fmt"
	"log"
	"net/http"

	"indexer/internal/apikeys"
	"indexer/internal/config"
	"indexer/internal/database"
	"indexer/internal/documents"
	"indexer/internal/opensearch"
	"indexer/internal/server"
	"indexer/internal/documentfilters"
)

func main() {

	cfg := config.Load()

	db, err := database.NewMySQL(cfg)

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

	fmt.Println(
		"Server running on",
		address,
	)

	err = http.ListenAndServe(
		address,
		app.Router(),
	)

	if err != nil {
		log.Fatal(err)
	}
}