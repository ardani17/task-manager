-- Add password_hash column to developers table for authentication
ALTER TABLE developers ADD COLUMN IF NOT EXISTS password_hash VARCHAR(255);

-- Update status to 'online' when developer logs in (will be managed by auth)
-- No additional indexes needed, email already indexed
