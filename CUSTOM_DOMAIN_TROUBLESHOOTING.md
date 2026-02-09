# Custom Domain Troubleshooting for multimarketclothingsearch.com

Your GitHub Pages site works at: `https://m3owle0.github.io/mmcs/` ✅

But `https://multimarketclothingsearch.com` shows Hostinger parking page ❌

## Step-by-Step Fix

### 1. Verify CNAME File Location in GitHub

The `CNAME` file **must** be in the same location as your `index.html`:

**If your repo structure is:**
```
mmcs/
  ├── index.html
  ├── miku.svg
  ├── mmcs.svg
  ├── bg.svg
  └── CNAME  ← Must be here!
```

**Check:**
- Go to: `https://github.com/m3owle0/mmcs/tree/main/mmcs` (or wherever your files are)
- Verify `CNAME` file exists in the same folder as `index.html`
- The `CNAME` file should contain: `multimarketclothingsearch.com`

### 2. Verify GitHub Pages Custom Domain Settings

1. Go to: `https://github.com/m3owle0/mmcs/settings/pages`
2. Under "Custom domain", it should show: `multimarketclothingsearch.com`
3. Status should show: ✅ DNS check successful
4. "Enforce HTTPS" should be checked

### 3. Disable Hostinger Hosting (CRITICAL)

If Hostinger is serving a parking page, you need to disable hosting:

**In Hostinger:**
1. Go to your Hostinger account dashboard
2. Find "Websites" or "Hosting" section
3. Look for `multimarketclothingsearch.com`
4. **Disable/Remove hosting** for this domain
5. **Keep DNS management active** (you still need DNS records)

**Why:** Hostinger hosting can override DNS and serve the parking page even if DNS points to GitHub.

### 4. Verify DNS Records (Again)

In Hostinger DNS management, confirm:

**A Records (4 total):**
- `@` → `185.199.108.153`
- `@` → `185.199.109.153`
- `@` → `185.199.110.153`
- `@` → `185.199.111.153`

**CNAME Record:**
- `www` → `m3owle0.github.io`

**Make sure:**
- No A records pointing to Hostinger IPs (like `84.32.84.32`)
- No conflicting records

### 5. Check DNS Propagation

Visit: https://dnschecker.org/#A/multimarketclothingsearch.com

**You should see:**
- Most locations showing: `185.199.108.153`, `185.199.109.153`, `185.199.110.153`, `185.199.111.153`

**If you see:**
- Different IPs → DNS hasn't propagated yet (wait 24-48 hours)
- Hostinger IPs → Hostinger hosting is still active (disable it)

### 6. Clear DNS Cache

**On your computer:**
```bash
# Windows (run as Administrator)
ipconfig /flushdns

# Mac/Linux
sudo dscacheutil -flushcache
# or
sudo systemd-resolve --flush-caches
```

**Or use a different DNS:**
- Try accessing from your phone (using mobile data, not WiFi)
- Use a VPN
- Use: https://www.whatsmydns.net/#A/multimarketclothingsearch.com

### 7. Wait for Propagation

DNS changes can take:
- **Minimum:** 5-30 minutes
- **Typical:** 1-4 hours
- **Maximum:** 24-48 hours

## Most Likely Issue

**Hostinger hosting is still active** and serving the parking page, even though DNS points to GitHub.

**Solution:** Disable hosting in Hostinger, keep DNS management active.

## Quick Test

After disabling Hostinger hosting:

1. Wait 10-30 minutes
2. Try: `https://multimarketclothingsearch.com` (use HTTPS, not HTTP)
3. Check: `https://www.multimarketclothingsearch.com`

## If Still Not Working

1. **Double-check CNAME file location** in GitHub repo
2. **Verify GitHub Pages custom domain** is set correctly
3. **Check DNS propagation** at dnschecker.org
4. **Contact Hostinger support** to ensure hosting is fully disabled
5. **Wait 24-48 hours** for full DNS propagation

---

**Remember:** The domain should work at the **root** (`multimarketclothingsearch.com`), not at `/mmcs`. GitHub Pages custom domains always serve from root, not subdirectories.
