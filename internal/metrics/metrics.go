package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	SearchRequestsTotal =
		prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "indexer_search_requests_total",
				Help: "Total search requests",
			},
		)

	UpsertDocumentsTotal =
		prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "indexer_upsert_documents_total",
				Help: "Total documents upserted",
			},
		)

	DeleteDocumentsTotal =
		prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "indexer_delete_documents_total",
				Help: "Total documents deleted",
			},
		)
)

func Register() {

	prometheus.MustRegister(
		SearchRequestsTotal,
	)

	prometheus.MustRegister(
		UpsertDocumentsTotal,
	)

	prometheus.MustRegister(
		DeleteDocumentsTotal,
	)
}