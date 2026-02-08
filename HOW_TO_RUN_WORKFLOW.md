# 🚀 How to Run GitHub Actions Workflow

## Step-by-Step Instructions:

### Step 1: Go to Your GitHub Repository

1. **Open your browser**
2. **Go to:** Your GitHub repository (the one with your code)
3. **Make sure you're logged in**

---

### Step 2: Open the Actions Tab

1. **Click:** "Actions" tab (top menu, next to "Code", "Issues", etc.)
2. **You'll see:** A list of workflows (or "No workflow runs" if none have run yet)

---

### Step 3: Find Your Workflow

1. **Look for:** "Deploy Supabase Edge Functions" workflow
2. **If you see it:** Click on it
3. **If you don't see it:** 
   - Make sure the workflow file exists: `.github/workflows/deploy-supabase-functions.yml`
   - If it doesn't exist, you need to commit and push it first

---

### Step 4: Run the Workflow

1. **Click:** "Run workflow" button (top right, blue button)
2. **Select branch:** Make sure it says `main` or `master` (your default branch)
3. **Click:** "Run workflow" button again (green button that appears)
4. **Wait:** You'll see a yellow dot appear - the workflow is starting

---

### Step 5: Watch It Run

1. **You'll see:** A new workflow run appear at the top of the list
2. **Click on it** to see the progress
3. **Watch the steps:**
   - ✅ Green checkmark = Success
   - ⏳ Yellow circle = Running
   - ❌ Red X = Failed

---

### Step 6: Check Results

**If successful:**
- ✅ All steps will have green checkmarks
- ✅ You'll see: "✅ Function deployed successfully!"
- ✅ Your Edge Function is now deployed!

**If failed:**
- ❌ One or more steps will have red X marks
- **Click on the failed step** to see the error
- **Common errors:**
  - Missing secrets → Add them in Settings → Secrets
  - Invalid token → Regenerate Supabase access token
  - Wrong project ref → Should be `wbpfuuiznsmysbskywdx`

---

## 📸 Visual Guide:

```
GitHub Repository
├── Code tab
├── Issues tab
├── Pull requests tab
├── Actions tab ← CLICK HERE
│   ├── "Deploy Supabase Edge Functions" ← CLICK THIS
│   │   └── "Run workflow" button ← CLICK THIS
│   │       └── "Run workflow" (confirm) ← CLICK THIS
│   └── [Workflow runs appear here]
└── Settings tab
```

---

## ⚠️ Before Running:

**Make sure you have:**
- ✅ Added all 5 GitHub secrets (see `ADD_GITHUB_SECRETS.md`)
- ✅ Workflow file exists (`.github/workflows/deploy-supabase-functions.yml`)

---

## 🔄 Automatic Runs:

**The workflow will also run automatically when:**
- ✅ You push changes to `supabase/functions/**` files
- ✅ You push changes to the workflow file itself
- ✅ You push to `main` or `master` branch

**But you can always run it manually using the steps above!**

---

## 🚨 Troubleshooting:

### "Run workflow" button is grayed out or missing?
- **Check:** Are you on the `main` or `master` branch?
- **Check:** Does the workflow file exist?
- **Check:** Do you have permission to run workflows?

### Workflow runs but fails immediately?
- **Check:** Are all secrets added?
- **Check:** Are secret names spelled correctly?
- **Check:** Do secrets have the correct values?

### Can't find the Actions tab?
- **Check:** You're in the right repository
- **Check:** You're logged into GitHub
- **Check:** You have access to the repository

---

**That's it!** Once you click "Run workflow", it will deploy your Edge Function automatically. 🎉
