-- Migrate Existing Authenticated Users to unlocked_users Table
-- Run this in Supabase SQL Editor to create records for all existing auth.users

-- Insert all existing auth.users into unlocked_users table
-- This handles users who signed up before the trigger was set up
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
        SPLIT_PART(u.email, '@', 1),  -- Use email prefix as fallback
        ''
    ) as username,
    false as verified,  -- Set to false by default - you can verify them manually or run UPDATE below
    false as notifications_subscription_active,
    NULL as notifications_subscription_expires_at,
    COALESCE(u.last_sign_in_at, u.created_at, NOW()) as last_active,
    u.created_at,
    NOW() as updated_at
FROM auth.users u
WHERE NOT EXISTS (
    -- Only insert if record doesn't already exist
    SELECT 1 
    FROM unlocked_users uu 
    WHERE uu.auth_user_id = u.id
)
ON CONFLICT (auth_user_id) DO NOTHING;

-- Optional: Verify all existing users automatically
-- Uncomment the line below if you want all existing users to be verified immediately
-- UPDATE unlocked_users SET verified = true WHERE verified = false;

-- Optional: Activate notifications for all existing users
-- Uncomment the lines below if you want all existing users to have active notifications
-- UPDATE unlocked_users 
-- SET notifications_subscription_active = true,
--     notifications_subscription_expires_at = NULL  -- Lifetime subscription
-- WHERE notifications_subscription_active = false;

-- Show summary of what was created
SELECT 
    COUNT(*) as total_users,
    COUNT(CASE WHEN verified = true THEN 1 END) as verified_users,
    COUNT(CASE WHEN verified = false THEN 1 END) as unverified_users,
    COUNT(CASE WHEN notifications_subscription_active = true THEN 1 END) as active_subscriptions
FROM unlocked_users;
