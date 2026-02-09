# Step-by-Step Setup Guide

Follow these steps in order to get your MMCS website up and running.

## Prerequisites Checklist
- [ ] A Supabase account (free tier works)
- [ ] A Discord account (for testing notifications)
- [ ] Go 1.23+ installed (for the notifier service)
- [ ] A text editor or IDE

---

## Step 1: Set Up Supabase Database

### 1.1 Create Supabase Project
1. Go to [supabase.com](https://supabase.com)
2. Sign up or log in
3. Click **"New Project"**
4. Fill in:
   - **Name**: MMCS (or any name)
   - **Database Password**: Create a strong password (save it!)
   - **Region**: Choose closest to you
5. Click **"Create new project"**
6. Wait 2-3 minutes for project to initialize

### 1.2 Run Database Schema
1. In Supabase dashboard, click **"SQL Editor"** (left sidebar)
2. Click **"New query"**
3. Open `database/schema.sql` in your text editor
4. Copy **ALL** the SQL code from `database/schema.sql`
5. Paste it into the Supabase SQL Editor
6. Click **"Run"** (or press Ctrl+Enter)
7. You should see "Success. No rows returned"
8. ✅ Database is now set up!

### 1.3 Get Your API Keys
1. In Supabase dashboard, click **"Project Settings"** (gear icon, bottom left)
2. Click **"API"** in the settings menu
3. Find these values and **copy them** (you'll need them):
   - **Project URL**: `https://xxxxx.supabase.co`
   - **anon public key**: `eyJhbGc...` (long string)
   - **service_role key**: `eyJhbGc...` (long string, keep secret!)

---

## Step 2: Configure the Website

### 2.1 Update Supabase Credentials
1. Open `index.html` in your text editor
2. Press **Ctrl+F** (or Cmd+F on Mac) to search
3. Search for: `SUPABASE_URL = 'https://wbpfuuiznsmysbskywdx.supabase.co'`
4. Replace with your Project URL:
   ```javascript
   const SUPABASE_URL = 'https://YOUR-PROJECT-ID.supabase.co';
   ```
5. Search for: `SUPABASE_ANON_KEY = 'sb_publishable_rIy_-DWT87Gj9ao1WvN3gA_WA6eME-x'`
6. Replace with your anon key:
   ```javascript
   const SUPABASE_ANON_KEY = 'YOUR-ANON-KEY-HERE';
   ```
7. **Important**: There are **2 places** to update (search for both occurrences)
   - One around line 1715
   - One around line 2573
8. Save the file

### 2.2 Test the Website Locally
1. Open `index.html` in your web browser (double-click the file)
2. You should see the website
3. Try clicking **"Unlock Site"** button
4. Try creating an account (sign up)
5. ✅ Website is configured!

---

## Step 3: Set Up Discord Notifications (Optional)

### 3.1 Create Discord Webhook
1. Open Discord (web or app)
2. Go to your Discord server (or create a test server)
3. Click server name → **"Server Settings"**
4. Click **"Integrations"** → **"Webhooks"**
5. Click **"New Webhook"**
6. Give it a name: "MMCS Notifications"
7. Choose a channel (or create a test channel)
8. Click **"Copy Webhook URL"**
9. Save this URL somewhere (you'll need it later)

### 3.2 Set Up the Notifier Service

#### Option A: Quick Test (Windows)
1. Open Command Prompt or PowerShell
2. Navigate to the notifier folder:
   ```cmd
   cd "C:\Users\puppiesandkittens\Downloads\mmcs revamp\notifier"
   ```
3. Install dependencies:
   ```cmd
   go mod download
   ```
4. Set environment variable and run:
   ```cmd
   set SUPABASE_SERVICE_ROLE_KEY=your-service-role-key-here
   go run .
   ```
   (Replace `your-service-role-key-here` with your actual service role key)

#### Option B: Quick Test (Mac/Linux)
1. Open Terminal
2. Navigate to the notifier folder:
   ```bash
   cd ~/Downloads/mmcs\ revamp/notifier
   ```
3. Install dependencies:
   ```bash
   go mod download
   ```
4. Set environment variable and run:
   ```bash
   export SUPABASE_SERVICE_ROLE_KEY="your-service-role-key-here"
   go run .
   ```
   (Replace `your-service-role-key-here` with your actual service role key)

5. You should see:
   ```
   🚀 Starting Discord Notifier
   📡 Supabase URL: https://...
   ✅ Sendico client initialized
   ```
6. ✅ Notifier is running!

---

## Step 4: Verify a User Account

### 4.1 Create Test Account
1. Open `index.html` in browser
2. Click **"Unlock Site"**
3. Click **"Sign Up"**
4. Enter:
   - Username: `testuser`
   - Email: `your-email@example.com`
   - Password: `testpassword123`
5. Click **"Sign Up"**
6. Check your email for verification (if email confirmation is enabled)

### 4.2 Verify User in Database
1. Go to Supabase dashboard
2. Click **"Table Editor"** (left sidebar)
3. Click **"unlocked_users"** table
4. Find your user (by email)
5. Click the row to edit
6. Change **"verified"** from `false` to `true`
7. Click **"Save"** (or press Enter)
8. ✅ User is now verified!

### 4.3 Test Login
1. Refresh the website
2. Click **"Unlock Site"**
3. Enter your email and password
4. Click **"Login"**
5. You should see "Welcome, [username]" in the header
6. ✅ Login works!

---

## Step 5: Configure Notifications (User Side)

### 5.1 Add Discord Webhook
1. While logged in, click **"🔔 Notifications"** button (top right)
2. In the "Discord Webhook URL" field, paste your webhook URL
3. Click **"Save Webhook URL"**
4. You should see "✓ Webhook URL saved successfully!"
5. ✅ Webhook configured!

### 5.2 Activate Subscription (Admin Side)
1. Go to Supabase dashboard → **"SQL Editor"**
2. Run this SQL (replace email with your test email):
   ```sql
   UPDATE unlocked_users 
   SET notifications_subscription_active = true,
       notifications_subscription_expires_at = NULL
   WHERE email = 'your-email@example.com';
   ```
3. Click **"Run"**
4. ✅ Subscription activated!

### 5.3 Create a Notification
1. On the website, click **"🔔 Notifications"**
2. Click **"+ Add Notification"**
3. Enter search term: `Supreme hoodie` (or any term)
4. Select markets (or leave empty for all)
5. Click **"Save"**
6. You should see your notification in the list
7. ✅ Notification created!

### 5.4 Test Notifications
1. Make sure the notifier service is running (Step 3.2)
2. Wait 1-2 minutes
3. Check your Discord channel
4. You should receive notifications if new items are found
5. ✅ Notifications working!

---

## Step 6: Deploy Website (Optional)

### Option A: Netlify (Easiest)
1. Go to [netlify.com](https://netlify.com)
2. Sign up/login
3. Drag and drop your project folder onto Netlify
4. Your site will be live at `https://random-name.netlify.app`
5. ✅ Deployed!

### Option B: GitHub Pages
1. Create a GitHub repository
2. Upload your files
3. Go to Settings → Pages
4. Select main branch
5. Your site will be live at `https://username.github.io/repo-name`
6. ✅ Deployed!

### Option C: Keep Local
- Just open `index.html` in browser whenever you want to use it
- ✅ No deployment needed!

---

## Step 7: Run Notifier Continuously (Production)

### Windows (Task Scheduler)
1. Create a batch file `start-notifier.bat`:
   ```batch
   @echo off
   cd /d "C:\Users\puppiesandkittens\Downloads\mmcs revamp\notifier"
   set SUPABASE_SERVICE_ROLE_KEY=your-service-role-key
   notifier.exe
   ```
2. Build the notifier: `go build -o notifier.exe`
3. Set up Task Scheduler to run the batch file on startup

### Mac/Linux (Systemd Service)
1. Build the notifier:
   ```bash
   cd notifier
   go build -o notifier
   ```
2. Create `/etc/systemd/system/mmcs-notifier.service`:
   ```ini
   [Unit]
   Description=MMCS Discord Notifier
   After=network.target

   [Service]
   Type=simple
   User=your-username
   WorkingDirectory=/path/to/mmcs/notifier
   Environment="SUPABASE_SERVICE_ROLE_KEY=your-service-role-key"
   ExecStart=/path/to/mmcs/notifier/notifier
   Restart=always

   [Install]
   WantedBy=multi-user.target
   ```
3. Enable and start:
   ```bash
   sudo systemctl enable mmcs-notifier
   sudo systemctl start mmcs-notifier
   ```

---

## Troubleshooting

### Website shows "Supabase not configured"
- ✅ Check that you updated both SUPABASE_URL and SUPABASE_ANON_KEY locations
- ✅ Make sure you copied the keys correctly (no extra spaces)

### Can't log in after signup
- ✅ Go to Supabase → Table Editor → unlocked_users
- ✅ Find your user and set `verified = true`

### Notifications not working
- ✅ Check notifier is running (you should see logs)
- ✅ Verify webhook URL is saved in database
- ✅ Check `notifications_subscription_active = true` in database
- ✅ Check Discord webhook URL is correct format

### Notifier shows errors
- ✅ Verify service role key is correct
- ✅ Check Supabase URL is correct
- ✅ Make sure database schema was run successfully

### "Error loading notifications"
- ✅ Make sure you're logged in
- ✅ Check browser console (F12) for errors
- ✅ Verify RLS policies are set up correctly

---

## Quick Reference

### Important URLs
- Supabase Dashboard: https://app.supabase.com
- SQL Editor: Dashboard → SQL Editor
- Table Editor: Dashboard → Table Editor
- API Settings: Dashboard → Settings → API

### Important Files
- `index.html` - Main website (update Supabase keys here)
- `database/schema.sql` - Database setup (run in Supabase)
- `notifier/main.go` - Notification service
- `README.md` - Full documentation

### Key Values to Save
- ✅ Supabase Project URL
- ✅ Supabase Anon Key (public)
- ✅ Supabase Service Role Key (secret!)
- ✅ Discord Webhook URL

---

## You're Done! 🎉

Your website should now be:
- ✅ Running locally
- ✅ Connected to Supabase
- ✅ Ready for users to sign up
- ✅ Sending Discord notifications (if notifier is running)

**Next Steps:**
1. Verify more users as they sign up
2. Keep the notifier running for continuous notifications
3. Customize the website (colors, markets, etc.)
4. Deploy to a hosting service

Need help? Check `README.md` for more details!
