-- Create unlocked_users table (if it doesn't exist)
-- Run this FIRST before migrating users

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

-- Create indexes for faster lookups
CREATE INDEX IF NOT EXISTS idx_unlocked_users_auth_user_id ON unlocked_users(auth_user_id);
CREATE INDEX IF NOT EXISTS idx_unlocked_users_verified ON unlocked_users(verified);
CREATE INDEX IF NOT EXISTS idx_unlocked_users_subscription_active ON unlocked_users(notifications_subscription_active);

-- Enable Row Level Security (RLS)
ALTER TABLE unlocked_users ENABLE ROW LEVEL SECURITY;

-- Create policy: Users can read their own data
DROP POLICY IF EXISTS "Users can read own unlocked_users data" ON unlocked_users;
CREATE POLICY "Users can read own unlocked_users data"
    ON unlocked_users
    FOR SELECT
    USING (auth.uid() = auth_user_id);

-- Create policy: Users can update their own data
DROP POLICY IF EXISTS "Users can update own unlocked_users data" ON unlocked_users;
CREATE POLICY "Users can update own unlocked_users data"
    ON unlocked_users
    FOR UPDATE
    USING (auth.uid() = auth_user_id);

-- Create policy: Users can insert their own data (for initial creation)
DROP POLICY IF EXISTS "Users can insert own unlocked_users data" ON unlocked_users;
CREATE POLICY "Users can insert own unlocked_users data"
    ON unlocked_users
    FOR INSERT
    WITH CHECK (auth.uid() = auth_user_id);

-- Grant necessary permissions
GRANT USAGE ON SCHEMA public TO anon, authenticated;
GRANT SELECT, INSERT, UPDATE ON unlocked_users TO authenticated;
GRANT SELECT ON unlocked_users TO anon;

-- Verify table was created
SELECT 'Table unlocked_users created successfully!' as status;
