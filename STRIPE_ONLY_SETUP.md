# Stripe-Only Subscription Setup Guide

This guide will walk you through setting up the subscription system using **only Stripe** for payments.

## 📋 Quick Overview

You'll create 3 Stripe payment links (one for each tier) and add them to your code. That's it!

---

## Step 1: Create Stripe Account (2 minutes)

1. Go to [Stripe.com](https://stripe.com)
2. Click **"Sign up"** or **"Log in"**
3. Complete account setup and verification
4. **Important**: Switch from **Test mode** to **Live mode** when ready for real payments

---

## Step 2: Create Products in Stripe (5 minutes)

### Step 2.1: Go to Products

1. Log in to [Stripe Dashboard](https://dashboard.stripe.com)
2. Click **"Products"** in the left sidebar
3. Click **"+ Add product"** button

### Step 2.2: Create Basic Plan ($5/month)

1. **Product name**: `MMCS Basic Plan`
2. **Description**: `Unlimited searches, premium badge, priority support`
3. **Pricing**:
   - **Price**: `5.00`
   - **Currency**: `USD`
   - **Billing period**: Select **"Recurring"** → **"Monthly"**
4. Click **"Save product"**
5. **Copy the product ID** (you'll see it after saving, looks like: `prod_xxxxxxxxxxxxx`)

### Step 2.3: Create Pro Plan ($10/month)

1. Click **"+ Add product"** again
2. **Product name**: `MMCS Pro Plan`
3. **Description**: `Everything in Basic + Discord notifications, advanced filters`
4. **Pricing**:
   - **Price**: `10.00`
   - **Currency**: `USD`
   - **Billing period**: **"Recurring"** → **"Monthly"**
5. Click **"Save product"**
6. **Copy the product ID**

### Step 2.4: Create Premium Plan ($20/month)

1. Click **"+ Add product"** again
2. **Product name**: `MMCS Premium Plan`
3. **Description**: `Everything in Pro + unlimited notifications, API access`
4. **Pricing**:
   - **Price**: `20.00`
   - **Currency**: `USD`
   - **Billing period**: **"Recurring"** → **"Monthly"**
5. Click **"Save product"**
6. **Copy the product ID**

---

## Step 3: Create Payment Links (5 minutes)

### Step 3.1: Go to Payment Links

1. In Stripe Dashboard, click **"Payment Links"** in the left sidebar
   - (If you don't see it, go to **Products** → Click on a product → **"Create payment link"**)

### Step 3.2: Create Basic Payment Link

1. Click **"+ Create payment link"**
2. Select your **"MMCS Basic Plan"** product
3. Make sure it's set to **"Recurring"** subscription
4. Click **"Create link"**
5. **Copy the link** - it will look like: `https://buy.stripe.com/xxxxxxxxxxxxx`
6. **Save this link** - you'll need it in Step 4

### Step 3.3: Create Pro Payment Link

1. Click **"+ Create payment link"** again
2. Select your **"MMCS Pro Plan"** product
3. Make sure it's set to **"Recurring"** subscription
4. Click **"Create link"**
5. **Copy the link**
6. **Save this link**

### Step 3.4: Create Premium Payment Link

1. Click **"+ Create payment link"** again
2. Select your **"MMCS Premium Plan"** product
3. Make sure it's set to **"Recurring"** subscription
4. Click **"Create link"**
5. **Copy the link**
6. **Save this link**

---

## Step 4: Update Your Code (2 minutes)

### Step 4.1: Open index.html

Open `c:\Users\puppiesandkittens\Downloads\mmcs\index.html` in your code editor.

### Step 4.2: Find Stripe Links Section

Press `Ctrl+F` (or `Cmd+F` on Mac) and search for:
```
STRIPE_CHECKOUT_LINKS
```

You should find it around **line 11010**. It will look like this:

```javascript
const STRIPE_CHECKOUT_LINKS = {
  basic: 'https://buy.stripe.com/YOUR_BASIC_LINK',
  pro: 'https://buy.stripe.com/YOUR_PRO_LINK',
  premium: 'https://buy.stripe.com/YOUR_PREMIUM_LINK',
  discord: 'https://buy.stripe.com/00w6oH2q5g6YdryekFcwg00'
};
```

### Step 4.3: Replace with Your Links

Replace the placeholder links with your actual Stripe payment links:

**Example:**
```javascript
const STRIPE_CHECKOUT_LINKS = {
  basic: 'https://buy.stripe.com/a1b2c3d4e5f6g7h8i9j0',
  pro: 'https://buy.stripe.com/x9y8z7w6v5u4t3s2r1q0',
  premium: 'https://buy.stripe.com/p0o9i8u7y6t5r4e3w2q1',
  discord: 'https://buy.stripe.com/00w6oH2q5g6YdryekFcwg00'
};
```

**Important**: 
- Replace `YOUR_BASIC_LINK` with your Basic plan link
- Replace `YOUR_PRO_LINK` with your Pro plan link  
- Replace `YOUR_PREMIUM_LINK` with your Premium plan link
- Keep the `discord` link as-is (it's already configured)

### Step 4.4: Save the File

Save `index.html` after making the changes.

---

## Step 5: Test the System (3 minutes)

### Step 5.1: Test Upgrade Modal

1. Open your website in a browser
2. Click the **"Unlock Site"** button or any upgrade prompt
3. You should see the upgrade modal with 4 tiers:
   - Free Trial
   - Basic ($5/month) - with ⭐ Recommended badge
   - Pro ($10/month)
   - Premium ($20/month)

### Step 5.2: Test Payment Buttons

1. Click on the **Basic** tier
2. Click the **"💳 Subscribe with Stripe"** button
3. It should open Stripe checkout in a new tab
4. **In Test Mode**: You can use test card `4242 4242 4242 4242` with any future expiry date
5. Verify it shows the correct plan and price ($5/month)

Repeat for Pro and Premium tiers to verify all links work.

### Step 5.3: Test Subscription Management

1. Log in to your account
2. Click your profile picture/avatar
3. Click **"💳 Manage Subscription"**
4. Verify the subscription management modal appears

---

## Step 6: Activate Subscriptions (Manual Process)

When a user subscribes through Stripe:

### Step 6.1: Check Stripe Dashboard

1. Go to [Stripe Dashboard](https://dashboard.stripe.com)
2. Click **"Customers"** in the left sidebar
3. Find the customer who just subscribed
4. Note their **email address**

### Step 6.2: Update User in Supabase

1. Go to [Supabase Dashboard](https://supabase.com/dashboard)
2. Select your project
3. Go to **Table Editor** → `unlocked_users`
4. Find the user by email
5. Update these fields:
   - `verified`: Set to `true`
   - `subscription_tier`: Set to `basic`, `pro`, or `premium` (based on which plan they chose)
   - `subscription_expires_at`: Set to 30 days from now
   - `payment_method`: Set to `stripe`

**Quick SQL Method:**
```sql
UPDATE unlocked_users
SET 
  verified = true,
  subscription_tier = 'basic',  -- or 'pro' or 'premium'
  subscription_expires_at = NOW() + INTERVAL '30 days',
  payment_method = 'stripe'
WHERE email = 'user@example.com';
```

---

## 🎯 Subscription Tiers Summary

| Tier | Price | Features |
|------|-------|----------|
| **Free Trial** | $0 | 75 free searches |
| **Basic** | $5/mo | Unlimited searches, premium badge, priority support |
| **Pro** | $10/mo | Basic + Discord notifications, advanced filters |
| **Premium** | $20/mo | Pro + unlimited notifications, API access |

---

## ⚠️ Important Notes

### Test Mode vs Live Mode

- **Test Mode**: Use for testing with test cards (no real charges)
- **Live Mode**: Switch when ready for real payments
- Toggle in Stripe Dashboard (top right corner)

### Test Cards (Test Mode Only)

- **Success**: `4242 4242 4242 4242`
- **Decline**: `4000 0000 0000 0002`
- Use any future expiry date and any CVC

### Subscription Renewals

- Stripe automatically handles monthly renewals
- You'll need to manually update `subscription_expires_at` each month (or set up webhooks for automation)

---

## 🔧 Troubleshooting

### Payment buttons don't work?

- **Check**: Links are correct (no `YOUR_` placeholders)
- **Check**: Stripe account is active
- **Check**: Browser console for errors (F12 → Console)

### Stripe checkout shows wrong price?

- **Check**: You selected the correct product when creating payment link
- **Check**: Product pricing is set correctly in Stripe

### Subscription status not updating?

- **Check**: User is logged in
- **Check**: Database columns exist (`subscription_tier`, `subscription_expires_at`)
- **Check**: You manually updated the user in Supabase

### Can't find Payment Links in Stripe?

- **Alternative**: Go to **Products** → Click a product → **"Create payment link"** button

---

## 📍 Quick Reference

### Where to Find Things

**In Stripe:**
- Dashboard: https://dashboard.stripe.com
- Products: Left sidebar → **Products**
- Payment Links: Left sidebar → **Payment Links** (or from Products page)
- Customers: Left sidebar → **Customers** (to see who subscribed)

**In Your Code:**
- Stripe links: Line ~11010 in `index.html`
- Search for: `STRIPE_CHECKOUT_LINKS`

**In Supabase:**
- Dashboard: https://supabase.com/dashboard
- Table: `unlocked_users`
- SQL Editor: Left sidebar → **SQL Editor**

---

## ✅ Setup Checklist

- [ ] Stripe account created and verified
- [ ] 3 products created (Basic $5, Pro $10, Premium $20)
- [ ] 3 payment links created (one for each product)
- [ ] Payment links copied
- [ ] Links updated in `index.html` (line ~11010)
- [ ] Code saved
- [ ] Tested upgrade modal
- [ ] Tested payment buttons (in test mode)
- [ ] Verified Stripe checkout opens correctly
- [ ] Know how to activate subscriptions manually

---

## 🚀 Next Steps (Optional)

### Automate Subscription Activation

Set up Stripe webhooks to automatically activate subscriptions:

1. **Create Webhook Endpoint**
   - Stripe Dashboard → **Developers** → **Webhooks**
   - Add endpoint URL (requires backend server)
   - Subscribe to: `checkout.session.completed` and `customer.subscription.created`

2. **Backend Handler**
   - Create endpoint that receives webhook events
   - Update Supabase when subscription is created
   - This requires backend code (not included in this guide)

### Monitor Subscriptions

- Check Stripe Dashboard regularly for new subscriptions
- Set up email notifications in Stripe for new customers
- Track subscription renewals and cancellations

---

**That's it!** Your Stripe-only subscription system is now set up and ready to use. 🎉
