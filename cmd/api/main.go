package main

import (
	"fmt"
	"log"
	"net/http"
	"indexer/internal/apikeys"
	"indexer/internal/config"
	"indexer/internal/database"
	"indexer/internal/server"
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

	router := server.NewRouter(
		db,
		apiKeyService,
	)

	address := ":" + cfg.AppPort

	fmt.Println(
		"Server running on",
		address,
	)

	err = http.ListenAndServe(
		address,
		router,
	)

	if err != nil {
		log.Fatal(err)
	}
}