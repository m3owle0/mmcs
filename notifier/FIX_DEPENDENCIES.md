# Fix Missing Go Dependencies

## The Problem

Go is trying to use a proxy at `127.0.0.1:9` which isn't working. You need to either:
1. Disable the proxy, or
2. Set it to use Go's public proxy directly

## Solution

Run these commands in PowerShell:

```powershell
# Disable proxy (use direct connection)
$env:GOPROXY="direct"

# Or use Go's public proxy directly
$env:GOPROXY="https://proxy.golang.org,direct"

# Then download dependencies
cd C:\Users\puppiesandkittens\Downloads\mmcs\notifier
go mod tidy
```

## Alternative: Check Your Go Proxy Settings

If you have a `GOPROXY` environment variable set, check it:

```powershell
echo $env:GOPROXY
```

If it shows something like `http://127.0.0.1:9`, that's the problem. Unset it:

```powershell
$env:GOPROXY=""
go mod tidy
```

## After Dependencies Are Downloaded

Once `go mod tidy` succeeds, you should see a `go.sum` file created. Then you can run:

```powershell
$env:SUPABASE_ANON_KEY="your_key_here"
go run .
```

## Quick Fix Script

Create a file `setup.ps1`:

```powershell
# Fix proxy and download dependencies
$env:GOPROXY="https://proxy.golang.org,direct"
cd C:\Users\puppiesandkittens\Downloads\mmcs\notifier
go mod tidy
Write-Host "Dependencies downloaded! Now run: go run ."
```

Then run: `.\setup.ps1`
