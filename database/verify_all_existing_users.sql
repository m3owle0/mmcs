-- Verify All Existing Users
-- Run this if you want to automatically verify all existing users

UPDATE unlocked_users 
SET verified = true 
WHERE verified = false;

-- Show updated count
SELECT 
    COUNT(*) as total_users,
    COUNT(CASE WHEN verified = true THEN 1 END) as verified_users
FROM unlocked_users;
