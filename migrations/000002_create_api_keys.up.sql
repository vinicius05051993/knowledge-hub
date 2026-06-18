CREATE TABLE api_keys (

    id BIGINT IDENTITY(1,1) PRIMARY KEY,

    name VARCHAR(255) NOT NULL,

    namespace VARCHAR(100) NOT NULL,

    api_key_hash CHAR(64) NOT NULL,

    permissions NVARCHAR(MAX) NULL,

    active BIT NOT NULL DEFAULT 1,

    created_at DATETIME NOT NULL,

    CONSTRAINT uk_api_key_hash
        UNIQUE (
            api_key_hash
        ),

    CONSTRAINT uk_namespace
        UNIQUE (
            namespace
        )
);