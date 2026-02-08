# Stripe Webhook Status Check

## ❌ Current Status: NOT Set Up Yet

The webhook code is **created** but **not deployed**. Here's what needs to be done:

---

## ✅ What's Already Done

- ✅ Webhook function code created (`supabase/functions/stripe-webhook/index.ts`)
- ✅ Function logic written (handles checkout.session.completed, subscription updates)
- ✅ Tier detection logic (Basic $5, Pro $10)

---

## ❌ What Still Needs to Be Done

### 1. Deploy the Function
The function code exists but hasn't been deployed to Supabase yet.

### 2. Configure Stripe Webhook
Stripe doesn't know where to send webhook events yet.

### 3. Set Environment Variables
The function needs your Stripe and Supabase keys.

---

## 🚀 Quick Setup (5 Minutes)

### Step 1: Deploy Function (Choose One Method)

**Method A: Using npx (No Installation)**
```powershell
cd C:\Users\puppiesandkittens\Downloads\mmcs

# Login
npx supabase@latest login

# Link project
npx supabase@latest link --project-ref wbpfuuiznsmysbskywdx

# Set secrets (get these from Stripe/Supabase dashboards)
npx supabase@latest secrets set STRIPE_SECRET_KEY=sk_live_xxxxx
npx supabase@latest secrets set STRIPE_WEBHOOK_SECRET=whsec_xxxxx
npx supabase@latest secrets set SUPABASE_URL=https://wbpfuuiznsmysbskywdx.supabase.co
npx supabase@latest secrets set SUPABASE_SERVICE_ROLE_KEY=your_service_role_key

# Deploy function
npx supabase@latest functions deploy stripe-webhook
```

**Method B: If You Installed Supabase CLI**
```powershell
supabase login
supabase link --project-ref wbpfuuiznsmysbskywdx
supabase secrets set STRIPE_SECRET_KEY=sk_live_xxxxx
supabase secrets set STRIPE_WEBHOOK_SECRET=whsec_xxxxx
supabase secrets set SUPABASE_URL=https://wbpfuuiznsmysbskywdx.supabase.co
supabase secrets set SUPABASE_SERVICE_ROLE_KEY=your_service_role_key
supabase functions deploy stripe-webhook
```

After deploying, you'll get a URL like:
`https://wbpfuuiznsmysbskywdx.supabase.co/functions/v1/stripe-webhook`

### Step 2: Configure Stripe Webhook

1. Go to [Stripe Dashboard](https://dashboard.stripe.com) → **Developers** → **Webhooks**
2. Click **"+ Add endpoint"**
3. **Endpoint URL:** Paste your function URL from Step 1
4. **Events to send:**
   - ✅ `checkout.session.completed`
   - ✅ `customer.subscription.updated`
   - ✅ `customer.subscription.deleted`
5. Click **"Add endpoint"**
6. **Copy the signing secret** (starts with `whsec_...`)
7. Update your secrets with the webhook secret:
   ```powershell
   npx supabase@latest secrets set STRIPE_WEBHOOK_SECRET=whsec_xxxxx
   ```

### Step 3: Test It

1. Go to Stripe Dashboard → **Webhooks** → Click your endpoint
2. Click **"Send test webhook"**
3. Select `checkout.session.completed`
4. Check if it succeeds

Or do a real test:
1. Subscribe with test card: `4242 4242 4242 4242`
2. Check Supabase → `unlocked_users` table
3. User should be upgraded automatically!

---

## 🔑 Where to Get Your Keys

### Stripe Secret Key:
1. Stripe Dashboard → **Developers** → **API keys**
2. Copy **"Secret key"** (starts with `sk_live_...` or `sk_test_...`)

### Stripe Webhook Secret:
1. Stripe Dashboard → **Developers** → **Webhooks**
2. Click your webhook endpoint
3. Copy **"Signing secret"** (starts with `whsec_...`)

### Supabase Service Role Key:
1. Supabase Dashboard → **Project Settings** → **API**
2. Copy **"service_role"** key (⚠️ Keep secret! Not the anon key)

---

## ✅ Checklist

- [ ] Function code exists (`supabase/functions/stripe-webhook/index.ts`)
- [ ] Supabase CLI installed OR using npx
- [ ] Logged into Supabase
- [ ] Project linked
- [ ] Secrets set (Stripe key, webhook secret, Supabase keys)
- [ ] Function deployed
- [ ] Got function URL
- [ ] Stripe webhook endpoint created
- [ ] Webhook events configured
- [ ] Tested with test subscription

---

## 🧪 How to Verify It's Working

### Test 1: Check Function Logs
```powershell
npx supabase@latest functions logs stripe-webhook
```

### Test 2: Check Stripe Webhook Logs
1. Stripe Dashboard → **Webhooks** → Your endpoint
2. Click **"Recent events"**
3. Should see successful webhook deliveries

### Test 3: Real Test Subscription
1. Use test card: `4242 4242 4242 4242`
2. Complete subscription
3. Check Supabase → `unlocked_users` table
4. User should have:
   - `verified = true`
   - `subscription_tier = 'basic'` or `'pro'`
   - `subscription_expires_at` = 30 days from now

---

## 🚨 Common Issues

### Function not receiving webhooks?
- ✅ Check function URL is correct in Stripe
- ✅ Check webhook secret matches
- ✅ Check function logs for errors

### User not getting upgraded?
- ✅ Check email matches exactly
- ✅ Check service role key is correct
- ✅ Check function logs: `npx supabase@latest functions logs stripe-webhook`

### Wrong tier assigned?
- ✅ Check `amount_total` logic (500 cents = Basic, 1000 cents = Pro)
- ✅ Verify Stripe prices match

---

## 📝 Current Status

**Right now:** Webhook is **NOT active** - users won't be upgraded automatically.

**After setup:** Webhook **WILL be active** - users upgrade automatically within seconds!

---

**Need help with any step?** Let me know where you're stuck!
