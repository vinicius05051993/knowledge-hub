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

func IndexDocument(
	ctx context.Context,
	client *Client,
	document *Document,
) error {

	payload, err := json.Marshal(
		document,
	)

	if err != nil {
		return err
	}

	documentID := document.ID

	if documentID == "" {
		documentID = document.DocumentKey
	}

	path := fmt.Sprintf(
		"%s/_doc/%s",
		DocumentsIndex,
		documentID,
	)

	if strings.Contains(
		document.Namespace,
		"test",
	) {
		path += "?refresh=true"
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		client.URL(path),
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

	documentKey :=
		namespace +
			":" +
			externalID

	return DeleteDocuments(
		ctx,
		client,
		[]string{
			documentKey,
		},
	)
}