-- Fix Corrupted Webhook URLs
-- This SQL fixes webhook URLs that have JSON data incorrectly appended to them
-- Run this in Supabase SQL Editor

-- Find and fix webhook URLs that have JSON appended (they contain "[{" or "{")
UPDATE unlocked_users
SET discord_webhook_url = SUBSTRING(
    discord_webhook_url 
    FROM 1 
    FOR CASE 
        WHEN POSITION('[{' IN discord_webhook_url) > 0 
        THEN POSITION('[{' IN discord_webhook_url) - 1
        WHEN POSITION('{"' IN discord_webhook_url) > 0 
        THEN POSITION('{"' IN discord_webhook_url) - 1
        ELSE LENGTH(discord_webhook_url)
    END
)
WHERE discord_webhook_url IS NOT NULL
  AND (discord_webhook_url LIKE '%[{%' OR discord_webhook_url LIKE '%{"%')
  AND POSITION('https://discord.com/api/webhooks/' IN discord_webhook_url) > 0;

-- Show how many were fixed
SELECT 
    COUNT(*) as fixed_count,
    'Webhook URLs cleaned (removed appended JSON)' as description
FROM unlocked_users
WHERE discord_webhook_url IS NOT NULL
  AND discord_webhook_url LIKE 'https://discord.com/api/webhooks/%'
  AND NOT (discord_webhook_url LIKE '%[{%' OR discord_webhook_url LIKE '%{"%');
