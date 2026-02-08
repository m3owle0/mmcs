# Quick 5-Minute Test Guide

## Step 1: Database Setup (30 seconds)

Run this in Supabase SQL Editor:

```sql
ALTER TABLE unlocked_users 
ADD COLUMN IF NOT EXISTS discord_webhook_url TEXT,
ADD COLUMN IF NOT EXISTS discord_notifications JSONB,
ADD COLUMN IF NOT EXISTS discord_subscription_active BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS discord_subscription_expires_at TIMESTAMPTZ;
```

## Step 2: Get Discord Webhook (1 minute)

1. Discord → Right-click channel → Edit Channel → Integrations → Webhooks → New Webhook
2. Copy the webhook URL

## Step 3: Activate Your Subscription (30 seconds)

Run this SQL (replace YOUR_EMAIL):

```sql
UPDATE unlocked_users 
SET 
  discord_subscription_active = TRUE,
  discord_subscription_expires_at = NOW() + INTERVAL '30 days'
WHERE email = 'YOUR_EMAIL';
```

## Step 4: Test in Website (1 minute)

1. Refresh website
2. Click profile → "🔔 Discord Notifications ($5/mo)"
3. Paste webhook URL
4. Click "Test Webhook"
5. ✅ Check Discord - you should see a message!

## Step 5: Add Test Notification (1 minute)

1. In Discord Notifications modal
2. Click "+ Add Search Term"
3. Enter "test" (or any term)
4. Click "Add"
5. Click "Save"

## Step 6: Run Golang Service (2 minutes)

### Windows:

```powershell
# 1. Install Go from https://go.dev/dl/

# 2. Create folder
mkdir C:\discord-notifier
cd C:\discord-notifier

# 3. Create main.go (copy code from TESTING_GUIDE.md)

# 4. Get your Supabase Anon Key:
#    Supabase Dashboard → Project Settings → API Keys → Copy "Publishable key" (sb_publishable_...) OR Legacy "anon" key

# 5. Run:
$env:SUPABASE_ANON_KEY="paste_your_key_here"
go run main.go
```

**You should see:**
```
✅ Found 1 subscriber(s)
👤 Processing: yourname (your@email.com)
   🔍 Checking: 'test'
   ✅ Sending test notification...
   ✅ Notification sent!
```

**Check Discord** - You should see the notification!

## That's It! 🎉

If you see the notification in Discord, everything is working!

## Next: Implement Real Search Logic

Replace the test code in `processUserNotifications()` with your actual market search implementation.
