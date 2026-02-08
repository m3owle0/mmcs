# Subscription Security Guide - Preventing Free Upgrades

## ⚠️ Current Security Issue

**The Problem:** Users can potentially modify their own `subscription_tier` field in the database if Row Level Security (RLS) isn't properly configured.

**The Risk:** Someone could change their tier from `basic` to `pro` without paying.

---

## 🔒 How It Currently Works

### Current Flow:
1. User pays via Stripe
2. **You manually update** `subscription_tier` in Supabase
3. Code reads `subscription_tier` from database
4. Features are unlocked based on tier

### The Vulnerability:
- If RLS policies allow users to UPDATE their own rows, they can change `subscription_tier`
- Client-side code can't be trusted - users can modify JavaScript
- Database must enforce security, not the frontend

---

## ✅ Solution: Row Level Security (RLS) Policies

You need to configure Supabase Row Level Security to **prevent users from modifying their subscription tier**.

### Step 1: Enable RLS on `unlocked_users` Table

1. Go to [Supabase Dashboard](https://supabase.com/dashboard)
2. Select your project
3. Go to **Table Editor** → `unlocked_users`
4. Click **"Enable RLS"** (if not already enabled)

### Step 2: Create RLS Policies

Go to **Authentication** → **Policies** (or **Table Editor** → `unlocked_users` → **Policies**)

#### Policy 1: Users can READ their own data
```sql
-- Allow users to read their own row
CREATE POLICY "Users can read own data"
ON unlocked_users
FOR SELECT
USING (auth.uid() = auth_user_id);
```

#### Policy 2: Users can UPDATE limited fields (NOT subscription_tier)
```sql
-- Allow users to update only safe fields (NOT subscription_tier, verified, etc.)
CREATE POLICY "Users can update safe fields"
ON unlocked_users
FOR UPDATE
USING (auth.uid() = auth_user_id)
WITH CHECK (
  auth.uid() = auth_user_id
  -- Explicitly prevent updating subscription-related fields
  AND subscription_tier IS NOT DISTINCT FROM (SELECT subscription_tier FROM unlocked_users WHERE auth_user_id = auth.uid())
  AND verified IS NOT DISTINCT FROM (SELECT verified FROM unlocked_users WHERE auth_user_id = auth.uid())
  AND subscription_expires_at IS NOT DISTINCT FROM (SELECT subscription_expires_at FROM unlocked_users WHERE auth_user_id = auth.uid())
);
```

**Better Approach:** Only allow updating specific safe fields:

```sql
-- More secure: Only allow updating non-sensitive fields
CREATE POLICY "Users can update profile fields only"
ON unlocked_users
FOR UPDATE
USING (auth.uid() = auth_user_id)
WITH CHECK (
  auth.uid() = auth_user_id
  -- Only allow updating these fields:
  -- username, description, profile_picture_url, discord_webhook_url, discord_notifications
  -- Everything else (subscription_tier, verified, etc.) is blocked
);
```

#### Policy 3: Users can INSERT their own row (for signup)
```sql
-- Allow users to create their own row during signup
CREATE POLICY "Users can insert own row"
ON unlocked_users
FOR INSERT
WITH CHECK (auth.uid() = auth_user_id);
```

#### Policy 4: BLOCK all subscription-related updates
```sql
-- Explicitly deny updates to subscription fields
CREATE POLICY "Block subscription field updates"
ON unlocked_users
FOR UPDATE
USING (false)
WITH CHECK (false)
-- This will be overridden by the more specific policy above, but adds extra security
```

---

## 🛡️ Recommended Secure Setup

### Option A: Service Role Only Updates (Most Secure)

**Only your backend/admin can update subscription fields:**

1. **Disable user updates entirely** - Users can only read
2. **Use Supabase Service Role** for subscription updates
3. **Create a backend function** or use Supabase Dashboard to update subscriptions

**RLS Policy:**
```sql
-- Users can only SELECT their own data
CREATE POLICY "Users can read own data"
ON unlocked_users
FOR SELECT
USING (auth.uid() = auth_user_id);

-- Users can only update profile fields (not subscription)
CREATE POLICY "Users can update profile only"
ON unlocked_users
FOR UPDATE
USING (auth.uid() = auth_user_id)
WITH CHECK (
  auth.uid() = auth_user_id
  -- Block subscription fields explicitly
  AND (subscription_tier IS NULL OR subscription_tier = (SELECT subscription_tier FROM unlocked_users WHERE auth_user_id = auth.uid()))
);
```

### Option B: Function-Based Updates (Secure + Flexible)

Create a database function that validates updates:

```sql
-- Function to safely update user profile (blocks subscription fields)
CREATE OR REPLACE FUNCTION update_user_profile(
  p_username TEXT,
  p_description TEXT,
  p_profile_picture_url TEXT,
  p_discord_webhook_url TEXT,
  p_discord_notifications JSONB
)
RETURNS void
LANGUAGE plpgsql
SECURITY DEFINER
AS $$
BEGIN
  UPDATE unlocked_users
  SET
    username = COALESCE(p_username, username),
    description = COALESCE(p_description, description),
    profile_picture_url = COALESCE(p_profile_picture_url, profile_picture_url),
    discord_webhook_url = COALESCE(p_discord_webhook_url, discord_webhook_url),
    discord_notifications = COALESCE(p_discord_notifications, discord_notifications),
    last_active = NOW()
  WHERE auth_user_id = auth.uid();
  
  -- subscription_tier, verified, subscription_expires_at are NOT updated here
END;
$$;
```

---

## 📋 Step-by-Step: Secure Your Database

### Step 1: Check Current RLS Status

1. Go to Supabase Dashboard
2. **Table Editor** → `unlocked_users`
3. Check if **"RLS Enabled"** is ON
4. If OFF, click **"Enable RLS"**

### Step 2: Review Existing Policies

1. Click **"Policies"** tab (or go to **Authentication** → **Policies**)
2. See what policies exist for `unlocked_users`
3. **Delete any policies that allow UPDATE without restrictions**

### Step 3: Create Secure Policies

Run this SQL in **SQL Editor**:

```sql
-- Step 1: Drop existing UPDATE policies (if any)
DROP POLICY IF EXISTS "Users can update own data" ON unlocked_users;
DROP POLICY IF EXISTS "Enable update for users" ON unlocked_users;

-- Step 2: Create read policy
DROP POLICY IF EXISTS "Users can read own data" ON unlocked_users;
CREATE POLICY "Users can read own data"
ON unlocked_users
FOR SELECT
USING (auth.uid() = auth_user_id);

-- Step 3: Create insert policy (for signup)
DROP POLICY IF EXISTS "Users can insert own row" ON unlocked_users;
CREATE POLICY "Users can insert own row"
ON unlocked_users
FOR INSERT
WITH CHECK (auth.uid() = auth_user_id);

-- Step 4: Create restricted update policy
-- Users can ONLY update: username, description, profile_picture_url, 
-- discord_webhook_url, discord_notifications, last_active
-- They CANNOT update: verified, subscription_tier, subscription_expires_at, payment_method
DROP POLICY IF EXISTS "Users can update profile only" ON unlocked_users;
CREATE POLICY "Users can update profile only"
ON unlocked_users
FOR UPDATE
USING (auth.uid() = auth_user_id)
WITH CHECK (
  auth.uid() = auth_user_id
  -- Prevent changes to subscription fields by checking they haven't changed
  AND (
    -- Get current values
    (SELECT verified FROM unlocked_users WHERE auth_user_id = auth.uid()) 
    IS NOT DISTINCT FROM verified
  )
  AND (
    (SELECT subscription_tier FROM unlocked_users WHERE auth_user_id = auth.uid()) 
    IS NOT DISTINCT FROM subscription_tier
  )
  AND (
    (SELECT subscription_expires_at FROM unlocked_users WHERE auth_user_id = auth.uid()) 
    IS NOT DISTINCT FROM subscription_expires_at
  )
  AND (
    (SELECT payment_method FROM unlocked_users WHERE auth_user_id = auth.uid()) 
    IS NOT DISTINCT FROM payment_method
  )
);
```

### Step 4: Test the Security

1. **Try to update subscription_tier as a user:**
   ```sql
   -- This should FAIL (run as authenticated user, not service role)
   UPDATE unlocked_users 
   SET subscription_tier = 'pro' 
   WHERE auth_user_id = auth.uid();
   ```

2. **Try to update profile fields:**
   ```sql
   -- This should SUCCEED
   UPDATE unlocked_users 
   SET username = 'NewUsername' 
   WHERE auth_user_id = auth.uid();
   ```

---

## 🔐 How to Update Subscriptions Securely

### Method 1: Use Supabase Dashboard (Manual)

1. Go to **Table Editor** → `unlocked_users`
2. Find user by email
3. Update `subscription_tier`, `verified`, `subscription_expires_at`
4. **Only you (admin) can do this** - users cannot

### Method 2: Use Service Role Key (Backend)

Create a backend script/function that uses the **Service Role** key (not the anon key):

```javascript
// Backend only - uses service role key
const { createClient } = require('@supabase/supabase-js');

const supabaseAdmin = createClient(
  'https://your-project.supabase.co',
  'YOUR_SERVICE_ROLE_KEY' // NOT the anon key!
);

// Update subscription (only works with service role)
async function activateSubscription(userEmail, tier) {
  const { data, error } = await supabaseAdmin
    .from('unlocked_users')
    .update({
      verified: true,
      subscription_tier: tier,
      subscription_expires_at: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString(),
      payment_method: 'stripe'
    })
    .eq('email', userEmail);
    
  return { data, error };
}
```

### Method 3: Stripe Webhooks (Automated)

Set up Stripe webhooks to automatically update subscriptions:

1. **Create webhook endpoint** (requires backend server)
2. **Verify webhook signature** (Stripe provides this)
3. **Use Service Role key** to update Supabase
4. **Update subscription_tier** based on Stripe subscription status

---

## ✅ Verification Checklist

After setting up RLS policies:

- [ ] RLS is enabled on `unlocked_users` table
- [ ] Users can READ their own data
- [ ] Users can UPDATE profile fields (username, description, etc.)
- [ ] Users CANNOT UPDATE `subscription_tier`
- [ ] Users CANNOT UPDATE `verified`
- [ ] Users CANNOT UPDATE `subscription_expires_at`
- [ ] Users CANNOT UPDATE `payment_method`
- [ ] Only admin/service role can update subscription fields
- [ ] Tested: User trying to change tier fails
- [ ] Tested: User updating profile succeeds

---

## 🧪 Testing Security

### Test 1: User Cannot Change Tier

1. Log in as a test user
2. Open browser console (F12)
3. Try to run:
   ```javascript
   const { data, error } = await supabase
     .from('unlocked_users')
     .update({ subscription_tier: 'pro' })
     .eq('auth_user_id', 'YOUR_USER_ID');
   console.log('Result:', data, error);
   ```
4. **Expected:** Should fail with permission error

### Test 2: User Can Update Profile

```javascript
const { data, error } = await supabase
  .from('unlocked_users')
  .update({ username: 'NewUsername' })
  .eq('auth_user_id', 'YOUR_USER_ID');
```
**Expected:** Should succeed

---

## 🚨 Important Notes

1. **Never expose Service Role key** in frontend code
2. **Always verify payments** before updating subscription_tier
3. **Check subscription_expires_at** - don't let expired subscriptions access Pro features
4. **Monitor for suspicious activity** - check for users with Pro tier but no payment record
5. **Regular audits** - periodically check subscription_tier matches payment records

---

## 📊 Recommended Database Schema

Make sure your `unlocked_users` table has:

```sql
-- Required columns for subscription security
subscription_tier TEXT DEFAULT 'trial',  -- 'trial', 'basic', 'pro'
verified BOOLEAN DEFAULT false,
subscription_expires_at TIMESTAMPTZ,
payment_method TEXT,  -- 'stripe', 'paypal', etc.
stripe_customer_id TEXT,  -- Store Stripe customer ID for verification
stripe_subscription_id TEXT,  -- Store Stripe subscription ID
```

**Add indexes:**
```sql
CREATE INDEX idx_unlocked_users_subscription_tier ON unlocked_users(subscription_tier);
CREATE INDEX idx_unlocked_users_stripe_customer ON unlocked_users(stripe_customer_id);
```

---

## 🔄 Subscription Verification Flow

### When User Subscribes:

1. **User pays via Stripe** → Stripe creates subscription
2. **You get notification** (email/webhook)
3. **You verify payment** in Stripe Dashboard
4. **You update Supabase** (manually or via webhook):
   ```sql
   UPDATE unlocked_users
   SET 
     verified = true,
     subscription_tier = 'pro',  -- or 'basic'
     subscription_expires_at = NOW() + INTERVAL '30 days',
     payment_method = 'stripe',
     stripe_customer_id = 'cus_xxxxx',
     stripe_subscription_id = 'sub_xxxxx'
   WHERE email = 'user@example.com';
   ```

### When Checking Access:

The code already does this correctly:
- Reads `subscription_tier` from database
- Checks `subscription_expires_at` if set
- Verifies `verified = true`
- **RLS prevents users from modifying these fields**

---

## 🎯 Quick Fix (5 Minutes)

**If you want to secure it RIGHT NOW:**

1. Go to Supabase Dashboard
2. **Table Editor** → `unlocked_users` → **Policies**
3. **Delete any UPDATE policies** that don't restrict subscription fields
4. **Create this policy:**

```sql
CREATE POLICY "Users can update profile only"
ON unlocked_users
FOR UPDATE
USING (auth.uid() = auth_user_id)
WITH CHECK (
  auth.uid() = auth_user_id
  -- Prevent subscription field changes
  AND subscription_tier = (SELECT subscription_tier FROM unlocked_users WHERE auth_user_id = auth.uid())
  AND verified = (SELECT verified FROM unlocked_users WHERE auth_user_id = auth.uid())
);
```

This ensures users can't change their subscription tier!

---

**Need help implementing this?** Let me know and I can guide you through the exact steps for your Supabase setup.
