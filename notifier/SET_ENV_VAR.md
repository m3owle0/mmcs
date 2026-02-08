# Setting SUPABASE_ANON_KEY in PowerShell

## The Problem

You had a syntax error in your PowerShell command. There was a **space** after `$env:` which is incorrect.

## Correct Syntax

### Option 1: PowerShell (No space after $env:)

```powershell
$env:SUPABASE_ANON_KEY="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6IndicGZ1dWl6bnNteXNic2t5d2R4Iiwicm9sZSI6ImFub24iLCJpYXQiOjE3NzAxNzUyMjMsImV4cCI6MjA4NTc1MTIyM30.t48b38QU8QpWfDyGu__hTKdCYbjVh1rhHcrt1D7mFWU"
```

**Important:** No space between `$env:` and `SUPABASE_ANON_KEY`

### Option 2: Use start.bat (Easier)

1. Open `start.bat` in notepad
2. Find the line: `set SUPABASE_ANON_KEY=YOUR_KEY_HERE`
3. Replace `YOUR_KEY_HERE` with your actual key
4. Save and run: `.\start.bat`

### Option 3: Set in PowerShell Session

```powershell
# Set the variable (no space after $env:)
$env:SUPABASE_ANON_KEY="your_complete_key_here"

# Verify it's set
echo $env:SUPABASE_ANON_KEY

# Then run the notifier
go run .
```

## Common Mistakes

❌ **Wrong:** `$env: SUPABASE_ANON_KEY="..."` (space after $env:)  
✅ **Correct:** `$env:SUPABASE_ANON_KEY="..."` (no space)

❌ **Wrong:** Key with spaces or line breaks  
✅ **Correct:** Complete key in one string, no spaces

## Your Key (from the error)

It looks like your key might be truncated or have a space. Make sure you copy the **complete** key from Supabase without any spaces or line breaks.

The key should look like:
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6IndicGZ1dWl6bnNteXNic2t5d2R4Iiwicm9sZSI6ImFub24iLCJpYXQiOjE3NzAxNzUyMjMsImV4cCI6MjA4NTc1MTIyM30.t48b38QU8QpWfDyGu__hTKdCYbjVh1rhHcrt1D7mFWU
```

(All one continuous string, no spaces)
