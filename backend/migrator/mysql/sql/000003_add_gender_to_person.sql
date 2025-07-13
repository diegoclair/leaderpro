-- Add gender field to tab_person table
ALTER TABLE tab_person 
ADD COLUMN gender ENUM('male', 'female', 'other') NULL AFTER has_kids;