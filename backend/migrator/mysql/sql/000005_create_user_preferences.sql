CREATE TABLE IF NOT EXISTS `user_preferences` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `user_id` INT NOT NULL,
    
    -- Appearance
    `theme` VARCHAR(20) NOT NULL DEFAULT 'light',
    
    -- Metadata
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_user_id` (`user_id`),
    CONSTRAINT `fk_user_preferences_user` FOREIGN KEY (`user_id`) REFERENCES `tab_user` (`user_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
