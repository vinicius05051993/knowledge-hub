CREATE TABLE document_filters (

    id BIGINT AUTO_INCREMENT PRIMARY KEY,

    document_key VARCHAR(255) NOT NULL,

    field_name VARCHAR(100) NOT NULL,

    field_value VARCHAR(255) NOT NULL,

    INDEX idx_document_key (
        document_key
    ),

    INDEX idx_name_value (
        field_name,
        field_value
    )
);