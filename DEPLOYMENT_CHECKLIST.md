# ✅ Deployment Checklist - Verify Everything Works

## 🎯 What Should Be Deployed

1. ✅ **Website** (index.html) - Deployed to Netlify/Vercel/etc.
2. ⚠️ **Supabase Edge Function** (Stripe webhook) - Needs GitHub Actions deployment
3. ⚠️ **Stripe Webhook Configuration** - Needs to be set up in Stripe Dashboard

---

## ✅ Step 1: Verify Website is Live

1. **Check your site URL** (Netlify/Vercel/etc.)
2. **Test the site:**
   - ✅ Page loads
   - ✅ Can search markets
   - ✅ Can sign up/login
   - ✅ Subscription buttons work

**If site is working:** ✅ **Website deployment successful!**

---

## ⚠️ Step 2: Deploy Supabase Edge Function (If Not Done)

Your Stripe webhook function needs to be deployed via GitHub Actions.

### Check if Already Deployed:

1. **Go to:** https://supabase.com/dashboard/project/wbpfuuiznsmysbskywdx/functions
2. **Look for:** `stripe-webhook` function
3. **If it exists:** ✅ Function is deployed!
4. **If it doesn't exist:** Follow steps below

### Deploy via GitHub Actions:

1. **Go to your GitHub repository**
2. **Click:** Actions tab
3. **Check if workflow has run:**
   - ✅ If you see "Deploy Supabase Edge Functions" workflow → Check if it succeeded
   - ❌ If no workflow exists → You need to set up secrets first

### If Workflow Failed:

**Check if secrets are set:**
1. Go to: Repository → Settings → Secrets and variables → Actions
2. Verify these secrets exist:
   - ✅ `SUPABASE_ACCESS_TOKEN`
   - ✅ `SUPABASE_PROJECT_REF`
   - ✅ `SUPABASE_URL`
   - ✅ `SUPABASE_SERVICE_ROLE_KEY`
   - ✅ `STRIPE_SECRET_KEY`
   - ⚠️ `STRIPE_WEBHOOK_SECRET` (can be set later)

**If secrets are missing:**
- See `GITHUB_ENTERPRISE_DEPLOY.md` for setup instructions

**If secrets exist but workflow failed:**
- Click on the failed workflow run
- Check error messages
- Common issues:
  - Invalid access token → Regenerate in Supabase Dashboard
  - Wrong project ref → Should be `wbpfuuiznsmysbskywdx`
  - Missing webhook secret → Can use placeholder for now

### Manual Trigger (If Needed):

1. Go to: Actions tab → "Deploy Supabase Edge Functions"
2. Click: "Run workflow" → "Run workflow"
3. Wait for completion
4. Check logs for errors

---

## ⚠️ Step 3: Configure Stripe Webhook

After the Edge Function is deployed, configure Stripe:

### Get Your Webhook URL:

```
https://wbpfuuiznsmysbskywdx.supabase.co/functions/v1/stripe-webhook
```

### Set Up in Stripe:

1. **Go to:** https://dashboard.stripe.com/webhooks
2. **Click:** "+ Add endpoint"
3. **Endpoint URL:** `https://wbpfuuiznsmysbskywdx.supabase.co/functions/v1/stripe-webhook`
4. **Events to send:**
   - ✅ `checkout.session.completed`
   - ✅ `customer.subscription.updated`
   - ✅ `customer.subscription.deleted`
5. **Click:** "Add endpoint"
6. **Copy the "Signing secret"** (starts with `whsec_...`)

### Update GitHub Secret:

1. **Go to:** GitHub → Repository → Settings → Secrets → Actions
2. **Update:** `STRIPE_WEBHOOK_SECRET` with the `whsec_...` value
3. **Redeploy:** Trigger the workflow again or push a change

---

## 🧪 Step 4: Test Everything

### Test Website:
- ✅ Site loads
- ✅ Can create account
- ✅ Can search markets
- ✅ Subscription buttons link to Stripe

### Test Stripe Webhook:

1. **Use Stripe Test Mode:**
   - Toggle "Test mode" in Stripe Dashboard
   - Use test card: `4242 4242 4242 4242`

2. **Complete a test subscription:**
   - Go to your site
   - Click "Subscribe" (Basic or Pro)
   - Complete checkout with test card
   - Use any future expiry date, any CVC

3. **Check Supabase:**
   - Go to: Supabase Dashboard → Table Editor → `unlocked_users`
   - Find your test user
   - Verify:
     - ✅ `verified = true`
     - ✅ `subscription_tier = 'basic'` or `'pro'`
     - ✅ `subscription_expires_at` = 30 days from now
     - ✅ `payment_method = 'stripe'`

4. **Check Stripe Webhook Logs:**
   - Go to: Stripe Dashboard → Webhooks → Your endpoint
   - Click: "Recent events"
   - Should see: `checkout.session.completed` event
   - Status should be: ✅ `200 OK`

5. **Check Function Logs:**
   - Go to: Supabase Dashboard → Edge Functions → `stripe-webhook` → Logs
   - Should see: "✅ Successfully activated subscription for [email]"

---

## ✅ Final Checklist

- [ ] Website is live and accessible
- [ ] Supabase Edge Function is deployed
- [ ] GitHub Actions workflow runs successfully
- [ ] Stripe webhook endpoint is created
- [ ] Stripe webhook secret is set in GitHub secrets
- [ ] Test subscription completes successfully
- [ ] User is automatically upgraded in Supabase
- [ ] Webhook events show in Stripe Dashboard

---

## 🚨 Common Issues

### Website works but subscriptions don't upgrade:
- ⚠️ Edge Function not deployed → Deploy via GitHub Actions
- ⚠️ Stripe webhook not configured → Set up in Stripe Dashboard
- ⚠️ Webhook secret missing → Add to GitHub secrets and redeploy
- ⚠️ Email mismatch → Check email in Stripe matches Supabase

### GitHub Actions workflow fails:
- ⚠️ Missing secrets → Add all required secrets
- ⚠️ Invalid access token → Regenerate in Supabase Dashboard
- ⚠️ Wrong project ref → Should be `wbpfuuiznsmysbskywdx`

### Stripe webhook returns errors:
- ⚠️ Function not deployed → Deploy via GitHub Actions
- ⚠️ Wrong webhook URL → Should end with `/functions/v1/stripe-webhook`
- ⚠️ Missing webhook secret → Add to GitHub secrets and redeploy

---

## 🎉 Success Indicators

**Everything is working when:**
- ✅ Website loads and functions correctly
- ✅ Users can subscribe via Stripe
- ✅ Subscriptions automatically upgrade users in Supabase
- ✅ Webhook events appear in Stripe Dashboard
- ✅ Function logs show successful processing

---

**Need help?** Check:
- `GITHUB_ENTERPRISE_DEPLOY.md` - GitHub Actions setup
- `SETUP_WEBHOOK_NOW.md` - Manual deployment guide
- `WEBHOOK_QUICK_START.md` - Quick webhook setup
