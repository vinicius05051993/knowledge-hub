CREATE TABLE document_filters (

    id BIGINT IDENTITY(1,1) PRIMARY KEY,

    document_key VARCHAR(255) NOT NULL,

    field_name VARCHAR(100) NOT NULL,

    field_value VARCHAR(255) NOT NULL
);

CREATE INDEX idx_document_key
ON document_filters(document_key);

CREATE INDEX idx_name_value
ON document_filters(
    field_name,
    field_value
);