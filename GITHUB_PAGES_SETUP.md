# GitHub Pages Setup for multimarketclothingsearch.com

## Step 1: Verify Required Files Are in Your Repository

Make sure these files are in the **root** of your GitHub repository:

✅ **Required Files:**
- `index.html` (main website file)
- `miku.svg` (favicon)
- `mmcs.svg` (logo)
- `bg.svg` (background image)
- `CNAME` (custom domain file - I just created this for you)

❌ **Do NOT upload these folders** (they're not needed for the website):
- `notifier/` folder
- `database/` folder

## Step 2: Upload Files to GitHub

1. **Go to your GitHub repository**
2. **Click "Add file" → "Upload files"**
3. **Upload these files:**
   - `index.html`
   - `miku.svg`
   - `mmcs.svg`
   - `bg.svg`
   - `CNAME` (the file I just created)
4. **Click "Commit changes"**

## Step 3: Configure GitHub Pages

1. **Go to your repository Settings** (top right of repo page)
2. **Click "Pages" in the left sidebar**
3. **Under "Source":**
   - Select: **Deploy from a branch**
   - Branch: **main** (or **master**)
   - Folder: **/ (root)**
4. **Click "Save"**
5. **Wait 1-2 minutes** for GitHub to build your site

## Step 4: Configure Custom Domain

1. **Still in Settings → Pages**
2. **Under "Custom domain":**
   - Enter: `multimarketclothingsearch.com`
   - Click "Save"
3. **GitHub will create/update the CNAME file** (or use the one I created)

## Step 5: Configure DNS at Your Domain Registrar

You need to add DNS records at wherever you bought your domain (GoDaddy, Namecheap, Cloudflare, etc.):

### Option A: CNAME Record (Recommended)
- **Type:** CNAME
- **Name:** @ (or leave blank, or `www`)
- **Value:** `your-username.github.io` (replace with your GitHub username)
- **TTL:** 3600 (or default)

### Option B: A Records (Alternative)
If CNAME doesn't work, use these A records:
- **Type:** A
- **Name:** @
- **Value:** `185.199.108.153`
- **TTL:** 3600

- **Type:** A
- **Name:** @
- **Value:** `185.199.109.153`

- **Type:** A
- **Name:** @
- **Value:** `185.199.110.153`

- **Type:** A
- **Name:** @
- **Value:** `185.199.111.153`

### For www subdomain (optional):
- **Type:** CNAME
- **Name:** www
- **Value:** `your-username.github.io`

## Step 6: Verify It's Working

1. **Wait 5-10 minutes** after DNS changes
2. **Check DNS propagation:** https://whatsmydns.net/#CNAME/multimarketclothingsearch.com
3. **Visit:** https://multimarketclothingsearch.com
4. **Check GitHub Pages status:**
   - Go to Settings → Pages
   - You should see a green checkmark ✅
   - If there's an error, click it to see details

## Troubleshooting

### Site shows 404 or "Page not found"
- Make sure `index.html` is in the root directory (not in a subfolder)
- Check GitHub Pages is enabled (Settings → Pages)
- Wait a few minutes for GitHub to rebuild

### Domain shows "Not found" or DNS error
- Verify DNS records are correct at your domain registrar
- DNS changes can take 24-48 hours to propagate
- Check: https://dnschecker.org/#CNAME/multimarketclothingsearch.com

### Site loads but images are broken
- Make sure `miku.svg`, `mmcs.svg`, and `bg.svg` are in the same directory as `index.html`
- Check file names match exactly (case-sensitive)

### GitHub Pages shows "Custom domain not verified"
- Make sure DNS records are set correctly
- Wait for DNS propagation (can take up to 48 hours)
- GitHub will automatically verify once DNS is correct

## Quick Checklist

- [ ] All files uploaded to GitHub repo root
- [ ] `CNAME` file created with `multimarketclothingsearch.com`
- [ ] GitHub Pages enabled (Settings → Pages)
- [ ] Custom domain added in GitHub Pages settings
- [ ] DNS records configured at domain registrar
- [ ] Waited for DNS propagation (check with dnschecker.org)
- [ ] Site accessible at https://multimarketclothingsearch.com

## Need Help?

- **GitHub Pages Docs:** https://docs.github.com/en/pages
- **Custom Domain Help:** https://docs.github.com/en/pages/configuring-a-custom-domain-for-your-github-pages-site

---

**Note:** The `CNAME` file I created should be committed to your repository. It tells GitHub Pages to use your custom domain.
