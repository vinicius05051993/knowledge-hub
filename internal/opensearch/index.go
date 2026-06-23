package opensearch

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
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

	url := client.URL(
		DocumentsIndex,
	)

	log.Printf(
		"CreateDocumentsIndex url=%s",
		url,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		url,
		bytes.NewBufferString(body),
	)

	if err != nil {

		log.Printf(
			"CreateDocumentsIndex request creation failed: %v",
			err,
		)

		return err
	}

	resp, err := client.Do(
		req,
	)

	if err != nil {

		log.Printf(
			"CreateDocumentsIndex request failed url=%s error=%v",
			url,
			err,
		)

		return fmt.Errorf(
			"create documents index request failed: %w",
			err,
		)
	}

	defer resp.Body.Close()

	content, _ := io.ReadAll(
		resp.Body,
	)

	log.Printf(
		"CreateDocumentsIndex response status=%d body=%s",
		resp.StatusCode,
		string(content),
	)

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	if resp.StatusCode == http.StatusBadRequest {

		if bytes.Contains(
			content,
			[]byte("resource_already_exists_exception"),
		) {

			log.Printf(
				"CreateDocumentsIndex index already exists",
			)

			return nil
		}
	}

	return fmt.Errorf(
		"opensearch create index failed status=%d body=%s",
		resp.StatusCode,
		string(content),
	)
}