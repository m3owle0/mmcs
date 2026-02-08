# 🔐 Add Secrets to GitHub - Step by Step

You have your keys! Now add them to GitHub:

## ✅ Step 1: Go to GitHub Secrets

1. **Go to your GitHub repository**
2. **Click:** Settings (top menu)
3. **Click:** Secrets and variables → Actions
4. **Click:** "New repository secret" button

---

## ✅ Step 2: Add Each Secret (One at a Time)

Add these 5 secrets:

### Secret 1: SUPABASE_ACCESS_TOKEN
- **Name:** `SUPABASE_ACCESS_TOKEN`
- **Value:** `sbp_1bd867a6e2ef2c7ce2bf69d14811475560cb787f`
- Click "Add secret"

### Secret 2: SUPABASE_PROJECT_REF
- **Name:** `SUPABASE_PROJECT_REF`
- **Value:** `wbpfuuiznsmysbskywdx`
- Click "Add secret"

### Secret 3: SUPABASE_URL
- **Name:** `SUPABASE_URL`
- **Value:** `https://wbpfuuiznsmysbskywdx.supabase.co`
- Click "Add secret"

### Secret 4: SUPABASE_SERVICE_ROLE_KEY
- **Name:** `SUPABASE_SERVICE_ROLE_KEY`
- **Value:** `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6IndicGZ1dWl6bnNteXNic2t5d2R4Iiwicm9sZSI6InNlcnZpY2Vfcm9sZSIsImlhdCI6MTc3MDE3NTIyMywiZXhwIjoyMDg1NzUxMjIzfQ.SV6k5LEDxCEh3cjwkfcTvA0lan1Lf8rReo5bGLVBgms`
- Click "Add secret"

### Secret 5: STRIPE_SECRET_KEY
- **Name:** `STRIPE_SECRET_KEY`
- **Value:** `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`
- Click "Add secret"

---

## ✅ Step 3: Deploy the Function

After adding all 5 secrets:

1. **Go to:** Actions tab (in your GitHub repo)
2. **Click:** "Deploy Supabase Edge Functions" workflow
3. **Click:** "Run workflow" button (top right)
4. **Click:** "Run workflow" again to confirm
5. **Wait 1-2 minutes** for it to complete

---

## ✅ Step 4: Check if It Worked

1. **Look at the workflow run** - all steps should have green checkmarks ✅
2. **If it failed:** Click on the failed step to see the error
3. **If it succeeded:** You'll see "✅ Function deployed successfully!"

---

## ⚠️ Note About STRIPE_WEBHOOK_SECRET

You'll add this **after** you create the Stripe webhook endpoint (Step 5 in the original instructions). For now, the function can deploy without it, but the webhook won't work until you add it.

---

**Once the workflow completes successfully, let me know and we'll set up the Stripe webhook!**
