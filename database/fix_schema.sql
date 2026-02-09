-- Fix Schema: Add missing columns to existing unlocked_users table
-- Run this if you get column errors

-- Add verified column if missing
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'unlocked_users' AND column_name = 'verified'
    ) THEN
        ALTER TABLE unlocked_users ADD COLUMN verified BOOLEAN DEFAULT false;
    END IF;
END $$;

-- Add discord_webhook_url if missing
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'unlocked_users' AND column_name = 'discord_webhook_url'
    ) THEN
        ALTER TABLE unlocked_users ADD COLUMN discord_webhook_url TEXT;
    END IF;
END $$;

-- Add discord_notifications if missing
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'unlocked_users' AND column_name = 'discord_notifications'
    ) THEN
        ALTER TABLE unlocked_users ADD COLUMN discord_notifications JSONB DEFAULT '[]'::jsonb;
    END IF;
END $$;

-- Add notifications_subscription_active if missing
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'unlocked_users' AND column_name = 'notifications_subscription_active'
    ) THEN
        ALTER TABLE unlocked_users ADD COLUMN notifications_subscription_active BOOLEAN DEFAULT false;
    END IF;
END $$;

-- Add notifications_subscription_expires_at if missing
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'unlocked_users' AND column_name = 'notifications_subscription_expires_at'
    ) THEN
        ALTER TABLE unlocked_users ADD COLUMN notifications_subscription_expires_at TIMESTAMPTZ;
    END IF;
END $$;

-- Add last_active if missing
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'unlocked_users' AND column_name = 'last_active'
    ) THEN
        ALTER TABLE unlocked_users ADD COLUMN last_active TIMESTAMPTZ DEFAULT NOW();
    END IF;
END $$;

-- Add created_at if missing
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'unlocked_users' AND column_name = 'created_at'
    ) THEN
        ALTER TABLE unlocked_users ADD COLUMN created_at TIMESTAMPTZ DEFAULT NOW();
    END IF;
END $$;

-- Add updated_at if missing
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'unlocked_users' AND column_name = 'updated_at'
    ) THEN
        ALTER TABLE unlocked_users ADD COLUMN updated_at TIMESTAMPTZ DEFAULT NOW();
    END IF;
END $$;

-- Create indexes if they don't exist
CREATE INDEX IF NOT EXISTS idx_unlocked_users_auth_user_id ON unlocked_users(auth_user_id);
CREATE INDEX IF NOT EXISTS idx_unlocked_users_verified ON unlocked_users(verified);
CREATE INDEX IF NOT EXISTS idx_unlocked_users_subscription_active ON unlocked_users(notifications_subscription_active);

-- Update NULL values to defaults
UPDATE unlocked_users SET verified = false WHERE verified IS NULL;
UPDATE unlocked_users SET notifications_subscription_active = false WHERE notifications_subscription_active IS NULL;
UPDATE unlocked_users SET discord_notifications = '[]'::jsonb WHERE discord_notifications IS NULL;
UPDATE unlocked_users SET created_at = NOW() WHERE created_at IS NULL;
UPDATE unlocked_users SET updated_at = NOW() WHERE updated_at IS NULL;
