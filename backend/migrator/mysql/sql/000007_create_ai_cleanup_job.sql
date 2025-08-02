-- ================================================
-- Migration 000007: AI Cleanup Job
-- ================================================

-- Note: This migration creates a stored procedure and event for cleanup
-- The DELIMITER syntax doesn't work well in migration tools, so we'll create it directly

-- Drop procedure if exists to avoid conflicts
DROP PROCEDURE IF EXISTS `cleanup_ai_conversations`;

-- Create procedure without DELIMITER (single statement)
CREATE PROCEDURE `cleanup_ai_conversations`()
DELETE FROM `ai_conversations` 
WHERE `expires_at` < NOW() 
LIMIT 1000;

-- Drop event if exists to avoid conflicts  
DROP EVENT IF EXISTS `cleanup_old_ai_data`;

-- Create event scheduler to run daily at 2:00 AM
CREATE EVENT `cleanup_old_ai_data`
ON SCHEDULE EVERY 1 DAY
STARTS TIMESTAMP(CURRENT_DATE + INTERVAL 1 DAY, '02:00:00')
DO CALL `cleanup_ai_conversations`();