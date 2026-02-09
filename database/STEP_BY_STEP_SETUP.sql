-- ============================================
-- STEP-BY-STEP SETUP FOR EXISTING USERS
-- ============================================
-- Run these scripts IN ORDER in Supabase SQL Editor

-- ============================================
-- STEP 1: Create the table (run this first!)
-- ============================================

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create unlocked_users table
CREATE TABLE IF NOT EXISTS unlocked_users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    auth_user_id UUID NOT NULL UNIQUE REFERENCES auth.users(id) ON DELETE CASCADE,
    email TEXT NOT NULL,
    username TEXT,
    verified BOOLEAN DEFAULT false,
    discord_webhook_url TEXT,
    discord_notifications JSONB DEFAULT '[]'::jsonb,
    notifications_subscription_active BOOLEAN DEFAULT false,
    notifications_subscription_expires_at TIMESTAMPTZ,
    last_active TIMESTAMPTZ DEFAULT NOW(),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_unlocked_users_auth_user_id ON unlocked_users(auth_user_id);
CREATE INDEX IF NOT EXISTS idx_unlocked_users_verified ON unlocked_users(verified);
CREATE INDEX IF NOT EXISTS idx_unlocked_users_subscription_active ON unlocked_users(notifications_subscription_active);

-- Enable Row Level Security
ALTER TABLE unlocked_users ENABLE ROW LEVEL SECURITY;

-- Create RLS policies
DROP POLICY IF EXISTS "Users can read own unlocked_users data" ON unlocked_users;
CREATE POLICY "Users can read own unlocked_users data"
    ON unlocked_users FOR SELECT
    USING (auth.uid() = auth_user_id);

DROP POLICY IF EXISTS "Users can update own unlocked_users data" ON unlocked_users;
CREATE POLICY "Users can update own unlocked_users data"
    ON unlocked_users FOR UPDATE
    USING (auth.uid() = auth_user_id);

DROP POLICY IF EXISTS "Users can insert own unlocked_users data" ON unlocked_users;
CREATE POLICY "Users can insert own unlocked_users data"
    ON unlocked_users FOR INSERT
    WITH CHECK (auth.uid() = auth_user_id);

-- Grant permissions
GRANT USAGE ON SCHEMA public TO anon, authenticated;
GRANT SELECT, INSERT, UPDATE ON unlocked_users TO authenticated;
GRANT SELECT ON unlocked_users TO anon;

-- ============================================
-- STEP 2: Migrate all existing users
-- ============================================

INSERT INTO unlocked_users (
    auth_user_id,
    email,
    username,
    verified,
    notifications_subscription_active,
    notifications_subscription_expires_at,
    last_active,
    created_at,
    updated_at
)
SELECT 
    u.id as auth_user_id,
    u.email,
    COALESCE(
        u.raw_user_meta_data->>'username',
        u.raw_user_meta_data->>'name',
        SPLIT_PART(u.email, '@', 1),
        ''
    ) as username,
    true as verified,  -- Auto-verify existing users
    true as notifications_subscription_active,  -- Auto-activate notifications
    NULL as notifications_subscription_expires_at,  -- Lifetime subscription
    COALESCE(u.last_sign_in_at, u.created_at, NOW()) as last_active,
    u.created_at,
    NOW() as updated_at
FROM auth.users u
WHERE NOT EXISTS (
    SELECT 1 
    FROM unlocked_users uu 
    WHERE uu.auth_user_id = u.id
)
ON CONFLICT (auth_user_id) DO UPDATE
SET 
    email = EXCLUDED.email,
    username = COALESCE(NULLIF(EXCLUDED.username, ''), unlocked_users.username),
    last_active = GREATEST(unlocked_users.last_active, EXCLUDED.last_active),
    updated_at = NOW();

-- ============================================
-- STEP 3: Show summary
-- ============================================

SELECT 
    COUNT(*) as total_users,
    COUNT(CASE WHEN verified = true THEN 1 END) as verified_users,
    COUNT(CASE WHEN notifications_subscription_active = true THEN 1 END) as active_subscriptions,
    COUNT(CASE WHEN discord_webhook_url IS NOT NULL THEN 1 END) as users_with_webhooks
FROM unlocked_users;
