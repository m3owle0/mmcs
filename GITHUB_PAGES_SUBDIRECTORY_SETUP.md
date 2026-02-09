# GitHub Pages Custom Domain Setup for Subdirectory (/mmcs)

Your site is deployed at: `https://m3owle0.github.io/mmcs/`

## Important: CNAME File Location

Since your site is in the `/mmcs` subdirectory, the `CNAME` file **must be in the `mmcs` folder** in your GitHub repository.

### Steps:

1. **In your GitHub repository**, make sure the `CNAME` file is located at:
   ```
   mmcs/CNAME
   ```
   (Not in the root of the repo, but inside the `mmcs` folder)

2. **The `CNAME` file should contain:**
   ```
   multimarketclothingsearch.com
   ```

3. **DNS Records** (in Hostinger) should be:
   
   **A Records (4 total):**
   - `@` → `185.199.108.153`
   - `@` → `185.199.109.153`
   - `@` → `185.199.110.153`
   - `@` → `185.199.111.153`
   
   **CNAME Record:**
   - `www` → `m3owle0.github.io` (your GitHub username, NOT the full path with /mmcs)

## Why This Matters

- GitHub Pages custom domains work at the **repository level**, not the subdirectory level
- When you add `multimarketclothingsearch.com` as a custom domain in GitHub Pages settings, GitHub will route it to your site
- The DNS CNAME for `www` should point to `m3owle0.github.io` (without `/mmcs`)
- GitHub automatically handles routing to the correct subdirectory

## Verify Setup

1. **Check CNAME file location:**
   - Go to: `https://github.com/m3owle0/mmcs/blob/main/mmcs/CNAME`
   - Or: `https://github.com/m3owle0/mmcs/blob/main/CNAME` (if mmcs is the repo name)
   - Should contain: `multimarketclothingsearch.com`

2. **Check GitHub Pages settings:**
   - Repo → Settings → Pages
   - Custom domain should show: `multimarketclothingsearch.com`
   - Status should be ✅ (green checkmark)

3. **Test the domain:**
   - Visit: `https://multimarketclothingsearch.com`
   - Should redirect/load your site from `/mmcs`

## If It's Still Not Working

### Option 1: Move files to root (if possible)
If you can restructure your repo:
- Move `index.html` and assets to the root of the repo
- Update GitHub Pages to deploy from root
- Then `CNAME` goes in root

### Option 2: Use a redirect (if GitHub Pages doesn't support subdirectory custom domains)
Some GitHub Pages configurations don't support custom domains for subdirectories. In that case:
- Keep DNS pointing to GitHub
- Use a redirect service or configure at the hosting level

## Quick Fix Checklist

- [ ] `CNAME` file is in the `mmcs` folder (or wherever your site files are)
- [ ] `CNAME` contains: `multimarketclothingsearch.com`
- [ ] DNS A records point to GitHub IPs (185.199.108.153, etc.)
- [ ] DNS CNAME for `www` points to `m3owle0.github.io`
- [ ] Custom domain added in GitHub Pages settings
- [ ] Waited 5-30 minutes for DNS propagation
