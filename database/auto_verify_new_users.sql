-- Auto-Verify New Users (But Keep Notifications Inactive)
-- This updates the trigger function to automatically verify new accounts
-- but keep notifications_subscription_active = false by default
-- Run this in Supabase SQL Editor

-- Update the function to auto-verify new users but keep notifications inactive
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
        true,  -- Auto-verify new accounts
        false, -- Keep notifications inactive (requires $5/month donation)
        NULL   -- No expiration date (notifications not active)
    )
    ON CONFLICT (auth_user_id) DO NOTHING;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Verify the update
SELECT 
    'Trigger function updated successfully' as status,
    'New users will be auto-verified but notifications will be inactive' as note;
