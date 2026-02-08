// Supabase Edge Function for Stripe Webhooks
// Deploy with: supabase functions deploy stripe-webhook

import { serve } from "https://deno.land/std@0.168.0/http/server.ts"
import { createClient } from 'https://esm.sh/@supabase/supabase-js@2'

// Import Stripe
const Stripe = await import('https://esm.sh/stripe@14.0.0')
const stripe = new Stripe.default(
  Deno.env.get('STRIPE_SECRET_KEY') || '',
  { apiVersion: '2023-10-16' }
)

serve(async (req) => {
  try {
    // Only allow POST requests
    if (req.method !== 'POST') {
      return new Response(JSON.stringify({ error: 'Method not allowed' }), {
        status: 405,
        headers: { 'Content-Type': 'application/json' }
      })
    }

    // Get webhook signature
    const signature = req.headers.get('stripe-signature')
    if (!signature) {
      return new Response(JSON.stringify({ error: 'No signature provided' }), {
        status: 400
      })
    }

    // Get webhook secret
    const webhookSecret = Deno.env.get('STRIPE_WEBHOOK_SECRET')
    if (!webhookSecret) {
      return new Response(JSON.stringify({ error: 'Webhook secret not configured' }), {
        status: 500
      })
    }

    // Get request body
    const body = await req.text()

    // Verify webhook signature
    let event
    try {
      event = stripe.webhooks.constructEvent(body, signature, webhookSecret)
    } catch (err) {
      console.error('Webhook signature verification failed:', err.message)
      return new Response(JSON.stringify({ error: `Webhook Error: ${err.message}` }), {
        status: 400,
        headers: { 'Content-Type': 'application/json' }
      })
    }

    // Initialize Supabase admin client (bypasses RLS)
    // Get Supabase URL from environment (set automatically by Supabase)
    const supabaseUrl = Deno.env.get('SUPABASE_URL') || Deno.env.get('SUPABASE_PROJECT_URL') || 'https://wbpfuuiznsmysbskywdx.supabase.co'
    // Use SERVICE_ROLE_KEY (not SUPABASE_SERVICE_ROLE_KEY - Supabase CLI blocks SUPABASE_ prefix)
    const supabaseServiceKey = Deno.env.get('SERVICE_ROLE_KEY')
    
    if (!supabaseUrl || !supabaseServiceKey) {
      return new Response(JSON.stringify({ error: 'Supabase credentials not configured' }), {
        status: 500
      })
    }

    const supabaseAdmin = createClient(supabaseUrl, supabaseServiceKey)

    console.log('Received webhook event:', event.type)

    // Handle checkout.session.completed (when user completes payment)
    if (event.type === 'checkout.session.completed') {
      const session = event.data.object
      
      // Get customer email
      const customerEmail = session.customer_details?.email || session.customer_email
      
      if (!customerEmail) {
        console.error('No email found in checkout session')
        return new Response(JSON.stringify({ error: 'No email found in session' }), {
          status: 400
        })
      }

      console.log(`Processing subscription for: ${customerEmail}`)

      // Determine subscription tier from amount (in cents)
      // Your Stripe links:
      // Basic: $5.00 = 500 cents
      // Pro: $10.00 = 1000 cents
      let subscriptionTier = 'basic' // Default
      
      if (session.amount_total === 1000) {
        subscriptionTier = 'pro'
      } else if (session.amount_total === 500) {
        subscriptionTier = 'basic'
      }

      console.log(`Detected tier: ${subscriptionTier} (amount: ${session.amount_total} cents)`)

      // Calculate expiry date (30 days from now for monthly subscriptions)
      const expiresAt = new Date()
      expiresAt.setDate(expiresAt.getDate() + 30)

      // Update user in Supabase
      const { data, error } = await supabaseAdmin
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
        return new Response(JSON.stringify({ 
          error: 'Failed to update user',
          details: error.message 
        }), {
          status: 500
        })
      }

      console.log(`✅ Successfully activated ${subscriptionTier} subscription for ${customerEmail}`)
      
      return new Response(JSON.stringify({ 
        success: true,
        tier: subscriptionTier,
        email: customerEmail,
        expiresAt: expiresAt.toISOString()
      }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' }
      })
    }

    // Handle subscription updates (renewals, changes)
    if (event.type === 'customer.subscription.updated') {
      const subscription = event.data.object
      const customerId = subscription.customer

      console.log(`Processing subscription update for customer: ${customerId}`)

      // Find user by Stripe customer ID
      const { data: user, error: findError } = await supabaseAdmin
        .from('unlocked_users')
        .select('email, subscription_tier')
        .eq('stripe_customer_id', customerId)
        .maybeSingle()

      if (findError || !user) {
        console.error('User not found for customer ID:', customerId)
        return new Response(JSON.stringify({ error: 'User not found' }), {
          status: 404
        })
      }

      // Update subscription expiry
      const expiresAt = new Date(subscription.current_period_end * 1000)
      const isActive = subscription.status === 'active' || subscription.status === 'trialing'

      const { error: updateError } = await supabaseAdmin
        .from('unlocked_users')
        .update({
          subscription_expires_at: expiresAt.toISOString(),
          verified: isActive
        })
        .eq('email', user.email)

      if (updateError) {
        console.error('Error updating subscription:', updateError)
        return new Response(JSON.stringify({ error: 'Failed to update subscription' }), {
          status: 500
        })
      }

      console.log(`✅ Updated subscription for ${user.email}`)
      return new Response(JSON.stringify({ success: true }), {
        status: 200
      })
    }

    // Handle subscription cancellations
    if (event.type === 'customer.subscription.deleted') {
      const subscription = event.data.object
      const customerId = subscription.customer

      console.log(`Processing subscription cancellation for customer: ${customerId}`)

      const { data: user } = await supabaseAdmin
        .from('unlocked_users')
        .select('email')
        .eq('stripe_customer_id', customerId)
        .maybeSingle()

      if (user) {
        await supabaseAdmin
          .from('unlocked_users')
          .update({
            verified: false,
            subscription_tier: 'trial',
            subscription_expires_at: null
          })
          .eq('email', user.email)

        console.log(`✅ Cancelled subscription for ${user.email}`)
      }

      return new Response(JSON.stringify({ success: true }), {
        status: 200
      })
    }

    // Acknowledge other events
    return new Response(JSON.stringify({ received: true }), {
      status: 200
    })

  } catch (error) {
    console.error('Webhook error:', error)
    return new Response(JSON.stringify({ 
      error: error.message,
      stack: error.stack 
    }), {
      status: 500
    })
  }
})
