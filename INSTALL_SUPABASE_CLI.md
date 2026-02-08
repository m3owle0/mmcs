# How to Install Supabase CLI on Windows

## ⚠️ Important: npm install -g Does NOT Work!

The Supabase CLI **cannot** be installed via `npm install -g supabase`. You'll get this error:
```
Installing Supabase CLI as a global module is not supported.
```

---

## ✅ Method 1: Direct Download (Easiest - No Dependencies)

### Step 1: Download
1. Go to: https://github.com/supabase/cli/releases/latest
2. Download: `supabase_windows_amd64.zip` (for most Windows PCs)
   - If you have ARM processor, download `supabase_windows_arm64.zip`

### Step 2: Extract
1. Extract the zip file
2. You'll get `supabase.exe`

### Step 3: Add to PATH (Optional but Recommended)

**Option A: Add to System PATH**
1. Copy `supabase.exe` to `C:\Windows\System32`
2. Now you can run `supabase` from anywhere

**Option B: Create Tools Folder**
1. Create folder: `C:\Tools`
2. Copy `supabase.exe` to `C:\Tools`
3. Add to PATH:
   - Press `Win + X` → **System** → **Advanced system settings**
   - Click **"Environment Variables"**
   - Under **"User variables"**, find `Path` → **Edit**
   - Click **"New"** → Add `C:\Tools`
   - Click **OK** on all dialogs

**Option C: Just Use Full Path**
- You can run it directly: `C:\path\to\supabase.exe login`

### Step 4: Verify Installation
Open PowerShell and run:
```powershell
supabase --version
```

You should see the version number.

---

## ✅ Method 2: Scoop (Package Manager - Recommended)

### Step 1: Install Scoop (If You Don't Have It)

Open PowerShell and run:
```powershell
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
irm get.scoop.sh | iex
```

### Step 2: Install Supabase CLI

```powershell
scoop bucket add supabase https://github.com/supabase/scoop-bucket.git
scoop install supabase
```

### Step 3: Verify
```powershell
supabase --version
```

---

## ✅ Method 3: Use npx (No Installation - Temporary)

You can use Supabase CLI **without installing** by using `npx`:

```bash
# Login
npx supabase@latest login

# Link project
npx supabase@latest link --project-ref wbpfuuiznsmysbskywdx

# Deploy function
npx supabase@latest functions deploy stripe-webhook
```

**Note:** This downloads the CLI each time, but works if you can't install it.

---

## ✅ Method 4: Chocolatey (If You Have It)

```powershell
choco install supabase
```

---

## 🎯 Recommended: Method 1 (Direct Download)

**Why?**
- ✅ No dependencies needed
- ✅ Works immediately
- ✅ No package manager setup
- ✅ Simple and reliable

**Steps:**
1. Download from GitHub releases
2. Extract `supabase.exe`
3. Copy to `C:\Windows\System32` (or add to PATH)
4. Done!

---

## 🧪 Test Installation

After installing, test it:

```powershell
supabase --version
supabase --help
```

You should see version info and help text.

---

## 📚 Next Steps

Once Supabase CLI is installed, continue with the webhook setup:

1. `supabase login`
2. `supabase link --project-ref wbpfuuiznsmysbskywdx`
3. `supabase functions deploy stripe-webhook`

See `WEBHOOK_QUICK_START.md` for full instructions.

---

## 🆘 Troubleshooting

### "supabase is not recognized"
- Make sure `supabase.exe` is in your PATH
- Or use full path: `C:\path\to\supabase.exe`
- Or use `npx supabase@latest` instead

### "Permission denied"
- Run PowerShell as Administrator
- Or copy `supabase.exe` to a folder you own (not System32)

### Still having issues?
Use **Method 3 (npx)** - it works without any installation!
