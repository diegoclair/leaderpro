CREATE TABLE IF NOT EXISTS tab_one_on_one (
    one_on_one_id INT NOT NULL AUTO_INCREMENT,
    one_on_one_uuid CHAR(36) NOT NULL,
    company_id INT NOT NULL,
    person_id INT NOT NULL,
    manager_id INT NOT NULL,
    scheduled_date TIMESTAMP NOT NULL,
    actual_date TIMESTAMP NULL,
    duration INT NULL DEFAULT 30,
    location VARCHAR(255) NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'scheduled',
    agenda TEXT NULL,
    discussion_notes TEXT NULL,
    action_items TEXT NULL,
    private_notes TEXT NULL,
    ai_context TEXT NULL,
    ai_suggestions TEXT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    completed_at TIMESTAMP NULL,
    
    PRIMARY KEY (one_on_one_id),
    UNIQUE INDEX one_on_one_id_UNIQUE (one_on_one_id ASC) VISIBLE,
    INDEX idx_one_on_one_company (company_id ASC) VISIBLE,
    INDEX idx_one_on_one_person (person_id ASC) VISIBLE,
    INDEX idx_one_on_one_manager (manager_id ASC) VISIBLE,
    INDEX idx_one_on_one_status (status ASC) VISIBLE,
    INDEX idx_one_on_one_scheduled (scheduled_date ASC) VISIBLE,
    
    CONSTRAINT fk_one_on_one_company
        FOREIGN KEY (company_id)
        REFERENCES tab_company (company_id)
        ON DELETE CASCADE
        ON UPDATE NO ACTION,
        
    CONSTRAINT fk_one_on_one_person
        FOREIGN KEY (person_id)
        REFERENCES tab_person (person_id)
        ON DELETE CASCADE
        ON UPDATE NO ACTION,
        
    CONSTRAINT fk_one_on_one_manager
        FOREIGN KEY (manager_id)
        REFERENCES tab_user (user_id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION
) ENGINE = InnoDB CHARACTER SET=utf8mb4;

CREATE TABLE IF NOT EXISTS tab_feedback (
    feedback_id INT NOT NULL AUTO_INCREMENT,
    feedback_uuid CHAR(36) NOT NULL,
    company_id INT NOT NULL,
    person_id INT NOT NULL,
    given_by INT NOT NULL,
    one_on_one_id INT NULL,
    type VARCHAR(50) NOT NULL,
    category VARCHAR(50) NOT NULL,
    content TEXT NOT NULL,
    mentioned_from VARCHAR(255) NULL,
    mentioned_date TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    is_private TINYINT(1) NOT NULL DEFAULT 0,
    
    PRIMARY KEY (feedback_id),
    UNIQUE INDEX feedback_id_UNIQUE (feedback_id ASC) VISIBLE,
    INDEX idx_feedback_company (company_id ASC) VISIBLE,
    INDEX idx_feedback_person (person_id ASC) VISIBLE,
    INDEX idx_feedback_given_by (given_by ASC) VISIBLE,
    INDEX idx_feedback_one_on_one (one_on_one_id ASC) VISIBLE,
    INDEX idx_feedback_type (type ASC) VISIBLE,
    
    CONSTRAINT fk_feedback_company
        FOREIGN KEY (company_id)
        REFERENCES tab_company (company_id)
        ON DELETE CASCADE
        ON UPDATE NO ACTION,
        
    CONSTRAINT fk_feedback_person
        FOREIGN KEY (person_id)
        REFERENCES tab_person (person_id)
        ON DELETE CASCADE
        ON UPDATE NO ACTION,
        
    CONSTRAINT fk_feedback_given_by
        FOREIGN KEY (given_by)
        REFERENCES tab_user (user_id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION,
        
    CONSTRAINT fk_feedback_one_on_one
        FOREIGN KEY (one_on_one_id)
        REFERENCES tab_one_on_one (one_on_one_id)
        ON DELETE SET NULL
        ON UPDATE NO ACTION
) ENGINE = InnoDB CHARACTER SET=utf8mb4;