# 🔧 Fix Notifier - Step by Step Guide

## Problem: Notifier shows "0 total user(s)"

Your notifier can't see any users because:
1. Either users aren't in the `unlocked_users` table, OR
2. The `anon` key is blocked by RLS (Row Level Security)

---

## ✅ STEP 1: Fix Corrupted Webhook URLs

1. Go to **Supabase Dashboard** → **SQL Editor**
2. Copy and paste this entire script:

```sql
-- Fix Corrupted Webhook URLs
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
```

3. Click **RUN**
4. You should see "Success. No rows returned" or a number of rows updated

---

## ✅ STEP 2: Make Sure Users Are Migrated

1. Still in **Supabase SQL Editor**
2. Copy and paste this entire script:

```sql
-- Migrate all existing auth.users to unlocked_users
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
    true as verified,
    true as notifications_subscription_active,
    NULL as notifications_subscription_expires_at,
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

-- Show how many users we have
SELECT 
    COUNT(*) as total_users,
    COUNT(CASE WHEN verified = true THEN 1 END) as verified_users,
    COUNT(CASE WHEN notifications_subscription_active = true THEN 1 END) as active_subscriptions,
    COUNT(CASE WHEN discord_webhook_url IS NOT NULL AND discord_webhook_url != '' THEN 1 END) as users_with_webhooks
FROM unlocked_users;
```

3. Click **RUN**
4. Check the results - you should see numbers like:
   - `total_users: 5`
   - `active_subscriptions: 5`
   - `users_with_webhooks: 1` (or however many have webhooks)

---

## ✅ STEP 3: Allow Notifier to Read All Users (Fix RLS)

The notifier needs to read ALL users, but RLS is blocking it. Fix this:

1. Still in **Supabase SQL Editor**
2. Copy and paste this entire script:

```sql
-- Allow service_role and anon to read all users for notifier
DROP POLICY IF EXISTS "Allow notifier to read all users" ON unlocked_users;
CREATE POLICY "Allow notifier to read all users"
ON unlocked_users
FOR SELECT
TO anon, authenticated, service_role
USING (true);
```

3. Click **RUN**
4. You should see "Success. No rows returned"

---

## ✅ STEP 4: Get Your Service Role Key (IMPORTANT!)

The `anon` key has limitations. Use the `service_role` key instead:

1. Go to **Supabase Dashboard** → **Settings** → **API**
2. Find **`service_role`** key (NOT the `anon` key)
3. Copy it (it's long, starts with `eyJ...`)
4. **⚠️ KEEP THIS SECRET** - Never share it or put it in frontend code!

---

## ✅ STEP 5: Update start.bat with Service Role Key

1. Open `notifier/start.bat` in a text editor
2. Find the line that says `set SUPABASE_ANON_KEY=...`
3. Replace it with:
   ```batch
   set SUPABASE_SERVICE_ROLE_KEY=YOUR_SERVICE_ROLE_KEY_HERE
   ```
   (Replace `YOUR_SERVICE_ROLE_KEY_HERE` with the actual key from Step 4)
4. Save the file

---

## ✅ STEP 6: Restart the Notifier

1. Close the notifier window (if it's running) - Press `Ctrl+C`
2. Double-click `notifier/start.bat` to run it again
3. You should now see:
   - `📈 Database stats: X total user(s), X active subscription(s)...`
   - `✅ Found X subscriber(s) ready for notifications...`
   - (Where X is a number greater than 0)

---

## ✅ STEP 7: Verify It's Working

Look at the notifier console output. You should see:
- ✅ `Found X subscriber(s) ready for notifications` (X > 0)
- ✅ `✓ Including user [email] ([username]) - webhook: https://discord.com/api/webhooks/...`
- ✅ No more "0 total user(s)" messages

---

## 🆘 If It Still Doesn't Work

1. **Check the console output** - What does it say now?
2. **Verify in Supabase**:
   - Go to **Table Editor** → `unlocked_users`
   - Do you see your users?
   - Does your user have `notifications_subscription_active = true`?
   - Does your user have a `discord_webhook_url` that starts with `https://discord.com/api/webhooks/`?
3. **Check the webhook URL** - Make sure it's clean (no JSON appended)

---

## 📝 Quick Checklist

- [ ] Step 1: Fixed corrupted webhooks (SQL)
- [ ] Step 2: Migrated users to unlocked_users (SQL)
- [ ] Step 3: Fixed RLS policy (SQL)
- [ ] Step 4: Got service_role key from Supabase
- [ ] Step 5: Updated start.bat with service_role key
- [ ] Step 6: Restarted notifier
- [ ] Step 7: Verified it's working

---

**That's it! Follow these steps in order and your notifier should work.**
