# Hostinger DNS Configuration for GitHub Pages

## Option 1: Use Hostinger DNS (Keep Current Nameservers) ✅ RECOMMENDED

**Keep your current Hostinger nameservers** (`ns1.dns-parking.com` and `ns2.dns-parking.com` or whatever Hostinger assigned).

Instead, configure **DNS Records** in Hostinger:

### Steps:

1. **In Hostinger, go to DNS Management** (not Nameservers section)
   - Look for "DNS Zone" or "DNS Management" or "Advanced DNS"
   - This is different from the Nameservers page

2. **Add these A Records** (for the root domain):

   | Type | Name | Value | TTL |
   |------|------|-------|-----|
   | A | @ | 185.199.108.153 | 3600 |
   | A | @ | 185.199.109.153 | 3600 |
   | A | @ | 185.199.110.153 | 3600 |
   | A | @ | 185.199.111.153 | 3600 |

   **Note:** Use `@` or leave Name blank for root domain

3. **Add CNAME Record** (for www subdomain - optional):

   | Type | Name | Value | TTL |
   |------|------|-------|-----|
   | CNAME | www | your-username.github.io | 3600 |

   **Replace `your-username` with your actual GitHub username**

4. **Save all records**

5. **Wait 5-30 minutes** for DNS to propagate

---

## Option 2: Use Cloudflare Nameservers (If Hostinger DNS is Limited)

If Hostinger doesn't allow easy DNS record management, use Cloudflare:

### Steps:

1. **Sign up for Cloudflare** (free): https://cloudflare.com

2. **Add your domain** to Cloudflare

3. **Cloudflare will give you nameservers** like:
   - `alice.ns.cloudflare.com`
   - `bob.ns.cloudflare.com`

4. **In Hostinger, change nameservers:**
   - Go to Nameservers section (where you saw the current ones)
   - Click "Change Nameservers"
   - Enter Cloudflare's nameservers
   - Save

5. **In Cloudflare DNS panel, add these records:**

   | Type | Name | Content | TTL |
   |------|------|---------|-----|
   | A | @ | 185.199.108.153 | Auto |
   | A | @ | 185.199.109.153 | Auto |
   | A | @ | 185.199.110.153 | Auto |
   | A | @ | 185.199.111.153 | Auto |
   | CNAME | www | your-username.github.io | Auto |

6. **Wait for nameserver changes** (can take 24-48 hours)

---

## Important Notes

- **GitHub Pages does NOT provide nameservers** - you must use DNS records
- **Option 1 is easier** if Hostinger allows DNS record management
- **DNS changes can take 24-48 hours** to fully propagate worldwide
- **Check DNS propagation:** https://dnschecker.org/#A/multimarketclothingsearch.com

## Verify Setup

1. **Check DNS records:**
   ```bash
   # In terminal/command prompt
   nslookup multimarketclothingsearch.com
   ```

2. **Check GitHub Pages:**
   - Go to your repo → Settings → Pages
   - Under "Custom domain" you should see: `multimarketclothingsearch.com`
   - Status should show ✅ (green checkmark) when DNS is correct

3. **Test the site:**
   - Visit: https://multimarketclothingsearch.com
   - Should load your GitHub Pages site

## Troubleshooting

### "Domain not verified" in GitHub Pages
- DNS records may not have propagated yet
- Wait 24-48 hours
- Verify records are correct at dnschecker.org

### Site shows "Not found"
- Make sure `CNAME` file is in your GitHub repo root
- Verify DNS records point to GitHub IPs
- Check GitHub Pages is enabled in repo settings

### Can't find DNS Management in Hostinger
- Look for "Advanced DNS" or "DNS Zone Editor"
- Contact Hostinger support if you can't find it
- Consider Option 2 (Cloudflare) if Hostinger doesn't support custom DNS records

---

**Quick Answer:** Keep Hostinger nameservers, but add A records pointing to GitHub's IP addresses (185.199.108.153, 185.199.109.153, 185.199.110.153, 185.199.111.153) in Hostinger's DNS management panel.
