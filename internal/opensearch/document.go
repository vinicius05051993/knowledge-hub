package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func IndexDocument(
	ctx context.Context,
	client *Client,
	document *Document,
) error {

	documentID := fmt.Sprintf(
		"%s:%s",
		document.Namespace,
		document.ExternalID,
	)

	payload, err := json.Marshal(
		document,
	)

	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		client.URL(
			fmt.Sprintf(
				"%s/_doc/%s?refresh=true",
				DocumentsIndex,
				documentID,
			),
		),
		bytes.NewBuffer(payload),
	)

	if err != nil {
		return err
	}

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {

		content, _ := io.ReadAll(
			resp.Body,
		)

		return fmt.Errorf(
			"opensearch error: %s",
			string(content),
		)
	}

	return nil
}

func DeleteDocument(
	ctx context.Context,
	client *Client,
	namespace string,
	externalID string,
) error {

	documentID := fmt.Sprintf(
		"%s:%s",
		namespace,
		externalID,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		client.URL(
			fmt.Sprintf(
				"%s/_doc/%s?refresh=true",
				DocumentsIndex,
				documentID,
			),
		),
		nil,
	)

	if err != nil {
		return err
	}

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {

		content, _ := io.ReadAll(
			resp.Body,
		)

		return fmt.Errorf(
			"opensearch error: %s",
			string(content),
		)
	}

	return nil
}