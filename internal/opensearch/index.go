package opensearch

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

const DocumentsIndex = "documents"

func CreateDocumentsIndex(
	ctx context.Context,
	client *Client,
) error {

	body := `
{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 0
  },
  "mappings": {
    "properties": {
      "namespace": {
        "type": "keyword"
      },
      "external_id": {
        "type": "keyword"
      },
      "title": {
        "type": "text"
      },
      "text": {
        "type": "text"
      }
    }
  }
}
`

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		client.URL(DocumentsIndex),
		bytes.NewBufferString(body),
	)

	if err != nil {
		return err
	}

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	content, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusBadRequest {

		if bytes.Contains(
			content,
			[]byte("resource_already_exists_exception"),
		) {
			return nil
		}
	}

	return fmt.Errorf(
		"opensearch error: %s",
		string(content),
	)
}

func IndexExists(
	ctx context.Context,
	client *Client,
	index string,
) (bool, error) {

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodHead,
		client.URL(index),
		nil,
	)

	if err != nil {
		return false, err
	}

	resp, err := client.Do(req)

	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	return false,
		fmt.Errorf(
			"unexpected status %d",
			resp.StatusCode,
		)
}