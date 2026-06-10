package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type searchResponse struct {
	Hits struct {
		Hits []struct {
			Score float64 `json:"_score"`

			Source struct {
				DocumentKey string `json:"document_key"`

				Namespace string `json:"namespace"`

				ExternalID string `json:"external_id"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func Search(
	ctx context.Context,
	client *Client,
	query string,
	offset int,
	limit int,
) ([]SearchResult, error) {

	body := map[string]interface{}{
		"from": offset,
		"size": limit,
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query": query,
				"fields": []string{
					"title^2",
					"text",
				},
			},
		},
	}

	payload, err := json.Marshal(
		body,
	)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		client.URL(
			fmt.Sprintf(
				"%s/_search",
				DocumentsIndex,
			),
		),
		bytes.NewBuffer(payload),
	)

	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {

		content, _ := io.ReadAll(
			resp.Body,
		)

		return nil, fmt.Errorf(
			"opensearch error: %s",
			string(content),
		)
	}

	var result searchResponse

	err = json.NewDecoder(
		resp.Body,
	).Decode(
		&result,
	)

	if err != nil {
		return nil, err
	}

	searchResults :=
		make(
			[]SearchResult,
			0,
			len(result.Hits.Hits),
		)

	for _, hit := range result.Hits.Hits {

		searchResults = append(
			searchResults,
			SearchResult{
				DocumentKey: hit.Source.DocumentKey,

				Namespace: hit.Source.Namespace,

				ExternalID: hit.Source.ExternalID,

				Score: hit.Score,
			},
		)
	}

	return searchResults, nil
}