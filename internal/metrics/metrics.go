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

	SearchDurationSeconds =
		prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Name: "indexer_search_duration_seconds",
				Help: "Search duration in seconds",
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

	SyncUpsertsTotal =
		prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "indexer_sync_upserts_total",
				Help: "Total documents synchronized to OpenSearch",
			},
		)

	SyncDeletesTotal =
		prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "indexer_sync_deletes_total",
				Help: "Total documents deleted from OpenSearch",
			},
		)
)

func Register() {

	prometheus.MustRegister(
		SearchRequestsTotal,
	)

	prometheus.MustRegister(
		SearchDurationSeconds,
	)

	prometheus.MustRegister(
		UpsertDocumentsTotal,
	)

	prometheus.MustRegister(
		DeleteDocumentsTotal,
	)

	prometheus.MustRegister(
		SyncUpsertsTotal,
	)

	prometheus.MustRegister(
		SyncDeletesTotal,
	)
}