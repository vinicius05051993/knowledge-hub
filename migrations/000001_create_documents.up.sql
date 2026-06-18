CREATE TABLE documents (

    id BIGINT IDENTITY(1,1) PRIMARY KEY,

    document_key VARCHAR(400) NOT NULL,

    namespace VARCHAR(100) NOT NULL,
    external_id VARCHAR(255) NOT NULL,

    title VARCHAR(500) NOT NULL,
    [text] NVARCHAR(MAX) NOT NULL,

    sync_status TINYINT NOT NULL DEFAULT 0,
    deleted_at DATETIME NULL,

    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,

    CONSTRAINT uk_document_key
        UNIQUE (document_key),

    CONSTRAINT uk_namespace_external_id
        UNIQUE (
            namespace,
            external_id
        )
);