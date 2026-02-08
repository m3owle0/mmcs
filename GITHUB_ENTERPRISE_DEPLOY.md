# 🚀 Deploy Supabase Edge Functions via GitHub Enterprise

This guide shows you how to automatically deploy your Stripe webhook function using GitHub Enterprise (GitHub Actions).

---

## 📋 Prerequisites

- ✅ GitHub Enterprise account/repository access
- ✅ Supabase project access
- ✅ Stripe API keys (already have these)
- ✅ Supabase Service Role Key

---

## 🔑 Step 1: Get Your Supabase Access Token

You need a Supabase access token for GitHub Actions to authenticate:

1. **Go to:** https://supabase.com/dashboard/account/tokens
2. **Click:** "Generate new token"
3. **Name it:** `github-actions-deploy` (or any name)
4. **Copy the token** (starts with `sbp_...`)

⚠️ **Keep this token secret!** You'll add it to GitHub secrets.

---

## 🔐 Step 2: Set Up GitHub Enterprise Secrets

Go to your GitHub Enterprise repository and add these secrets:

### Navigate to Secrets:
1. Go to your repository on GitHub Enterprise
2. Click **Settings** → **Secrets and variables** → **Actions**
3. Click **"New repository secret"**

### Add These Secrets:

#### 1. `SUPABASE_ACCESS_TOKEN`
- **Value:** Your Supabase access token from Step 1 (starts with `sbp_...`)
- **Purpose:** Allows GitHub Actions to authenticate with Supabase

#### 2. `SUPABASE_PROJECT_REF`
- **Value:** `wbpfuuiznsmysbskywdx`
- **Purpose:** Identifies your Supabase project

#### 3. `SUPABASE_URL`
- **Value:** `https://wbpfuuiznsmysbskywdx.supabase.co`
- **Purpose:** Your Supabase project URL

#### 4. `SUPABASE_SERVICE_ROLE_KEY`
- **Value:** Your service role key from Supabase Dashboard
- **How to get:** Supabase Dashboard → Project Settings → API → Copy "service_role" key
- **Purpose:** Allows the function to bypass RLS and update user data

#### 5. `STRIPE_SECRET_KEY`
- **Value:** `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`
- **Purpose:** Stripe API authentication for webhook verification

#### 6. `STRIPE_WEBHOOK_SECRET`
- **Value:** Your webhook signing secret from Stripe (starts with `whsec_...`)
- **How to get:** 
  1. Deploy function first (or use a placeholder)
  2. Create webhook endpoint in Stripe Dashboard
  3. Copy the "Signing secret"
- **Purpose:** Verifies webhook requests are from Stripe

---

## 📁 Step 3: Commit the Workflow File

The GitHub Actions workflow file is already created at:
```
.github/workflows/deploy-supabase-functions.yml
```

### ⚠️ Important: Before Committing

**DO NOT commit files with actual secrets!** The file `SETUP_WEBHOOK_NOW.md` contains your actual Stripe keys. Either:
- **Option A:** Remove the keys from that file before committing
- **Option B:** Add it to `.gitignore` (already done)
- **Option C:** Don't commit that file at all

### Commit and Push:

```bash
# Add the workflow file
git add .github/workflows/deploy-supabase-functions.yml

# Add .gitignore (to protect secrets)
git add .gitignore

# Commit
git commit -m "Add GitHub Actions workflow for Supabase Edge Functions deployment"

# Push to GitHub Enterprise
git push origin main
```

**Or if you haven't initialized git yet:**

```bash
cd C:\Users\puppiesandkittens\Downloads\mmcs

# Initialize git (if needed)
git init

# Add remote (replace with your GitHub Enterprise URL)
git remote add origin https://github-enterprise.yourcompany.com/your-org/mmcs.git

# Add all files
git add .

# Commit
git commit -m "Initial commit with Supabase Edge Function deployment"

# Push
git push -u origin main
```

---

## 🚀 Step 4: Trigger Deployment

### Automatic Deployment:
The workflow runs automatically when you:
- ✅ Push changes to `supabase/functions/**` files
- ✅ Push changes to the workflow file itself
- ✅ Push to `main` or `master` branch

### Manual Deployment:
1. Go to your repository on GitHub Enterprise
2. Click **Actions** tab
3. Select **"Deploy Supabase Edge Functions"** workflow
4. Click **"Run workflow"** → **"Run workflow"**

---

## ✅ Step 5: Verify Deployment

### Check GitHub Actions:
1. Go to **Actions** tab in your repository
2. Click on the latest workflow run
3. Verify all steps completed successfully ✅

### Get Your Webhook URL:
After deployment, your webhook URL will be:
```
https://wbpfuuiznsmysbskywdx.supabase.co/functions/v1/stripe-webhook
```

### Configure Stripe Webhook:
1. Go to: https://dashboard.stripe.com/webhooks
2. Click **"+ Add endpoint"**
3. **Endpoint URL:** `https://wbpfuuiznsmysbskywdx.supabase.co/functions/v1/stripe-webhook`
4. **Events to send:**
   - ✅ `checkout.session.completed`
   - ✅ `customer.subscription.updated`
   - ✅ `customer.subscription.deleted`
5. Click **"Add endpoint"**
6. **Copy the "Signing secret"** (starts with `whsec_...`)
7. **Update GitHub secret:** Go back to GitHub → Settings → Secrets → Update `STRIPE_WEBHOOK_SECRET`
8. **Redeploy:** Either push a change or manually trigger the workflow

---

## 🔄 Step 6: Update Webhook Secret (After First Deployment)

After you get the webhook secret from Stripe:

1. **Update GitHub Secret:**
   - Go to: Repository → Settings → Secrets → Actions
   - Click on `STRIPE_WEBHOOK_SECRET`
   - Click **"Update"**
   - Paste your `whsec_...` secret
   - Click **"Update secret"**

2. **Redeploy Function:**
   - Go to: Actions tab
   - Click **"Deploy Supabase Edge Functions"**
   - Click **"Run workflow"** → **"Run workflow"**

---

## 🧪 Testing

### Test the Deployment:

1. **Complete a test subscription** using Stripe test mode
2. **Check Stripe Dashboard** → Webhooks → Recent events
3. **Check Supabase** → `unlocked_users` table
4. **Verify user upgraded** automatically

### View Function Logs:

```bash
# If you have Supabase CLI installed locally:
supabase functions logs stripe-webhook

# Or check in Supabase Dashboard:
# Dashboard → Edge Functions → stripe-webhook → Logs
```

---

## 📝 Workflow Details

### What the Workflow Does:

1. **Checks out your code** from GitHub
2. **Sets up Supabase CLI** in the GitHub Actions runner
3. **Logs in to Supabase** using your access token
4. **Links to your project** using project ref
5. **Sets all secrets** (Stripe keys, Supabase keys, etc.)
6. **Deploys the function** (`stripe-webhook`)
7. **Verifies deployment** and outputs the webhook URL

### When It Runs:

- ✅ **On push** to `main`/`master` branch (if `supabase/functions/**` changed)
- ✅ **Manually** via "Run workflow" button
- ✅ **On workflow file changes**

### What Files Are Deployed:

- `supabase/functions/stripe-webhook/index.ts` → Deployed as Edge Function

---

## 🚨 Troubleshooting

### "Authentication failed" or "Invalid token"
- ✅ Verify `SUPABASE_ACCESS_TOKEN` is correct
- ✅ Make sure token hasn't expired
- ✅ Check token has proper permissions

### "Project not found"
- ✅ Verify `SUPABASE_PROJECT_REF` matches your project
- ✅ Check project ref in Supabase Dashboard URL

### "Secret not found" errors
- ✅ Verify all 6 secrets are set in GitHub
- ✅ Check secret names match exactly (case-sensitive)
- ✅ Make sure secrets don't have extra spaces

### Function deploys but webhook doesn't work
- ✅ Verify `STRIPE_WEBHOOK_SECRET` is set correctly
- ✅ Check Stripe webhook endpoint URL matches deployment URL
- ✅ Verify webhook events are configured in Stripe
- ✅ Check function logs for errors

### "Workflow not running"
- ✅ Check you're pushing to `main` or `master` branch
- ✅ Verify workflow file is in `.github/workflows/` directory
- ✅ Check file is named correctly: `deploy-supabase-functions.yml`
- ✅ Verify GitHub Actions is enabled for your repository

---

## 🔒 Security Best Practices

1. **Never commit secrets** to your repository
2. **Use GitHub Secrets** for all sensitive values
3. **Rotate tokens regularly** (Supabase access tokens, Stripe keys)
4. **Limit access** to repository secrets (only trusted collaborators)
5. **Monitor workflow runs** for unauthorized access
6. **Use branch protection** to prevent accidental deployments

---

## 📚 Additional Resources

- **Supabase CLI Docs:** https://supabase.com/docs/guides/cli
- **GitHub Actions Docs:** https://docs.github.com/en/actions
- **Stripe Webhooks:** https://stripe.com/docs/webhooks
- **Supabase Edge Functions:** https://supabase.com/docs/guides/functions

---

## ✅ Quick Checklist

- [ ] Supabase access token generated
- [ ] All 6 GitHub secrets added
- [ ] Workflow file committed to repository
- [ ] Code pushed to GitHub Enterprise
- [ ] Workflow runs successfully
- [ ] Function deployed and accessible
- [ ] Stripe webhook endpoint created
- [ ] Webhook secret added to GitHub secrets
- [ ] Function redeployed with webhook secret
- [ ] Test subscription completed
- [ ] User automatically upgraded in Supabase

---

**That's it!** Your Supabase Edge Function will now deploy automatically whenever you push changes to `supabase/functions/` or manually trigger the workflow. 🎉
