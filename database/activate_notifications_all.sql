-- Activate Notifications for All Users (Lifetime)
-- Run this if you want all users to have active notification subscriptions

UPDATE unlocked_users 
SET notifications_subscription_active = true,
    notifications_subscription_expires_at = NULL  -- NULL = lifetime subscription
WHERE notifications_subscription_active = false;

-- Show updated count
SELECT 
    COUNT(*) as total_users,
    COUNT(CASE WHEN notifications_subscription_active = true THEN 1 END) as active_subscriptions
FROM unlocked_users;
