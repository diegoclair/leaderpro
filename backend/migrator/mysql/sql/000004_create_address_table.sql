CREATE TABLE IF NOT EXISTS tab_address (
    address_id INT NOT NULL AUTO_INCREMENT,
    address_uuid CHAR(36) NOT NULL,
    person_id INT NOT NULL,
    city VARCHAR(100) NULL COMMENT 'City where person is located',
    state VARCHAR(100) NULL COMMENT 'State/Province where person is located',
    country VARCHAR(100) NULL DEFAULT 'Brazil' COMMENT 'Country where person is located',
    is_primary TINYINT(1) NOT NULL DEFAULT 1 COMMENT 'Whether this is the primary address',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    active TINYINT(1) NOT NULL DEFAULT 1,
    
    PRIMARY KEY (address_id),
    UNIQUE INDEX address_id_UNIQUE (address_id ASC) VISIBLE,
    UNIQUE INDEX address_uuid_UNIQUE (address_uuid ASC) VISIBLE,
    INDEX idx_address_person (person_id ASC) VISIBLE,
    INDEX idx_address_active (active ASC) VISIBLE,
    
    CONSTRAINT fk_address_person
        FOREIGN KEY (person_id)
        REFERENCES tab_person (person_id)
        ON DELETE CASCADE
        ON UPDATE NO ACTION
) ENGINE = InnoDB CHARACTER SET=utf8mb4;