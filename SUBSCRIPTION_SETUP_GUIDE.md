# Smart Subscription System - Complete Setup Guide

This guide will walk you through setting up the smart subscription system step by step.

## 📋 Table of Contents
1. [PayPal Subscription Setup](#1-paypal-subscription-setup)
2. [Stripe Checkout Links Setup](#2-stripe-checkout-links-setup)
3. [Database Columns Setup (Optional)](#3-database-columns-setup-optional)
4. [Update Code with Payment Links](#4-update-code-with-payment-links)
5. [Testing the System](#5-testing-the-system)
6. [Managing Subscriptions](#6-managing-subscriptions)

---

## 1. PayPal Subscription Setup

### Step 1.1: Create PayPal Business Account
- Go to [PayPal Business](https://www.paypal.com/business)
- Sign up or log in to your business account
- Complete business verification if needed

### Step 1.2: Create Subscription Plans

1. **Go to PayPal Dashboard**
   - Log in to [PayPal Business Dashboard](https://www.paypal.com/businessmanage/account/home)
   - Navigate to **Products** → **Subscriptions**

2. **Create Basic Plan ($5/month)**
   - Click **"Create Plan"** or **"New Product"**
   - **Plan Name**: `MMCS Basic Plan`
   - **Description**: `Unlimited searches, premium badge, priority support`
   - **Billing Cycle**: 
     - Frequency: `Monthly`
     - Price: `$5.00 USD`
   - **Plan ID**: Copy this ID (you'll need it later)
   - Click **"Save"**

3. **Create Pro Plan ($10/month)**
   - Click **"Create Plan"** again
   - **Plan Name**: `MMCS Pro Plan`
   - **Description**: `Everything in Basic + Discord notifications, advanced filters`
   - **Billing Cycle**: 
     - Frequency: `Monthly`
     - Price: `$10.00 USD`
   - **Plan ID**: Copy this ID
   - Click **"Save"**

4. **Create Premium Plan ($20/month)**
   - Click **"Create Plan"** again
   - **Plan Name**: `MMCS Premium Plan`
   - **Description**: `Everything in Pro + unlimited notifications, API access`
   - **Billing Cycle**: 
     - Frequency: `Monthly`
     - Price: `$20.00 USD`
   - **Plan ID**: Copy this ID
   - Click **"Save"**

### Step 1.3: Get Subscription Links

For each plan, you need to create a subscription link:

1. **Option A: Use PayPal Subscription Buttons**
   - Go to **Tools** → **PayPal Buttons**
   - Select **"Subscription"** button type
   - Choose your plan
   - Copy the generated link

2. **Option B: Use PayPal API**
   - The link format is: `https://www.paypal.com/webapps/billing/subscriptions?plan_id=YOUR_PLAN_ID`
   - Replace `YOUR_PLAN_ID` with the Plan ID you copied

**Save these 3 links** - you'll need them in Step 4.

---

## 2. Stripe Checkout Links Setup

### Step 2.1: Create Stripe Account
- Go to [Stripe](https://stripe.com)
- Sign up or log in
- Complete account verification

### Step 2.2: Create Products and Prices

1. **Go to Stripe Dashboard**
   - Log in to [Stripe Dashboard](https://dashboard.stripe.com)
   - Navigate to **Products**

2. **Create Basic Product**
   - Click **"+ Add product"**
   - **Name**: `MMCS Basic Plan`
   - **Description**: `Unlimited searches, premium badge, priority support`
   - **Pricing**: 
     - Price: `$5.00`
     - Billing period: `Monthly`
   - Click **"Save product"**

3. **Create Pro Product**
   - Click **"+ Add product"** again
   - **Name**: `MMCS Pro Plan`
   - **Description**: `Everything in Basic + Discord notifications, advanced filters`
   - **Pricing**: 
     - Price: `$10.00`
     - Billing period: `Monthly`
   - Click **"Save product"**

4. **Create Premium Product**
   - Click **"+ Add product"** again
   - **Name**: `MMCS Premium Plan`
   - **Description**: `Everything in Pro + unlimited notifications, API access`
   - **Pricing**: 
     - Price: `$20.00`
     - Billing period: `Monthly`
   - Click **"Save product"**

### Step 2.3: Create Payment Links

1. **Go to Payment Links**
   - Navigate to **Products** → **Payment Links** (or use the left sidebar)
   - Click **"+ Create payment link"**

2. **Create Basic Payment Link**
   - Select your **Basic product**
   - Set as **Recurring** subscription
   - Click **"Create link"**
   - Copy the link (format: `https://buy.stripe.com/...`)

3. **Create Pro Payment Link**
   - Repeat for Pro product
   - Copy the link

4. **Create Premium Payment Link**
   - Repeat for Premium product
   - Copy the link

**Save these 3 links** - you'll need them in Step 4.

---

## 3. Database Columns Setup (Optional)

This step is optional but recommended for better subscription tracking.

### Step 3.1: Access Supabase Dashboard
- Go to [Supabase Dashboard](https://supabase.com/dashboard)
- Select your project
- Go to **SQL Editor**

### Step 3.2: Add Subscription Columns

Run this SQL in the SQL Editor:

```sql
-- Add subscription tier column
ALTER TABLE unlocked_users 
ADD COLUMN IF NOT EXISTS subscription_tier TEXT DEFAULT 'trial';

-- Add subscription expiry date
ALTER TABLE unlocked_users 
ADD COLUMN IF NOT EXISTS subscription_expires_at TIMESTAMPTZ;

-- Add payment method column
ALTER TABLE unlocked_users 
ADD COLUMN IF NOT EXISTS payment_method TEXT;

-- Add index for faster queries
CREATE INDEX IF NOT EXISTS idx_unlocked_users_subscription_tier 
ON unlocked_users(subscription_tier);

CREATE INDEX IF NOT EXISTS idx_unlocked_users_subscription_expires 
ON unlocked_users(subscription_expires_at);

-- Add comments for documentation
COMMENT ON COLUMN unlocked_users.subscription_tier IS 'Subscription tier: trial, basic, pro, premium';
COMMENT ON COLUMN unlocked_users.subscription_expires_at IS 'When the subscription expires (for recurring subscriptions)';
COMMENT ON COLUMN unlocked_users.payment_method IS 'Payment method used: paypal, stripe, cashapp, bitcoin, one-time';
```

### Step 3.3: Verify Columns Added

- Go to **Table Editor** → `unlocked_users`
- Verify you see the new columns:
  - `subscription_tier`
  - `subscription_expires_at`
  - `payment_method`

---

## 4. Update Code with Payment Links

### Step 4.1: Open index.html

Open `c:\Users\puppiesandkittens\Downloads\mmcs\index.html` in your editor.

### Step 4.2: Find PayPal Subscription Links

Search for (around line 10850):
```javascript
const PAYPAL_SUBSCRIPTION_LINKS = {
  basic: 'https://www.paypal.com/webapps/billing/subscriptions?plan_id=YOUR_BASIC_PLAN_ID',
  pro: 'https://www.paypal.com/webapps/billing/subscriptions?plan_id=YOUR_PRO_PLAN_ID',
  premium: 'https://www.paypal.com/webapps/billing/subscriptions?plan_id=YOUR_PREMIUM_PLAN_ID'
};
```

### Step 4.3: Replace PayPal Links

Replace `YOUR_BASIC_PLAN_ID`, `YOUR_PRO_PLAN_ID`, and `YOUR_PREMIUM_PLAN_ID` with the actual Plan IDs you copied from PayPal in Step 1.2.

**Example:**
```javascript
const PAYPAL_SUBSCRIPTION_LINKS = {
  basic: 'https://www.paypal.com/webapps/billing/subscriptions?plan_id=P-5ML4271244454362WXNWU5NQ',
  pro: 'https://www.paypal.com/webapps/billing/subscriptions?plan_id=P-5ML4271244454362WXNWU5NQ',
  premium: 'https://www.paypal.com/webapps/billing/subscriptions?plan_id=P-5ML4271244454362WXNWU5NQ'
};
```

### Step 4.4: Find Stripe Checkout Links

Search for (around line 10855):
```javascript
const STRIPE_CHECKOUT_LINKS = {
  basic: 'https://buy.stripe.com/YOUR_BASIC_LINK',
  pro: 'https://buy.stripe.com/YOUR_PRO_LINK',
  premium: 'https://buy.stripe.com/YOUR_PREMIUM_LINK',
  discord: 'https://buy.stripe.com/00w6oH2q5g6YdryekFcwg00' // Existing Discord notifications
};
```

### Step 4.5: Replace Stripe Links

Replace `YOUR_BASIC_LINK`, `YOUR_PRO_LINK`, and `YOUR_PREMIUM_LINK` with the actual Stripe payment links you copied in Step 2.3.

**Example:**
```javascript
const STRIPE_CHECKOUT_LINKS = {
  basic: 'https://buy.stripe.com/test_basic123',
  pro: 'https://buy.stripe.com/test_pro456',
  premium: 'https://buy.stripe.com/test_premium789',
  discord: 'https://buy.stripe.com/00w6oH2q5g6YdryekFcwg00' // Existing Discord notifications
};
```

### Step 4.6: Save the File

Save `index.html` after making these changes.

---

## 5. Testing the System

### Step 5.1: Test Upgrade Modal

1. Open your website
2. Click the **"Unlock Site"** button (or any upgrade prompt)
3. Verify the upgrade modal shows all 4 tiers:
   - Free Trial
   - Basic (with ⭐ Recommended badge)
   - Pro
   - Premium

### Step 5.2: Test Payment Buttons

1. Click on a subscription tier (e.g., Basic)
2. Click **"Subscribe via PayPal"**
   - Should open PayPal in a new tab
   - Verify it shows the correct plan and price
3. Click **"Subscribe via Stripe"**
   - Should open Stripe checkout in a new tab
   - Verify it shows the correct plan and price
4. Click **"One-Time Payment"**
   - Should show payment options

### Step 5.3: Test Subscription Management

1. Log in to your account
2. Click your profile picture/avatar
3. Click **"💳 Manage Subscription"**
4. Verify the subscription management modal shows:
   - Current plan
   - Subscription status badge
   - Next billing date (if applicable)
   - Payment method
   - Option to change plan
   - Cancel subscription button

### Step 5.4: Test Subscription Status Display

1. Check the header for subscription badges
2. Verify premium badge shows for verified users
3. Verify trial searches display shows for trial users

---

## 6. Managing Subscriptions

### Step 6.1: Manual Verification (Current Method)

When a user makes a payment:

1. **PayPal Subscription:**
   - Check PayPal dashboard for new subscriptions
   - Note the user's email from the subscription
   - Update their `verified` status in Supabase

2. **Stripe Subscription:**
   - Check Stripe dashboard for new subscriptions
   - Note the user's email
   - Update their `verified` status in Supabase

3. **One-Time Payment:**
   - User messages you on Instagram with their email
   - Update their `verified` status in Supabase

### Step 6.2: Update User Subscription in Supabase

1. Go to **Supabase Dashboard** → **Table Editor** → `unlocked_users`
2. Find the user by email
3. Update these fields:
   - `verified`: Set to `true`
   - `subscription_tier`: Set to `basic`, `pro`, or `premium`
   - `subscription_expires_at`: Set to 30 days from now (for monthly subscriptions)
   - `payment_method`: Set to `paypal`, `stripe`, `cashapp`, or `bitcoin`

**Example SQL:**
```sql
UPDATE unlocked_users
SET 
  verified = true,
  subscription_tier = 'basic',
  subscription_expires_at = NOW() + INTERVAL '30 days',
  payment_method = 'paypal'
WHERE email = 'user@example.com';
```

### Step 6.3: Set Up Webhooks (Advanced - Optional)

For automatic subscription activation:

1. **PayPal Webhooks:**
   - Go to PayPal Dashboard → **Webhooks**
   - Create webhook endpoint
   - Subscribe to `BILLING.SUBSCRIPTION.CREATED` and `BILLING.SUBSCRIPTION.UPDATED`
   - Point to your backend endpoint

2. **Stripe Webhooks:**
   - Go to Stripe Dashboard → **Developers** → **Webhooks**
   - Add endpoint
   - Subscribe to `checkout.session.completed` and `customer.subscription.created`
   - Point to your backend endpoint

3. **Backend Handler:**
   - Create an endpoint that receives webhook events
   - Update Supabase when subscription is created/updated
   - This requires backend code (not included in this guide)

---

## 7. Troubleshooting

### Issue: Payment links don't work

**Solution:**
- Verify you copied the correct links
- Check that links are not in test mode (if using Stripe)
- Ensure PayPal/Stripe accounts are verified

### Issue: Subscription status not updating

**Solution:**
- Check browser console for errors
- Verify database columns exist
- Check Supabase connection

### Issue: Upgrade modal not showing

**Solution:**
- Clear browser cache
- Check JavaScript console for errors
- Verify `index.html` was saved correctly

### Issue: Users can't see their subscription status

**Solution:**
- Verify user is logged in
- Check `subscription_tier` column in database
- Ensure `updateSubscriptionDisplay()` function is called

---

## 8. Next Steps

1. **Monitor Subscriptions:**
   - Regularly check PayPal/Stripe dashboards
   - Track subscription renewals
   - Handle cancellations

2. **Automate (Optional):**
   - Set up webhooks for automatic activation
   - Create a backend service to handle payments
   - Implement subscription expiry checks

3. **Marketing:**
   - Promote subscription tiers
   - Highlight premium features
   - Offer limited-time discounts

---

## 📝 Quick Reference

### PayPal Plan IDs Location
- PayPal Dashboard → Products → Subscriptions → [Your Plan] → Plan ID

### Stripe Payment Links Location
- Stripe Dashboard → Products → Payment Links → [Your Link]

### Database Columns
- `subscription_tier`: `trial`, `basic`, `pro`, `premium`
- `subscription_expires_at`: Timestamp (for recurring subscriptions)
- `payment_method`: `paypal`, `stripe`, `cashapp`, `bitcoin`, `one-time`

### Code Locations
- PayPal links: Line ~10850 in `index.html`
- Stripe links: Line ~10855 in `index.html`
- Subscription functions: Lines ~10870-11000 in `index.html`

---

## ✅ Setup Checklist

- [ ] PayPal business account created
- [ ] 3 PayPal subscription plans created (Basic, Pro, Premium)
- [ ] PayPal Plan IDs copied
- [ ] Stripe account created
- [ ] 3 Stripe products created (Basic, Pro, Premium)
- [ ] 3 Stripe payment links created
- [ ] Database columns added (optional)
- [ ] PayPal links updated in code
- [ ] Stripe links updated in code
- [ ] Code saved
- [ ] Tested upgrade modal
- [ ] Tested payment buttons
- [ ] Tested subscription management
- [ ] Verified subscription status display

---

**Need Help?** Check the troubleshooting section or review the code comments in `index.html`.
