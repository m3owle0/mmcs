# Free Cloud Hosting for MMCS Notifier

Your notifier needs:
- ✅ Environment variables (SUPABASE_SERVICE_ROLE_KEY)
- ✅ Continuous running (polls every 1 minute)
- ✅ HTTP requests (Supabase API, Sendico API, Discord webhooks)
- ✅ Go runtime support

## Best Free Options

### Option 1: Railway (Recommended) ⭐

**Why:** Best free tier, supports continuous running, easy deployment

**Free Tier:**
- $5 free credit/month (usually enough for small apps)
- 500 hours runtime/month
- Environment variables supported
- Auto-deploy from GitHub

**Steps:**

1. **Sign up:** https://railway.app (use GitHub to sign in)

2. **Create new project:**
   - Click "New Project"
   - Select "Deploy from GitHub repo"
   - Choose your `mmcs` repository
   - Select the `notifier` folder

3. **Configure:**
   - Railway will auto-detect Go
   - Set environment variable:
     - Go to project → Variables
     - Add: `SUPABASE_SERVICE_ROLE_KEY` = `your-service-role-key`
   - (Optional) Add: `SUPABASE_ANON_KEY` as fallback

4. **Deploy:**
   - Railway will build and deploy automatically
   - Check logs to verify it's running

**Cost:** Free (within $5 credit limit)

---

### Option 2: Render

**Why:** Free tier available, but spins down after inactivity

**Free Tier:**
- 750 hours/month
- Spins down after 15 minutes of inactivity
- Wakes up on first request (but your notifier needs to run continuously)

**Steps:**

1. **Sign up:** https://render.com

2. **Create Web Service:**
   - New → Web Service
   - Connect GitHub repo
   - Select `notifier` folder

3. **Configure:**
   - **Build Command:** `go mod download && go build -o notifier`
   - **Start Command:** `./notifier`
   - **Environment Variables:**
     - `SUPABASE_SERVICE_ROLE_KEY` = `your-key`

4. **Important:** Free tier spins down, so this might not work well for continuous polling

**Cost:** Free (but may not work for continuous running)

---

### Option 3: Fly.io ⭐

**Why:** Great for Go apps, generous free tier

**Free Tier:**
- 3 shared VMs
- 3GB persistent storage
- Continuous running supported

**Steps:**

1. **Install Fly CLI:**
   ```bash
   # Windows (PowerShell)
   iwr https://fly.io/install.ps1 -useb | iex
   
   # Mac/Linux
   curl -L https://fly.io/install.sh | sh
   ```

2. **Login:**
   ```bash
   fly auth login
   ```

3. **Create app:**
   ```bash
   cd notifier
   fly launch
   ```
   - Follow prompts
   - Don't deploy yet

4. **Create `fly.toml`** (in `notifier` folder):
   ```toml
   app = "mmcs-notifier"
   primary_region = "iad"  # Choose closest region
   
   [build]
     builder = "paketobuildpacks/builder:base"
   
   [env]
     SUPABASE_SERVICE_ROLE_KEY = "your-key-here"
   
   [[services]]
     internal_port = 8080
     protocol = "tcp"
   ```

5. **Deploy:**
   ```bash
   fly deploy
   ```

**Cost:** Free (within limits)

---

### Option 4: GitHub Actions (Cron Job) ⭐⭐

**Why:** Completely free, runs on schedule, no hosting needed

**How it works:** Runs your notifier as a scheduled job every minute

**Steps:**

1. **Create `.github/workflows/notifier.yml`** in your repo root:
   ```yaml
   name: MMCS Notifier
   
   on:
     schedule:
       - cron: '*/1 * * * *'  # Every minute
     workflow_dispatch:  # Manual trigger
   
   jobs:
     notify:
       runs-on: ubuntu-latest
       steps:
         - uses: actions/checkout@v3
         
         - name: Set up Go
           uses: actions/setup-go@v4
           with:
             go-version: '1.23'
         
         - name: Download dependencies
           working-directory: ./notifier
           run: go mod download
         
         - name: Build notifier
           working-directory: ./notifier
           run: go build -o notifier
         
         - name: Run notifier
           working-directory: ./notifier
           env:
             SUPABASE_SERVICE_ROLE_KEY: ${{ secrets.SUPABASE_SERVICE_ROLE_KEY }}
           run: ./notifier
   ```

2. **Add secret to GitHub:**
   - Repo → Settings → Secrets and variables → Actions
   - New repository secret
   - Name: `SUPABASE_SERVICE_ROLE_KEY`
   - Value: Your service role key

3. **Note:** GitHub Actions has a limit of 1 minute minimum between runs, so this works perfectly!

**Cost:** Completely FREE (2,000 minutes/month free for private repos, unlimited for public)

---

### Option 5: Google Cloud Run (Serverless)

**Why:** Free tier, but needs adaptation for continuous running

**Free Tier:**
- 2 million requests/month
- 360,000 GB-seconds compute
- 180,000 vCPU-seconds

**Note:** Cloud Run is serverless (runs on demand), so you'd need to adapt the code or use Cloud Scheduler to ping it every minute.

**Cost:** Free (within limits)

---

## Recommendation

**Best Option:** **GitHub Actions** (Option 4)
- ✅ Completely free
- ✅ No hosting setup needed
- ✅ Runs every minute (perfect for your use case)
- ✅ Easy to configure
- ✅ Automatic from GitHub

**Second Best:** **Railway** (Option 1)
- ✅ Easy setup
- ✅ Continuous running
- ✅ $5 free credit/month (usually enough)

---

## Quick Start: GitHub Actions (Recommended)

1. Create `.github/workflows/notifier.yml` in your repo
2. Add `SUPABASE_SERVICE_ROLE_KEY` as a GitHub secret
3. Push to GitHub
4. Done! It will run every minute automatically

Want me to create the GitHub Actions workflow file for you?
