-- Migration: Add verified column if it doesn't exist
-- Run this if you get "column verified does not exist" error

-- Add verified column if it doesn't exist
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_name = 'unlocked_users' 
        AND column_name = 'verified'
    ) THEN
        ALTER TABLE unlocked_users ADD COLUMN verified BOOLEAN DEFAULT false;
    END IF;
END $$;

-- Create index if it doesn't exist
CREATE INDEX IF NOT EXISTS idx_unlocked_users_verified ON unlocked_users(verified);

-- Update existing rows to have verified = false if NULL
UPDATE unlocked_users SET verified = false WHERE verified IS NULL;
