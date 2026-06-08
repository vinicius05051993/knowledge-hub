CREATE TABLE documents (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,

    namespace VARCHAR(100) NOT NULL,
    external_id VARCHAR(255) NOT NULL,

    title VARCHAR(500) NOT NULL,
    text LONGTEXT NOT NULL,

    payload JSON,

    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,

    UNIQUE KEY uk_namespace_external_id (
        namespace,
        external_id
    ),

    KEY idx_namespace (
        namespace
    )
);