# 🔑 Finding Your Supabase Service Role Key

You're on the right page! Here's exactly where to find it:

## 📍 Current Page:
**Supabase Dashboard → Project Settings → API**
URL: `https://supabase.com/dashboard/project/wbpfuuiznsmysbskywdx/settings/api`

## 🔍 What to Do:

1. **Scroll down** on the page you're currently viewing
2. **Look for a section called:** "Project API keys" or "API Keys"
3. **You'll see two keys:**
   - **anon/public** key (starts with `eyJ...`) - ❌ NOT this one
   - **service_role** key (also starts with `eyJ...`) - ✅ **THIS ONE!**

## ⚠️ Important:
- The **service_role** key is the one you need
- It's usually shown with a warning icon (⚠️) because it has full access
- It's longer than the anon key
- **Copy the entire key** - it's a long string

## 📋 What It Looks Like:

```
Project API keys

anon / public
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6IndicGZ1dWl6bnNteXNic2t5d2R4Iiwicm9sZSI6ImFub24iLCJpYXQiOjE3MDk...

service_role  ⚠️
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6IndicGZ1dWl6bnNteXNic2t5d2R4Iiwicm9sZSI6InNlcnZpY2Vfcm9sZSIsImlhdCI6MTcwOT...  [Copy]
```

## ✅ Once You Find It:

1. **Click the "Copy" button** next to `service_role`
2. **Paste it** into GitHub Secrets as `SUPABASE_SERVICE_ROLE_KEY`
3. **Continue with Step 3** of the setup instructions

---

**Can't find it?** It might be collapsed or you need to click "Reveal" to show it.
