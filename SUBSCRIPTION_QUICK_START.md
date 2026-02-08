# Subscription System - Quick Start Guide

## 🚀 5-Minute Setup (Minimum Required)

### Step 1: Get PayPal Subscription Links (2 minutes)

1. Go to [PayPal Business Dashboard](https://www.paypal.com/businessmanage/account/home)
2. Navigate to **Products** → **Subscriptions**
3. Create 3 plans:
   - **Basic**: $5/month
   - **Pro**: $10/month  
   - **Premium**: $20/month
4. Copy the Plan IDs (looks like: `P-5ML4271244454362WXNWU5NQ`)

### Step 2: Get Stripe Payment Links (2 minutes)

1. Go to [Stripe Dashboard](https://dashboard.stripe.com)
2. Navigate to **Products** → **Payment Links**
3. Create 3 payment links:
   - **Basic**: $5/month recurring
   - **Pro**: $10/month recurring
   - **Premium**: $20/month recurring
4. Copy the links (looks like: `https://buy.stripe.com/...`)

### Step 3: Update Code (1 minute)

1. Open `index.html`
2. Find line ~10850, replace PayPal links:
   ```javascript
   const PAYPAL_SUBSCRIPTION_LINKS = {
     basic: 'https://www.paypal.com/webapps/billing/subscriptions?plan_id=YOUR_BASIC_ID',
     pro: 'https://www.paypal.com/webapps/billing/subscriptions?plan_id=YOUR_PRO_ID',
     premium: 'https://www.paypal.com/webapps/billing/subscriptions?plan_id=YOUR_PREMIUM_ID'
   };
   ```

3. Find line ~10855, replace Stripe links:
   ```javascript
   const STRIPE_CHECKOUT_LINKS = {
     basic: 'https://buy.stripe.com/YOUR_BASIC_LINK',
     pro: 'https://buy.stripe.com/YOUR_PRO_LINK',
     premium: 'https://buy.stripe.com/YOUR_PREMIUM_LINK',
     discord: 'https://buy.stripe.com/00w6oH2q5g6YdryekFcwg00'
   };
   ```

4. Save the file

**Done!** The subscription system is now active.

---

## 📍 Where to Find Things

### In Your Code (`index.html`)

- **PayPal Links**: Line ~10850
- **Stripe Links**: Line ~10855
- **Subscription Functions**: Lines ~10870-11000
- **Upgrade Modal HTML**: Line ~3970
- **Subscription Management Modal**: Line ~4025

### In PayPal

- **Dashboard**: https://www.paypal.com/businessmanage/account/home
- **Subscriptions**: Products → Subscriptions
- **Plan ID Format**: `P-XXXXXXXXXXXXX`

### In Stripe

- **Dashboard**: https://dashboard.stripe.com
- **Payment Links**: Products → Payment Links
- **Link Format**: `https://buy.stripe.com/...`

### In Supabase

- **Dashboard**: https://supabase.com/dashboard
- **Table**: `unlocked_users`
- **SQL Editor**: Left sidebar → SQL Editor

---

## 🔧 Manual Subscription Activation

When a user pays, manually activate their subscription:

1. **Find User in Supabase**
   - Go to Table Editor → `unlocked_users`
   - Search by email

2. **Update Subscription**
   ```sql
   UPDATE unlocked_users
   SET 
     verified = true,
     subscription_tier = 'basic',  -- or 'pro' or 'premium'
     subscription_expires_at = NOW() + INTERVAL '30 days',
     payment_method = 'paypal'  -- or 'stripe', 'cashapp', 'bitcoin'
   WHERE email = 'user@example.com';
   ```

---

## 🎯 Subscription Tiers

| Tier | Price | Features |
|------|-------|----------|
| **Free Trial** | $0 | 75 free searches |
| **Basic** | $5/mo | Unlimited searches, premium badge, priority support |
| **Pro** | $10/mo | Basic + Discord notifications, advanced filters |
| **Premium** | $20/mo | Pro + unlimited notifications, API access |

---

## ⚠️ Important Notes

1. **Test Mode**: Stripe links work in test mode - switch to live mode for real payments
2. **PayPal Verification**: PayPal accounts need business verification for subscriptions
3. **Manual Activation**: Currently requires manual activation - webhooks can automate this later
4. **Database Columns**: Optional but recommended for better tracking

---

## 🆘 Quick Troubleshooting

**Payment buttons don't work?**
- Check links are correct (no `YOUR_` placeholders)
- Verify PayPal/Stripe accounts are active

**Subscription status not showing?**
- Check user is logged in
- Verify database has `subscription_tier` column (optional)

**Upgrade modal not appearing?**
- Clear browser cache
- Check JavaScript console for errors

---

## 📚 Full Documentation

See `SUBSCRIPTION_SETUP_GUIDE.md` for complete step-by-step instructions.
