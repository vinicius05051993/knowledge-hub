ALTER TABLE documents
DROP INDEX uk_document_key;

ALTER TABLE documents
DROP COLUMN document_key;