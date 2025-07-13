CREATE TABLE IF NOT EXISTS tab_note (
    note_id INT NOT NULL AUTO_INCREMENT,
    note_uuid CHAR(36) NOT NULL,
    company_id INT NOT NULL,
    person_id INT NOT NULL,
    user_id INT NOT NULL,
    type VARCHAR(50) NOT NULL,
    content TEXT NOT NULL,
    feedback_type VARCHAR(50) NULL,
    feedback_category VARCHAR(50) NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    PRIMARY KEY (note_id),
    UNIQUE INDEX note_id_UNIQUE (note_id ASC) VISIBLE,
    UNIQUE INDEX note_uuid_UNIQUE (note_uuid ASC) VISIBLE,
    INDEX idx_note_company (company_id ASC) VISIBLE,
    INDEX idx_note_person (person_id ASC) VISIBLE,
    INDEX idx_note_user (user_id ASC) VISIBLE,
    INDEX idx_note_type (type ASC) VISIBLE,
    INDEX idx_note_created (created_at ASC) VISIBLE,
    INDEX idx_note_deleted (deleted_at ASC) VISIBLE,
    
    CONSTRAINT fk_note_company
        FOREIGN KEY (company_id)
        REFERENCES tab_company (company_id)
        ON DELETE CASCADE
        ON UPDATE NO ACTION,
        
    CONSTRAINT fk_note_person
        FOREIGN KEY (person_id)
        REFERENCES tab_person (person_id)
        ON DELETE CASCADE
        ON UPDATE NO ACTION,
        
    CONSTRAINT fk_note_user
        FOREIGN KEY (user_id)
        REFERENCES tab_user (user_id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION
) ENGINE = InnoDB CHARACTER SET=utf8mb4;

CREATE TABLE IF NOT EXISTS tab_note_mention (
    mention_id INT NOT NULL AUTO_INCREMENT,
    mention_uuid CHAR(36) NOT NULL,
    note_id INT NOT NULL,
    mentioned_person_id INT NOT NULL,
    source_person_id INT NOT NULL,
    full_content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (mention_id),
    UNIQUE INDEX mention_id_UNIQUE (mention_id ASC) VISIBLE,
    UNIQUE INDEX mention_uuid_UNIQUE (mention_uuid ASC) VISIBLE,
    INDEX idx_mention_note (note_id ASC) VISIBLE,
    INDEX idx_mention_mentioned_person (mentioned_person_id ASC) VISIBLE,
    INDEX idx_mention_source_person (source_person_id ASC) VISIBLE,
    INDEX idx_mention_created (created_at ASC) VISIBLE,
    
    CONSTRAINT fk_mention_note
        FOREIGN KEY (note_id)
        REFERENCES tab_note (note_id)
        ON DELETE CASCADE
        ON UPDATE NO ACTION,
        
    CONSTRAINT fk_mention_mentioned_person
        FOREIGN KEY (mentioned_person_id)
        REFERENCES tab_person (person_id)
        ON DELETE CASCADE
        ON UPDATE NO ACTION,
        
    CONSTRAINT fk_mention_source_person
        FOREIGN KEY (source_person_id)
        REFERENCES tab_person (person_id)
        ON DELETE CASCADE
        ON UPDATE NO ACTION
) ENGINE = InnoDB CHARACTER SET=utf8mb4;