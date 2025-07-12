-- Add is_default column to tab_company_user
ALTER TABLE tab_company_user 
ADD COLUMN is_default TINYINT(1) NOT NULL DEFAULT 0 AFTER role;

-- Add index for faster queries on default companies
CREATE INDEX idx_user_default_company ON tab_company_user (user_id, is_default);

-- Add unique constraint to ensure only one default company per user
ALTER TABLE tab_company_user 
ADD CONSTRAINT uk_user_default_company 
UNIQUE KEY (user_id, is_default) 
COMMENT 'Ensures only one default company per user';