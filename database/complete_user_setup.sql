-- Complete User Setup: Migrate + Verify + Activate Notifications
-- This is a one-stop script that does everything

-- Step 1: Migrate all existing auth.users to unlocked_users
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

-- Step 2: Verify any users that weren't verified
UPDATE unlocked_users 
SET verified = true 
WHERE verified = false;

-- Step 3: Activate notifications for any users that don't have it
UPDATE unlocked_users 
SET notifications_subscription_active = true,
    notifications_subscription_expires_at = NULL
WHERE notifications_subscription_active = false;

-- Show final summary
SELECT 
    COUNT(*) as total_users,
    COUNT(CASE WHEN verified = true THEN 1 END) as verified_users,
    COUNT(CASE WHEN notifications_subscription_active = true THEN 1 END) as active_subscriptions,
    COUNT(CASE WHEN discord_webhook_url IS NOT NULL THEN 1 END) as users_with_webhooks
FROM unlocked_users;
