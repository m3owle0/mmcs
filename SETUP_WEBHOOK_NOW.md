# 🚀 Quick Webhook Setup - Ready to Run!

## ✅ What You Have:
- ✅ Stripe Secret Key: `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`
- ✅ Stripe Publishable Key: `pk_live_51SvasDEA9uVvtrPe4G0GoaW5ZRdjpZgrPpHscMBktFdZT0rMghg56eY2lFtD50bk2UmbNOPwmspKrMyDOR2qcPY5008RBgdDF5`
- ✅ Supabase URL: `https://wbpfuuiznsmysbskywdx.supabase.co`

## ⚠️ What You Still Need:
- ❌ Supabase Service Role Key (get from Supabase Dashboard → Project Settings → API)
- ❌ Stripe Webhook Secret (get after creating webhook endpoint)

---

## 📋 Step-by-Step Setup

### Step 1: Install Supabase CLI (if not already installed)

**Option A: Use npx (No Installation Needed)**
```powershell
# Just use npx - no installation required!
npx supabase@latest login
```

**Option B: Direct Download**
1. Go to: https://github.com/supabase/cli/releases/latest
2. Download: `supabase_windows_amd64.zip`
3. Extract and copy `supabase.exe` to `C:\Windows\System32`

**See `INSTALL_SUPABASE_CLI.md` for more options.**

### Step 2: Login to Supabase

```powershell
supabase login
# OR if using npx:
npx supabase@latest login
```

### Step 3: Link Your Project

```powershell
cd C:\Users\puppiesandkittens\Downloads\mmcs
supabase link --project-ref wbpfuuiznsmysbskywdx
# OR if using npx:
npx supabase@latest link --project-ref wbpfuuiznsmysbskywdx
```

### Step 4: Create the Function (if not already created)

```powershell
supabase functions new stripe-webhook
# OR if using npx:
npx supabase@latest functions new stripe-webhook
```

The function code is already in `supabase/functions/stripe-webhook/index.ts` ✅

### Step 5: Get Your Supabase Service Role Key

1. Go to: https://supabase.com/dashboard/project/wbpfuuiznsmysbskywdx/settings/api
2. Scroll down to **"Project API keys"**
3. Find **"service_role"** key (⚠️ NOT the anon key!)
4. Copy it (starts with `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...`)

### Step 6: Set Secrets (Run These Commands)

**Replace `YOUR_SERVICE_ROLE_KEY_HERE` with the key from Step 5:**

```powershell
# Set Stripe Secret Key
supabase secrets set STRIPE_SECRET_KEY=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

# Set Supabase URL
supabase secrets set SUPABASE_URL=https://wbpfuuiznsmysbskywdx.supabase.co

# Set Supabase Service Role Key (REPLACE WITH YOUR KEY!)
supabase secrets set SUPABASE_SERVICE_ROLE_KEY=YOUR_SERVICE_ROLE_KEY_HERE

# Webhook secret will be set later (after creating webhook endpoint)
```

**If using npx, add `npx supabase@latest` before each command:**
```powershell
npx supabase@latest secrets set STRIPE_SECRET_KEY=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
npx supabase@latest secrets set SUPABASE_URL=https://wbpfuuiznsmysbskywdx.supabase.co
npx supabase@latest secrets set SUPABASE_SERVICE_ROLE_KEY=YOUR_SERVICE_ROLE_KEY_HERE
```

### Step 7: Deploy the Function

```powershell
supabase functions deploy stripe-webhook
# OR if using npx:
npx supabase@latest functions deploy stripe-webhook
```

**You'll get a URL like:**
```
https://wbpfuuiznsmysbskywdx.supabase.co/functions/v1/stripe-webhook
```

**Copy this URL!** You'll need it for the next step.

### Step 8: Create Stripe Webhook Endpoint

1. Go to: https://dashboard.stripe.com/webhooks
2. Click **"+ Add endpoint"**
3. **Endpoint URL:** Paste the URL from Step 7:
   ```
   https://wbpfuuiznsmysbskywdx.supabase.co/functions/v1/stripe-webhook
   ```
4. **Events to send:** Select these events:
   - ✅ `checkout.session.completed`
   - ✅ `customer.subscription.updated`
   - ✅ `customer.subscription.deleted`
5. Click **"Add endpoint"**
6. **Copy the "Signing secret"** (starts with `whsec_...`)

### Step 9: Set Webhook Secret

```powershell
supabase secrets set STRIPE_WEBHOOK_SECRET=whsec_YOUR_WEBHOOK_SECRET_HERE
# OR if using npx:
npx supabase@latest secrets set STRIPE_WEBHOOK_SECRET=whsec_YOUR_WEBHOOK_SECRET_HERE
```

---

## ✅ Testing

### Test with Stripe Test Mode:

1. **Switch to Test Mode** in Stripe Dashboard (toggle in top right)
2. **Use test card:** `4242 4242 4242 4242`
3. **Complete a test subscription** using one of your Stripe checkout links
4. **Check Supabase:** Go to `unlocked_users` table
5. **Verify user has:**
   - ✅ `verified = true`
   - ✅ `subscription_tier = 'basic'` or `'pro'`
   - ✅ `subscription_expires_at` = 30 days from now
   - ✅ `payment_method = 'stripe'`

### Check Webhook Logs:

```powershell
supabase functions logs stripe-webhook
# OR if using npx:
npx supabase@latest functions logs stripe-webhook
```

---

## 🎯 Quick Command Reference

**All commands in one place (replace YOUR_SERVICE_ROLE_KEY and YOUR_WEBHOOK_SECRET):**

```powershell
# Login
npx supabase@latest login

# Link project
cd C:\Users\puppiesandkittens\Downloads\mmcs
npx supabase@latest link --project-ref wbpfuuiznsmysbskywdx

# Set secrets
npx supabase@latest secrets set STRIPE_SECRET_KEY=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
npx supabase@latest secrets set SUPABASE_URL=https://wbpfuuiznsmysbskywdx.supabase.co
npx supabase@latest secrets set SUPABASE_SERVICE_ROLE_KEY=YOUR_SERVICE_ROLE_KEY_HERE

# Deploy function
npx supabase@latest functions deploy stripe-webhook

# After creating webhook endpoint in Stripe:
npx supabase@latest secrets set STRIPE_WEBHOOK_SECRET=whsec_YOUR_WEBHOOK_SECRET_HERE
```

---

## 🚨 Troubleshooting

### "Command not found" or "supabase: command not recognized"
- Use `npx supabase@latest` instead of just `supabase`
- Or install CLI properly (see `INSTALL_SUPABASE_CLI.md`)

### Function deployment fails
- Make sure you're logged in: `npx supabase@latest login`
- Make sure project is linked: `npx supabase@latest link --project-ref wbpfuuiznsmysbskywdx`
- Check function code exists: `supabase/functions/stripe-webhook/index.ts`

### Webhook not receiving events
- Check Stripe Dashboard → Webhooks → Recent events
- Verify endpoint URL matches exactly
- Check function logs: `npx supabase@latest functions logs stripe-webhook`

### User not getting upgraded
- Verify email matches exactly in Supabase
- Check service role key is correct (not anon key!)
- Check function logs for errors
- Verify webhook secret is set correctly

---

## 🎉 You're Done!

Once all steps are complete, users will be **automatically upgraded** when they subscribe via Stripe. No manual database updates needed!

**The webhook will:**
- ✅ Upgrade users to Basic/Pro tier when they pay
- ✅ Handle subscription renewals automatically
- ✅ Handle cancellations automatically
- ✅ Update expiry dates automatically

---

**Need help?** Check:
- `WEBHOOK_QUICK_START.md` - Full detailed guide
- `STRIPE_WEBHOOK_SETUP.md` - Alternative deployment methods
- `WEBHOOK_STATUS_CHECK.md` - Status checklist
