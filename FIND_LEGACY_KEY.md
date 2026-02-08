# 🔑 Finding the Service Role Key (Legacy Tab)

You're on the API Keys page, but you need the **Legacy** keys!

## ✅ What to Do:

1. **Look at the tabs** at the top of the API Keys section
2. **Click on:** "Legacy anon, service_role API keys" tab
3. **You'll see:**
   - `anon` key - ❌ NOT this one
   - `service_role` key - ✅ **THIS ONE!**
4. **Click the eye icon** or "Reveal" to show the full key
5. **Copy the entire key**

## 📍 Where You Are:
- Current tab: "Publishable and secret API keys" 
- **Switch to:** "Legacy anon, service_role API keys" tab

The `service_role` key is in the Legacy section, not the new API keys section!

---

**Once you have it:** Add it to GitHub Secrets as `SUPABASE_SERVICE_ROLE_KEY`
