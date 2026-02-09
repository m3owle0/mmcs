-- Multi-Market Clothing Search Database Schema
-- Run this SQL in your Supabase SQL Editor

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

-- Create index for faster lookups
CREATE INDEX IF NOT EXISTS idx_unlocked_users_auth_user_id ON unlocked_users(auth_user_id);
CREATE INDEX IF NOT EXISTS idx_unlocked_users_verified ON unlocked_users(verified);
CREATE INDEX IF NOT EXISTS idx_unlocked_users_subscription_active ON unlocked_users(notifications_subscription_active);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to auto-update updated_at
CREATE TRIGGER update_unlocked_users_updated_at
    BEFORE UPDATE ON unlocked_users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Create function to automatically create unlocked_users record when user signs up
CREATE OR REPLACE FUNCTION handle_new_user()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO public.unlocked_users (
        auth_user_id, 
        email, 
        username, 
        verified,
        notifications_subscription_active,
        notifications_subscription_expires_at
    )
    VALUES (
        NEW.id,
        NEW.email,
        COALESCE(NEW.raw_user_meta_data->>'username', SPLIT_PART(NEW.email, '@', 1), ''),
        true,  -- Auto-verify new accounts (site access is free)
        false, -- Keep notifications inactive (requires $5/month donation)
        NULL   -- No expiration date (notifications not active)
    )
    ON CONFLICT (auth_user_id) DO NOTHING;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Create trigger to auto-create unlocked_users on signup
CREATE TRIGGER on_auth_user_created
    AFTER INSERT ON auth.users
    FOR EACH ROW
    EXECUTE FUNCTION handle_new_user();

-- Enable Row Level Security (RLS)
ALTER TABLE unlocked_users ENABLE ROW LEVEL SECURITY;

-- Create policy: Users can read their own data
CREATE POLICY "Users can read own unlocked_users data"
    ON unlocked_users
    FOR SELECT
    USING (auth.uid() = auth_user_id);

-- Create policy: Users can update their own data
CREATE POLICY "Users can update own unlocked_users data"
    ON unlocked_users
    FOR UPDATE
    USING (auth.uid() = auth_user_id);

-- Create policy: Users can insert their own data (for initial creation)
CREATE POLICY "Users can insert own unlocked_users data"
    ON unlocked_users
    FOR INSERT
    WITH CHECK (auth.uid() = auth_user_id);

-- Grant necessary permissions
GRANT USAGE ON SCHEMA public TO anon, authenticated;
GRANT SELECT, INSERT, UPDATE ON unlocked_users TO authenticated;
GRANT SELECT ON unlocked_users TO anon;

-- Note: The notifier service will need to use the service_role key to read all users
-- The service_role key bypasses RLS, so it can read all records for notification processing
