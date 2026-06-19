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

func BulkIndexDocuments(
	ctx context.Context,
	client *Client,
	documents []*Document,
) error {

	if len(documents) == 0 {
		return nil
	}

	path := "_bulk"

	if strings.Contains(
		documents[0].Namespace,
		"test",
	) {
		path += "?refresh=true"
	}

	var body bytes.Buffer

	for _, document := range documents {

		action := map[string]any{
			"index": map[string]any{
		        "_index": DocumentsIndex,
		        "_id":    document.ID,
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

		data, err = json.Marshal(
			document,
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

	resp, err :=
		client.httpClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {

		content, _ := io.ReadAll(
			resp.Body,
		)

		return fmt.Errorf(
			"opensearch bulk index error: %s",
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

		return fmt.Errorf(
			"opensearch bulk index returned errors",
		)
	}

	return nil
}