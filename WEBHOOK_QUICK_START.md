# Stripe Webhook Quick Start - Auto-Upgrade Users

## 🎯 Goal
Automatically upgrade users to Basic/Pro tier when they pay via Stripe - no manual database updates needed!

---

## ⚡ Fastest Method: Supabase Edge Function (15 minutes)

### Step 0: Install Supabase CLI

**⚠️ `npm install -g supabase` does NOT work!**

**Easiest Method - Direct Download:**
1. Go to: https://github.com/supabase/cli/releases/latest
2. Download: `supabase_windows_amd64.zip`
3. Extract and copy `supabase.exe` to `C:\Windows\System32`
4. Or use `npx supabase@latest` (no installation needed)

**See `INSTALL_SUPABASE_CLI.md` for detailed installation instructions.**

### Step 1: Verify CLI Installation

**⚠️ IMPORTANT:** `npm install -g supabase` does NOT work! Use one of these methods:

**Option A: Scoop (Recommended for Windows)**
```powershell
# Install Scoop first (if you don't have it)
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
irm get.scoop.sh | iex

# Then install Supabase CLI
scoop bucket add supabase https://github.com/supabase/scoop-bucket.git
scoop install supabase
```

**Option B: Direct Download (Easiest)**
1. Go to: https://github.com/supabase/cli/releases/latest
2. Download: `supabase_windows_amd64.zip` (or `supabase_windows_arm64.zip` for ARM)
3. Extract the zip file
4. Copy `supabase.exe` to a folder in your PATH (e.g., `C:\Windows\System32` or create `C:\Tools` and add it to PATH)
5. Or just run it from the extracted folder

**Option C: Chocolatey (If you have it)**
```powershell
choco install supabase
```

**Option D: Use npx (No Installation Needed)**
You can use Supabase CLI without installing:
```bash
npx supabase@latest login
npx supabase@latest link --project-ref your-project-ref
npx supabase@latest functions deploy stripe-webhook
```

### Step 2: Login to Supabase

```bash
supabase login
```

### Step 3: Link Your Project

```bash
cd C:\Users\puppiesandkittens\Downloads\mmcs
supabase link --project-ref your-project-ref
```

(Get project ref from Supabase Dashboard URL: `https://supabase.com/dashboard/project/your-project-ref`)

### Step 4: Create Function

```bash
supabase functions new stripe-webhook
```

### Step 5: Copy Function Code

Copy the code from `supabase/functions/stripe-webhook/index.ts` (I've created this file for you)

### Step 6: Set Secrets

```bash
# Get these from Stripe Dashboard and Supabase Dashboard
supabase secrets set STRIPE_SECRET_KEY=sk_live_xxxxx
supabase secrets set STRIPE_WEBHOOK_SECRET=whsec_xxxxx
supabase secrets set SUPABASE_URL=https://wbpfuuiznsmysbskywdx.supabase.co
supabase secrets set SUPABASE_SERVICE_ROLE_KEY=your_service_role_key_here
```

### Step 7: Deploy

```bash
supabase functions deploy stripe-webhook
```

You'll get a URL like: `https://wbpfuuiznsmysbskywdx.supabase.co/functions/v1/stripe-webhook`

### Step 8: Configure Stripe Webhook

1. Go to [Stripe Dashboard](https://dashboard.stripe.com) → **Developers** → **Webhooks**
2. Click **"+ Add endpoint"**
3. **Endpoint URL:** `https://wbpfuuiznsmysbskywdx.supabase.co/functions/v1/stripe-webhook`
4. **Events to send:**
   - ✅ `checkout.session.completed`
   - ✅ `customer.subscription.updated`
   - ✅ `customer.subscription.deleted`
5. Click **"Add endpoint"**
6. **Copy the signing secret** (starts with `whsec_...`)
7. Update your secrets:
   ```bash
   supabase secrets set STRIPE_WEBHOOK_SECRET=whsec_xxxxx
   ```

---

## 🔑 Getting Your Keys

### Stripe Secret Key:
1. Stripe Dashboard → **Developers** → **API keys**
2. Copy **"Secret key"** (starts with `sk_live_...` or `sk_test_...`)

### Stripe Webhook Secret:
1. Stripe Dashboard → **Developers** → **Webhooks**
2. Click your webhook endpoint
3. Copy **"Signing secret"** (starts with `whsec_...`)

### Supabase Service Role Key:
1. Supabase Dashboard → **Project Settings** → **API**
2. Copy **"service_role"** key (⚠️ Keep this secret! Not the anon key)

---

## 🧪 Testing

### Test Mode:
1. Use Stripe **test mode** (toggle in Stripe Dashboard)
2. Use test webhook endpoint
3. Subscribe with test card: `4242 4242 4242 4242`
4. Check Stripe Dashboard → **Webhooks** → **Recent events** for logs

### Verify It Works:
1. Complete a test subscription
2. Check Supabase → `unlocked_users` table
3. User should have:
   - ✅ `verified = true`
   - ✅ `subscription_tier = 'basic'` or `'pro'`
   - ✅ `subscription_expires_at` = 30 days from now
   - ✅ `payment_method = 'stripe'`

---

## 📋 What Happens When User Subscribes

1. **User clicks "Subscribe with Stripe"** → Opens Stripe checkout
2. **User completes payment** → Stripe processes payment
3. **Stripe sends webhook** → Your Supabase function receives event
4. **Function verifies signature** → Security check
5. **Function updates Supabase** → User upgraded automatically!
6. **User refreshed page** → They see their new tier!

**Total time:** Usually less than 5 seconds!

---

## 🎯 Tier Detection Logic

The function determines tier from payment amount:

- **$5.00 (500 cents)** → `basic` tier
- **$10.00 (1000 cents)** → `pro` tier

This matches your Stripe payment links:
- Basic: `https://buy.stripe.com/00w6oH2q5g6YdryekFcwg00` ($5)
- Pro: `https://buy.stripe.com/3cIbJ1e8Ng6Y1IQfoJcwg01` ($10)

---

## ✅ Checklist

- [ ] Supabase CLI installed
- [ ] Supabase project linked
- [ ] Edge function created
- [ ] Function code copied
- [ ] Secrets set (Stripe key, webhook secret, Supabase keys)
- [ ] Function deployed
- [ ] Stripe webhook endpoint created
- [ ] Webhook events configured
- [ ] Test subscription completed
- [ ] Verified user upgraded in Supabase

---

## 🚨 Troubleshooting

### Function not receiving webhooks?
- Check Stripe Dashboard → Webhooks → Recent events
- Verify endpoint URL is correct
- Check function logs: `supabase functions logs stripe-webhook`

### User not getting upgraded?
- Check if email matches exactly
- Verify service role key is correct
- Check function logs for errors
- Verify webhook secret is correct

### Wrong tier assigned?
- Check `amount_total` in webhook payload
- Verify your Stripe prices match the logic ($5 = Basic, $10 = Pro)

---

## 📚 Full Documentation

See `STRIPE_WEBHOOK_SETUP.md` for:
- Alternative deployment methods (Vercel, custom server)
- More detailed explanations
- Advanced features (subscription renewals, cancellations)

---

**That's it!** Once set up, users will be automatically upgraded when they subscribe. No manual database updates needed! 🎉
