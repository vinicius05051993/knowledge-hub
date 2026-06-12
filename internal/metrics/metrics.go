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

	SearchErrorsTotal =
		prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "indexer_search_errors_total",
				Help: "Total search errors",
			},
		)

	UpsertDocumentsTotal =
		prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "indexer_upsert_documents_total",
				Help: "Total documents upserted",
			},
		)

	UpsertErrorsTotal =
		prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "indexer_upsert_errors_total",
				Help: "Total upsert errors",
			},
		)

	DeleteDocumentsTotal =
		prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "indexer_delete_documents_total",
				Help: "Total documents deleted",
			},
		)

	DeleteErrorsTotal =
		prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "indexer_delete_errors_total",
				Help: "Total delete errors",
			},
		)

	SyncUpsertsTotal =
		prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "indexer_sync_upserts_total",
				Help: "Total documents synchronized to OpenSearch",
			},
		)

	SyncUpsertErrorsTotal =
		prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "indexer_sync_upsert_errors_total",
				Help: "Total sync upsert errors",
			},
		)

	SyncDeletesTotal =
		prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "indexer_sync_deletes_total",
				Help: "Total documents deleted from OpenSearch",
			},
		)

	SyncDeleteErrorsTotal =
		prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "indexer_sync_delete_errors_total",
				Help: "Total sync delete errors",
			},
		)

	PendingUpserts =
		prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "indexer_pending_upserts",
				Help: "Current pending upserts",
			},
		)

	PendingDeletes =
		prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "indexer_pending_deletes",
				Help: "Current pending deletes",
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
		SearchErrorsTotal,
	)

	prometheus.MustRegister(
		UpsertDocumentsTotal,
	)

	prometheus.MustRegister(
		UpsertErrorsTotal,
	)

	prometheus.MustRegister(
		DeleteDocumentsTotal,
	)

	prometheus.MustRegister(
		DeleteErrorsTotal,
	)

	prometheus.MustRegister(
		SyncUpsertsTotal,
	)

	prometheus.MustRegister(
		SyncUpsertErrorsTotal,
	)

	prometheus.MustRegister(
		SyncDeletesTotal,
	)

	prometheus.MustRegister(
		SyncDeleteErrorsTotal,
	)

	prometheus.MustRegister(
		PendingUpserts,
	)

	prometheus.MustRegister(
		PendingDeletes,
	)
}