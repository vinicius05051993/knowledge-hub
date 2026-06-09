ALTER TABLE documents
ADD COLUMN document_key VARCHAR(400);

UPDATE documents
SET document_key =
	CONCAT(
		namespace,
		':',
		external_id
	);

ALTER TABLE documents
MODIFY document_key
VARCHAR(400)
NOT NULL;

ALTER TABLE documents
ADD UNIQUE KEY uk_document_key (
	document_key
);