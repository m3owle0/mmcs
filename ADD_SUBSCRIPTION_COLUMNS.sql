-- Add Subscription Columns to unlocked_users Table
-- Run this SQL in Supabase SQL Editor

-- Step 1: Add subscription_tier column
ALTER TABLE unlocked_users 
ADD COLUMN IF NOT EXISTS subscription_tier TEXT DEFAULT 'trial';

-- Step 2: Add subscription expiry date
ALTER TABLE unlocked_users 
ADD COLUMN IF NOT EXISTS subscription_expires_at TIMESTAMPTZ;

-- Step 3: Add payment method tracking
ALTER TABLE unlocked_users 
ADD COLUMN IF NOT EXISTS payment_method TEXT;

-- Step 4: Add Stripe customer/subscription IDs (optional but recommended)
ALTER TABLE unlocked_users 
ADD COLUMN IF NOT EXISTS stripe_customer_id TEXT;
ALTER TABLE unlocked_users 
ADD COLUMN IF NOT EXISTS stripe_subscription_id TEXT;

-- Step 5: Add indexes for faster queries
CREATE INDEX IF NOT EXISTS idx_unlocked_users_subscription_tier 
ON unlocked_users(subscription_tier);

CREATE INDEX IF NOT EXISTS idx_unlocked_users_subscription_expires 
ON unlocked_users(subscription_expires_at);

CREATE INDEX IF NOT EXISTS idx_unlocked_users_stripe_customer 
ON unlocked_users(stripe_customer_id);

-- Step 6: Add comments for documentation
COMMENT ON COLUMN unlocked_users.subscription_tier IS 'Subscription tier: trial, basic, or pro';
COMMENT ON COLUMN unlocked_users.subscription_expires_at IS 'When the subscription expires (for recurring subscriptions)';
COMMENT ON COLUMN unlocked_users.payment_method IS 'Payment method used: stripe, paypal, cashapp, bitcoin, one-time';
COMMENT ON COLUMN unlocked_users.stripe_customer_id IS 'Stripe customer ID for subscription verification';
COMMENT ON COLUMN unlocked_users.stripe_subscription_id IS 'Stripe subscription ID for tracking';

-- Step 7: Set default values for existing users
-- All existing users should be on 'trial' tier
UPDATE unlocked_users 
SET subscription_tier = 'trial' 
WHERE subscription_tier IS NULL;

-- Step 8: Verify columns were added
SELECT 
  column_name, 
  data_type, 
  column_default,
  is_nullable
FROM information_schema.columns
WHERE table_name = 'unlocked_users'
  AND column_name IN ('subscription_tier', 'subscription_expires_at', 'payment_method', 'stripe_customer_id', 'stripe_subscription_id')
ORDER BY column_name;
