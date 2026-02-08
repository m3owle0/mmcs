# Quick Fix: Missing subscription_tier Column

## ⚠️ Error You're Seeing

```
ERROR: 42703: column "subscription_tier" does not exist
```

This means the database column doesn't exist yet, but your code is trying to use it.

---

## ✅ Quick Fix (2 Minutes)

### Step 1: Open Supabase SQL Editor

1. Go to [Supabase Dashboard](https://supabase.com/dashboard)
2. Select your project
3. Click **"SQL Editor"** in the left sidebar
4. Click **"New query"**

### Step 2: Run This SQL

Copy and paste this entire SQL script:

```sql
-- Add subscription_tier column
ALTER TABLE unlocked_users 
ADD COLUMN IF NOT EXISTS subscription_tier TEXT DEFAULT 'trial';

-- Add subscription expiry date
ALTER TABLE unlocked_users 
ADD COLUMN IF NOT EXISTS subscription_expires_at TIMESTAMPTZ;

-- Add payment method tracking
ALTER TABLE unlocked_users 
ADD COLUMN IF NOT EXISTS payment_method TEXT;

-- Add indexes for faster queries
CREATE INDEX IF NOT EXISTS idx_unlocked_users_subscription_tier 
ON unlocked_users(subscription_tier);

CREATE INDEX IF NOT EXISTS idx_unlocked_users_subscription_expires 
ON unlocked_users(subscription_expires_at);

-- Set all existing users to 'trial' tier
UPDATE unlocked_users 
SET subscription_tier = 'trial' 
WHERE subscription_tier IS NULL;
```

### Step 3: Click "Run"

Click the **"Run"** button (or press `Ctrl+Enter`)

### Step 4: Verify It Worked

You should see:
- ✅ Success message
- ✅ No errors

---

## 🧪 Test It

After running the SQL, test that it works:

```sql
-- Check if column exists
SELECT subscription_tier, verified, email 
FROM unlocked_users 
LIMIT 5;
```

This should return results without errors.

---

## 📋 What This Does

1. **Adds `subscription_tier` column** - Stores 'trial', 'basic', or 'pro'
2. **Adds `subscription_expires_at`** - Tracks when subscription expires
3. **Adds `payment_method`** - Tracks how they paid
4. **Creates indexes** - Makes queries faster
5. **Sets defaults** - All existing users get 'trial' tier

---

## 🔒 Next Step: Secure It

After adding the columns, **set up Row Level Security** so users can't change their own tier.

See `SUBSCRIPTION_SECURITY_GUIDE.md` for instructions.

---

## ✅ Done!

Once you run this SQL, your code will work correctly and users will be properly tracked by subscription tier.
