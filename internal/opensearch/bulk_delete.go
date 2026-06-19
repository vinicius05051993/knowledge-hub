package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func DeleteDocuments(
	ctx context.Context,
	client *Client,
	documentKeys []string,
) error {

	if len(documentKeys) == 0 {
		return nil
	}

	path := fmt.Sprintf(
		"%s/_delete_by_query",
		DocumentsIndex,
	)

	namespace, _, found :=
		strings.Cut(
			documentKeys[0],
			":",
		)

	if found && strings.Contains(namespace, "test") {
		path += "?refresh=true"
	}

	body := map[string]any{
		"query": map[string]any{
			"terms": map[string]any{
				"document_key": documentKeys,
			},
		},
	}

	payload, err := json.Marshal(
		body,
	)

	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		client.URL(path),
		bytes.NewReader(payload),
	)

	if err != nil {
		return err
	}

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := client.httpClient.Do(
		req,
	)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {

		content, _ := io.ReadAll(
			resp.Body,
		)

		return fmt.Errorf(
			"opensearch delete_by_query error: %s",
			string(content),
		)
	}

	var result struct {
		Deleted int64 `json:"deleted"`
		Failures []any `json:"failures"`
	}

	err = json.NewDecoder(
		resp.Body,
	).Decode(
		&result,
	)

	if err != nil {
		return err
	}

	if len(result.Failures) > 0 {

		return fmt.Errorf(
			"opensearch delete_by_query returned %d failures",
			len(result.Failures),
		)
	}

	return nil
}