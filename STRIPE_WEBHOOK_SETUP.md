# Automatic Subscription Activation via Stripe Webhooks

This guide will show you how to automatically upgrade users when they subscribe via Stripe, without manual database updates.

---

## 🎯 How It Works

1. **User subscribes** via Stripe checkout link
2. **Stripe sends webhook** to your server when payment succeeds
3. **Your server verifies** the webhook signature (security)
4. **Server updates Supabase** automatically with subscription details
5. **User is upgraded** instantly - no manual work needed!

---

## 📋 Prerequisites

- Stripe account with your payment links set up
- Supabase project with `unlocked_users` table
- A way to host a backend server (Vercel, Railway, Render, etc.)

---

## Option 1: Serverless Function (Recommended - Easiest)

### Step 1: Create a Vercel/Netlify Function

**Using Vercel (Free):**

1. **Install Vercel CLI:**
   ```bash
   npm i -g vercel
   ```

2. **Create project structure:**
   ```
   mmcs/
   ├── api/
   │   └── stripe-webhook.js
   ├── vercel.json
   └── package.json
   ```

3. **Create `api/stripe-webhook.js`:**
   ```javascript
   const { createClient } = require('@supabase/supabase-js');
   const stripe = require('stripe')(process.env.STRIPE_SECRET_KEY);

   // Map Stripe price IDs to subscription tiers
   const PRICE_TO_TIER = {
     'price_basic_id': 'basic',      // Replace with your Basic price ID
     'price_pro_id': 'pro',          // Replace with your Pro price ID
   };

   module.exports = async (req, res) => {
     // Only allow POST requests
     if (req.method !== 'POST') {
       return res.status(405).json({ error: 'Method not allowed' });
     }

     const sig = req.headers['stripe-signature'];
     const webhookSecret = process.env.STRIPE_WEBHOOK_SECRET;

     let event;

     try {
       // Verify webhook signature
       event = stripe.webhooks.constructEvent(req.body, sig, webhookSecret);
     } catch (err) {
       console.error('Webhook signature verification failed:', err.message);
       return res.status(400).send(`Webhook Error: ${err.message}`);
     }

     // Initialize Supabase with service role (bypasses RLS)
     const supabase = createClient(
       process.env.SUPABASE_URL,
       process.env.SUPABASE_SERVICE_ROLE_KEY // NOT anon key!
     );

     // Handle the event
     if (event.type === 'checkout.session.completed') {
       const session = event.data.object;
      
       // Get customer email from session
       const customerEmail = session.customer_details?.email || session.customer_email;
      
       if (!customerEmail) {
         console.error('No email found in session');
         return res.status(400).json({ error: 'No email found' });
       }

       // Determine subscription tier from price ID
       const priceId = session.line_items?.data[0]?.price?.id || session.amount_total;
       let subscriptionTier = 'basic'; // Default
      
       // Check which product they subscribed to
       if (session.amount_total === 500) { // $5.00 = 500 cents (Basic)
         subscriptionTier = 'basic';
       } else if (session.amount_total === 1000) { // $10.00 = 1000 cents (Pro)
         subscriptionTier = 'pro';
       }

       // Calculate expiry date (30 days from now)
       const expiresAt = new Date();
       expiresAt.setDate(expiresAt.getDate() + 30);

       // Update user in Supabase
       const { data, error } = await supabase
         .from('unlocked_users')
         .update({
           verified: true,
           subscription_tier: subscriptionTier,
           subscription_expires_at: expiresAt.toISOString(),
           payment_method: 'stripe',
           stripe_customer_id: session.customer || null,
           stripe_subscription_id: session.subscription || null
         })
         .eq('email', customerEmail);

       if (error) {
         console.error('Error updating user:', error);
         return res.status(500).json({ error: 'Failed to update user' });
       }

       console.log(`✅ Activated ${subscriptionTier} subscription for ${customerEmail}`);
       return res.json({ received: true, tier: subscriptionTier });
     }

     // Handle subscription updates/renewals
     if (event.type === 'customer.subscription.updated') {
       const subscription = event.data.object;
       const customerEmail = subscription.metadata?.email;

       if (customerEmail) {
         const expiresAt = new Date(subscription.current_period_end * 1000);
         
         await supabase
           .from('unlocked_users')
           .update({
             subscription_expires_at: expiresAt.toISOString(),
             verified: subscription.status === 'active'
           })
           .eq('email', customerEmail);
       }
     }

     // Handle subscription cancellations
     if (event.type === 'customer.subscription.deleted') {
       const subscription = event.data.object;
       const customerEmail = subscription.metadata?.email;

       if (customerEmail) {
         await supabase
           .from('unlocked_users')
           .update({
             verified: false,
             subscription_tier: 'trial'
           })
           .eq('email', customerEmail);
       }
     }

     res.json({ received: true });
   };
   ```

4. **Create `package.json`:**
   ```json
   {
     "name": "mmcs-webhook",
     "version": "1.0.0",
     "dependencies": {
       "@supabase/supabase-js": "^2.39.0",
       "stripe": "^14.0.0"
     }
   }
   ```

5. **Create `vercel.json`:**
   ```json
   {
     "functions": {
       "api/stripe-webhook.js": {
         "maxDuration": 10
       }
     }
   }
   ```

6. **Deploy to Vercel:**
   ```bash
   vercel
   ```

7. **Set Environment Variables in Vercel:**
   - `STRIPE_SECRET_KEY` - Your Stripe secret key
   - `STRIPE_WEBHOOK_SECRET` - Webhook signing secret (get from Stripe)
   - `SUPABASE_URL` - Your Supabase project URL
   - `SUPABASE_SERVICE_ROLE_KEY` - Supabase service role key (NOT anon key!)

---

## Option 2: Supabase Edge Function (No External Server Needed!)

This is the **easiest option** - uses Supabase's built-in serverless functions.

### Step 1: Install Supabase CLI

```bash
npm install -g supabase
```

### Step 2: Initialize Supabase Functions

```bash
cd mmcs
supabase init
supabase functions new stripe-webhook
```

### Step 3: Create the Webhook Function

Edit `supabase/functions/stripe-webhook/index.ts`:

```typescript
import { serve } from "https://deno.land/std@0.168.0/http/server.ts"
import { createClient } from 'https://esm.sh/@supabase/supabase-js@2'

const stripe = new (await import('https://esm.sh/stripe@14.0.0')).default(
  Deno.env.get('STRIPE_SECRET_KEY') || '',
  { apiVersion: '2023-10-16' }
)

serve(async (req) => {
  try {
    // Only allow POST
    if (req.method !== 'POST') {
      return new Response(JSON.stringify({ error: 'Method not allowed' }), {
        status: 405,
        headers: { 'Content-Type': 'application/json' }
      })
    }

    const signature = req.headers.get('stripe-signature')
    if (!signature) {
      return new Response(JSON.stringify({ error: 'No signature' }), {
        status: 400
      })
    }

    const body = await req.text()
    const webhookSecret = Deno.env.get('STRIPE_WEBHOOK_SECRET')

    // Verify webhook signature
    let event
    try {
      event = stripe.webhooks.constructEvent(body, signature, webhookSecret!)
    } catch (err) {
      return new Response(JSON.stringify({ error: `Webhook Error: ${err.message}` }), {
        status: 400
      })
    }

    // Initialize Supabase admin client
    const supabaseAdmin = createClient(
      Deno.env.get('SUPABASE_URL') ?? '',
      Deno.env.get('SUPABASE_SERVICE_ROLE_KEY') ?? ''
    )

    // Handle checkout.session.completed
    if (event.type === 'checkout.session.completed') {
      const session = event.data.object
      const customerEmail = session.customer_details?.email || session.customer_email

      if (!customerEmail) {
        return new Response(JSON.stringify({ error: 'No email found' }), {
          status: 400
        })
      }

      // Determine tier from amount (in cents)
      let subscriptionTier = 'basic'
      if (session.amount_total === 1000) { // $10.00 = Pro
        subscriptionTier = 'pro'
      } else if (session.amount_total === 500) { // $5.00 = Basic
        subscriptionTier = 'basic'
      }

      // Calculate expiry (30 days)
      const expiresAt = new Date()
      expiresAt.setDate(expiresAt.getDate() + 30)

      // Update user
      const { error } = await supabaseAdmin
        .from('unlocked_users')
        .update({
          verified: true,
          subscription_tier: subscriptionTier,
          subscription_expires_at: expiresAt.toISOString(),
          payment_method: 'stripe',
          stripe_customer_id: session.customer || null,
          stripe_subscription_id: session.subscription || null
        })
        .eq('email', customerEmail)

      if (error) {
        console.error('Error updating user:', error)
        return new Response(JSON.stringify({ error: 'Failed to update user' }), {
          status: 500
        })
      }

      return new Response(JSON.stringify({ 
        success: true, 
        tier: subscriptionTier,
        email: customerEmail 
      }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' }
      })
    }

    // Handle subscription updates
    if (event.type === 'customer.subscription.updated') {
      const subscription = event.data.object
      const customerId = subscription.customer

      // Find user by Stripe customer ID
      const { data: user } = await supabaseAdmin
        .from('unlocked_users')
        .select('email')
        .eq('stripe_customer_id', customerId)
        .single()

      if (user) {
        const expiresAt = new Date(subscription.current_period_end * 1000)
        await supabaseAdmin
          .from('unlocked_users')
          .update({
            subscription_expires_at: expiresAt.toISOString(),
            verified: subscription.status === 'active'
          })
          .eq('email', user.email)
      }
    }

    return new Response(JSON.stringify({ received: true }), {
      status: 200
    })
  } catch (error) {
    return new Response(JSON.stringify({ error: error.message }), {
      status: 500
    })
  }
})
```

### Step 4: Deploy Function

```bash
supabase functions deploy stripe-webhook
```

### Step 5: Set Secrets

```bash
supabase secrets set STRIPE_SECRET_KEY=sk_live_xxxxx
supabase secrets set STRIPE_WEBHOOK_SECRET=whsec_xxxxx
supabase secrets set SUPABASE_URL=https://your-project.supabase.co
supabase secrets set SUPABASE_SERVICE_ROLE_KEY=your_service_role_key
```

---

## Option 3: Simple Node.js Server (If You Have a Server)

Create `webhook-server.js`:

```javascript
const express = require('express');
const stripe = require('stripe')(process.env.STRIPE_SECRET_KEY);
const { createClient } = require('@supabase/supabase-js');

const app = express();
app.use(express.raw({ type: 'application/json' }));

const supabase = createClient(
  process.env.SUPABASE_URL,
  process.env.SUPABASE_SERVICE_ROLE_KEY
);

app.post('/webhook', async (req, res) => {
  const sig = req.headers['stripe-signature'];
  const webhookSecret = process.env.STRIPE_WEBHOOK_SECRET;

  let event;
  try {
    event = stripe.webhooks.constructEvent(req.body, sig, webhookSecret);
  } catch (err) {
    return res.status(400).send(`Webhook Error: ${err.message}`);
  }

  if (event.type === 'checkout.session.completed') {
    const session = event.data.object;
    const email = session.customer_details?.email;
    
    if (!email) return res.status(400).json({ error: 'No email' });

    // Determine tier from amount
    let tier = 'basic';
    if (session.amount_total === 1000) tier = 'pro'; // $10
    
    const expiresAt = new Date();
    expiresAt.setDate(expiresAt.getDate() + 30);

    await supabase
      .from('unlocked_users')
      .update({
        verified: true,
        subscription_tier: tier,
        subscription_expires_at: expiresAt.toISOString(),
        payment_method: 'stripe'
      })
      .eq('email', email);
  }

  res.json({ received: true });
});

app.listen(3000, () => console.log('Webhook server running on port 3000'));
```

---

## 🔧 Setting Up Stripe Webhook

### Step 1: Get Your Webhook Endpoint URL

**For Vercel:** `https://your-project.vercel.app/api/stripe-webhook`  
**For Supabase Edge Function:** `https://your-project.supabase.co/functions/v1/stripe-webhook`  
**For Custom Server:** `https://your-domain.com/webhook`

### Step 2: Configure Webhook in Stripe

1. Go to [Stripe Dashboard](https://dashboard.stripe.com)
2. Click **"Developers"** → **"Webhooks"**
3. Click **"+ Add endpoint"**
4. **Endpoint URL:** Paste your webhook URL
5. **Events to send:** Select:
   - `checkout.session.completed` ✅
   - `customer.subscription.created` ✅
   - `customer.subscription.updated` ✅
   - `customer.subscription.deleted` ✅
6. Click **"Add endpoint"**

### Step 3: Copy Webhook Signing Secret

1. After creating the endpoint, click on it
2. Find **"Signing secret"**
3. Copy it (starts with `whsec_...`)
4. Add to your environment variables as `STRIPE_WEBHOOK_SECRET`

---

## 🔑 Getting Your Keys

### Stripe Keys:
1. **Stripe Dashboard** → **Developers** → **API keys**
2. **Secret key:** `sk_live_...` (use this, not publishable key)
3. **Webhook secret:** From webhook endpoint page

### Supabase Keys:
1. **Supabase Dashboard** → **Project Settings** → **API**
2. **Project URL:** Copy this
3. **Service Role Key:** Copy this (⚠️ Keep secret! Not the anon key)

---

## 🧪 Testing

### Test Mode:
1. Use Stripe test mode
2. Use test webhook endpoint: `https://your-endpoint.com/webhook`
3. Use test card: `4242 4242 4242 4242`
4. Check Stripe Dashboard → **Webhooks** → **Recent events** for logs

### Verify It Works:
1. Subscribe with test account
2. Check Supabase → `unlocked_users` table
3. User should have:
   - `verified = true`
   - `subscription_tier = 'basic'` or `'pro'`
   - `subscription_expires_at` = 30 days from now

---

## 🎯 Which Option Should You Use?

- **Option 1 (Vercel):** Best if you want easy deployment, free hosting
- **Option 2 (Supabase Edge Function):** Best if you want everything in one place, no external services
- **Option 3 (Custom Server):** Best if you already have a server running

**Recommendation:** Use **Option 2 (Supabase Edge Function)** - it's the simplest and keeps everything in Supabase!

---

## 📝 Important Notes

1. **Service Role Key:** Must use service role key (not anon key) to bypass RLS
2. **Webhook Security:** Always verify webhook signatures
3. **Email Matching:** Webhook matches users by email - make sure emails match!
4. **Error Handling:** Log errors and monitor webhook failures in Stripe Dashboard

---

## 🚨 Troubleshooting

### Webhook not receiving events?
- Check Stripe Dashboard → Webhooks → Recent events
- Verify endpoint URL is correct
- Check server logs for errors

### User not getting upgraded?
- Check if email matches exactly (case-sensitive)
- Verify service role key is correct
- Check Supabase logs for update errors

### Wrong tier assigned?
- Check `amount_total` logic matches your Stripe prices
- Verify price IDs in your code match Stripe products

---

Need help setting this up? Let me know which option you prefer and I can provide more specific guidance!
