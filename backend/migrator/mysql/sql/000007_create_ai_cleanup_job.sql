-- ================================================
-- Migration 000007: AI Cleanup Job
-- ================================================

-- Procedure for automatic cleanup of old conversations
DELIMITER //

CREATE PROCEDURE IF NOT EXISTS `cleanup_ai_conversations`()
BEGIN
    DELETE FROM `ai_conversations` 
    WHERE `expires_at` < NOW() 
    LIMIT 1000; -- Process in batches to avoid blocking database
END//

DELIMITER ;

-- Event scheduler to run daily at 2:00 AM
CREATE EVENT IF NOT EXISTS `cleanup_old_ai_data`
ON SCHEDULE EVERY 1 DAY
STARTS TIMESTAMP(CURRENT_DATE + INTERVAL 1 DAY, '02:00:00')
DO CALL `cleanup_ai_conversations`();

-- Enable event scheduler if not enabled
SET GLOBAL event_scheduler = ON;