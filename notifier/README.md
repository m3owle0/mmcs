# Discord Notifier - Quick Start

## ✅ 5 Markets Working!

The notifier now has **real search functionality** for 5 Japanese markets via Sendico:
- mercari-jp (Mercari Japan)
- paypay-fleamarket (Yahoo PayPay Flea)
- rakuma (Rakuten Rakuma)
- rakuten-jp (Rakuten)
- yahoo-auctions (Yahoo Auctions)

## Setup (2 minutes)

1. **Install Go:** https://go.dev/dl/

2. **Get Supabase Key:**
   - Go to: https://supabase.com/dashboard
   - Select your project
   - Click "Project Settings" (gear icon) or "Connect" button
   - Go to "API Keys" section
   - Copy the **"Publishable key"** (starts with `sb_publishable_...`) OR the Legacy **"anon"** key

3. **Install Dependencies:**
   ```powershell
   go mod download
   ```

4. **Edit `start.bat`:**
   - Open `start.bat` in notepad
   - Replace `YOUR_KEY_HERE` with your Supabase key

5. **Run:**
   ```powershell
   .\start.bat
   ```

## Or Run Manually

```powershell
$env:SUPABASE_ANON_KEY="your_key_here"
go run .
```

## Build Executable

```powershell
go build -o discord-notifier.exe .
$env:SUPABASE_ANON_KEY="your_key"
.\discord-notifier.exe
```

## What It Does

- ✅ Checks Supabase every 5 minutes
- ✅ Finds users with active subscriptions
- ✅ **Searches Sendico API** for the 5 supported markets
- ✅ **Auto-translates** English → Japanese
- ✅ **Tracks seen items** to avoid duplicates
- ✅ Sends **real notifications** with item data to Discord

## Features

- Real searches (not test data!)
- Automatic translation
- Item deduplication
- Rich Discord embeds with images

## See Also

- `SENDICO_INTEGRATION.md` - Details about the Sendico integration
- `BACKEND_SETUP.md` - Deployment options
