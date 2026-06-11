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

	path := "_bulk"

	namespace, _, found :=
		strings.Cut(
			documentKeys[0],
			":",
		)

	if found && strings.Contains(namespace,"test") {
		path += "?refresh=true"
	}

	var body bytes.Buffer

	for _, documentKey := range documentKeys {

		action := map[string]any{
			"delete": map[string]any{
				"_index": DocumentsIndex,
				"_id":    documentKey,
			},
		}

		data, err := json.Marshal(
			action,
		)

		if err != nil {
			return err
		}

		body.Write(data)
		body.WriteByte('\n')
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		client.URL(path),
		&body,
	)

	if err != nil {
		return err
	}

	req.Header.Set(
		"Content-Type",
		"application/x-ndjson",
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
			"opensearch bulk delete error: %s",
			string(content),
		)
	}

	var result struct {
		Errors bool `json:"errors"`
	}

	err = json.NewDecoder(
		resp.Body,
	).Decode(
		&result,
	)

	if err != nil {
		return err
	}

	if result.Errors {

		content, _ := io.ReadAll(
			resp.Body,
		)

		return fmt.Errorf(
			"opensearch bulk delete returned errors: %s",
			string(content),
		)
	}

	return nil
}