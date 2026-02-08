# How to Find Your Supabase API Key (2026)

## Quick Steps

1. **Go to:** https://supabase.com/dashboard
2. **Select your project** (click on it)
3. **Click one of these:**
   - **"Project Settings"** (gear icon ⚙️ in left sidebar)
   - **OR** the **"Connect"** button (usually at top)
4. **Go to "API Keys" section**
5. **Copy the key:**
   - **"Publishable key"** (new format, starts with `sb_publishable_...`) ← **Use this if available**
   - **OR** Legacy **"anon"** key (if you see a "Legacy API Keys" tab)

## Visual Guide

```
Supabase Dashboard
  └─ Select Your Project
      └─ Project Settings (⚙️) OR Connect Button
          └─ API Keys Section
              ├─ Publishable key (sb_publishable_xxx) ← Use this
              └─ OR Legacy API Keys tab → anon key
```

## Important Notes

- ✅ **Use:** Publishable key (`sb_publishable_...`) or Legacy anon key
- ❌ **Don't use:** service_role key (has admin privileges, unsafe for client-side)
- 🔑 **Key format:** New keys start with `sb_publishable_`, old keys are shorter

## If You Can't Find It

1. Make sure you're logged into the correct Supabase account
2. Make sure you've selected the correct project
3. Try the "Connect" button - it often shows the key directly
4. Check if you have the right permissions (project owner/admin)

## Your Current Key Format

If your key starts with `sb_publishable_`, you're using the new format (correct!).
