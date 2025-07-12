CREATE TABLE IF NOT EXISTS tab_user (
    user_id INT NOT NULL AUTO_INCREMENT,
    user_uuid CHAR(36) NOT NULL,
    email VARCHAR(255) NOT NULL,
    name VARCHAR(450) NOT NULL,
    password VARCHAR(200) NOT NULL,
    phone VARCHAR(20) NULL,
    profile_photo VARCHAR(500) NULL,
    plan VARCHAR(50) NULL DEFAULT 'trial',
    trial_ends_at TIMESTAMP NULL,
    subscribed_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    last_login_at TIMESTAMP NULL,
    active TINYINT(1) NOT NULL DEFAULT 1,
    email_verified TINYINT(1) NOT NULL DEFAULT 0,
    
    PRIMARY KEY (user_id),
    UNIQUE INDEX user_id_UNIQUE (user_id ASC) VISIBLE,
    UNIQUE INDEX email_UNIQUE (email ASC) VISIBLE,
    INDEX idx_user_active (active ASC) VISIBLE
) ENGINE = InnoDB CHARACTER SET=utf8mb4;


CREATE TABLE IF NOT EXISTS tab_session (
    session_id INT NOT NULL AUTO_INCREMENT,
    session_uuid CHAR(36) NOT NULL,
    user_id INT NULL,
    refresh_token VARCHAR(1500) NOT NULL,
    user_agent VARCHAR(1000) NOT NULL,
    client_ip VARCHAR(500) NOT NULL,
    is_blocked TINYINT(1) NOT NULL DEFAULT 0,
    refresh_token_expires_at TIMESTAMP NOT NULL ,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (session_id),
    UNIQUE INDEX session_id_UNIQUE (session_id ASC) VISIBLE,
    INDEX fk_tab_session_tab_user_idx (user_id ASC) VISIBLE,

    CONSTRAINT fk_tab_session_tab_user
        FOREIGN KEY (user_id)
        REFERENCES tab_user (user_id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION
) ENGINE = InnoDB CHARACTER SET=utf8;

CREATE TABLE IF NOT EXISTS tab_company (
    company_id INT NOT NULL AUTO_INCREMENT,
    company_uuid CHAR(36) NOT NULL,
    name VARCHAR(255) NOT NULL,
    industry VARCHAR(100) NULL,
    size VARCHAR(50) NULL,
    role VARCHAR(200) NULL,
    is_default TINYINT(1) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    user_owner_id INT NOT NULL,
    active TINYINT(1) NOT NULL DEFAULT 1,
    
    PRIMARY KEY (company_id),
    UNIQUE INDEX company_id_UNIQUE (company_id ASC) VISIBLE,
    INDEX idx_company_active (active ASC) VISIBLE,
    
    CONSTRAINT fk_company_owner
        FOREIGN KEY (user_owner_id)
        REFERENCES tab_user (user_id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION
) ENGINE = InnoDB CHARACTER SET=utf8mb4;

CREATE TABLE IF NOT EXISTS tab_person (
    person_id INT NOT NULL AUTO_INCREMENT,
    person_uuid CHAR(36) NOT NULL,
    company_id INT NOT NULL,
    name VARCHAR(450) NOT NULL,
    email VARCHAR(255) NULL,
    position VARCHAR(200) NULL,
    department VARCHAR(200) NULL,
    phone VARCHAR(20) NULL,
    birthday DATE NULL,
    start_date DATE NULL,
    is_manager TINYINT(1) NOT NULL DEFAULT 0,
    manager_id INT NULL,
    notes TEXT NULL,
    has_kids TINYINT(1) NOT NULL DEFAULT 0,
    interests TEXT NULL,
    personality TEXT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by INT NOT NULL,
    active TINYINT(1) NOT NULL DEFAULT 1,
    
    PRIMARY KEY (person_id),
    UNIQUE INDEX person_id_UNIQUE (person_id ASC) VISIBLE,
    INDEX idx_person_company (company_id ASC) VISIBLE,
    INDEX idx_person_manager (manager_id ASC) VISIBLE,
    INDEX idx_person_active (active ASC) VISIBLE,
    
    CONSTRAINT fk_person_company
        FOREIGN KEY (company_id)
        REFERENCES tab_company (company_id)
        ON DELETE CASCADE
        ON UPDATE NO ACTION,
        
    CONSTRAINT fk_person_manager
        FOREIGN KEY (manager_id)
        REFERENCES tab_person (person_id)
        ON DELETE SET NULL
        ON UPDATE NO ACTION,
        
    CONSTRAINT fk_person_created_by
        FOREIGN KEY (created_by)
        REFERENCES tab_user (user_id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION
) ENGINE = InnoDB CHARACTER SET=utf8mb4;