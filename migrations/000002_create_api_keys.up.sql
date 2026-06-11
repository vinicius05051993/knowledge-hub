CREATE TABLE api_keys (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,

    name VARCHAR(255) NOT NULL,

    namespace VARCHAR(100) NOT NULL,

    api_key_hash CHAR(64) NOT NULL,

    permissions JSON NULL,

    active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at DATETIME NOT NULL,

    UNIQUE KEY uk_api_key_hash (
        api_key_hash
    ),

    UNIQUE KEY uk_namespace (
        namespace
    )
);