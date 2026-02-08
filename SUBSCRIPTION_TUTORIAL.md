# Complete Tutorial: Implementing Subscriptions & Discord Notifier

## Part 1: Implementing Subscriptions ($5/month)

### Option A: Stripe Integration (Recommended - Easiest)

#### Step 1: Create Stripe Account

1. Go to https://stripe.com
2. Sign up for an account
3. Get your API keys:
   - Go to Developers → API keys
   - Copy your **Publishable key** (starts with `pk_`)
   - Copy your **Secret key** (starts with `sk_`) - Keep this secret!

#### Step 2: Create a Product and Price in Stripe

1. In Stripe Dashboard, go to **Products**
2. Click **+ Add product**
3. Fill in:
   - **Name:** Discord Notifications
   - **Description:** Monthly subscription for Discord notifications
   - **Pricing:** 
     - Type: Recurring
     - Price: $5.00 USD
     - Billing period: Monthly
4. Click **Save product**
5. Copy the **Price ID** (starts with `price_`) - You'll need this!

#### Step 3: Add Stripe to Your Website

You have two options:

**Option 3A: Stripe Checkout (Easiest - No backend needed)**

Add this to your HTML (create a subscription button/page):

```html
<!-- Add this button somewhere in your website -->
<button id="subscribe-discord-stripe-btn" style="background: #635BFF; color: white; padding: 12px 24px; border: none; border-radius: 8px; cursor: pointer; font-weight: bold;">
  Subscribe to Discord Notifications - $5/month
</button>

<!-- Add Stripe.js library -->
<script src="https://js.stripe.com/v3/"></script>
```

Add this JavaScript:

```javascript
// Stripe Subscription Handler
(function() {
  const SUPABASE_URL = 'https://wbpfuuiznsmysbskywdx.supabase.co';
  const SUPABASE_ANON_KEY = 'sb_publishable_rIy_-DWT87Gj9ao1WvN3gA_WA6eME-x';
  const STRIPE_PUBLISHABLE_KEY = 'pk_test_...'; // Your Stripe publishable key
  const STRIPE_PRICE_ID = 'price_...'; // Your Stripe price ID from Step 2
  
  const stripe = Stripe(STRIPE_PUBLISHABLE_KEY);
  const supabase = window.supabase ? window.supabase.createClient(SUPABASE_URL, SUPABASE_ANON_KEY) : null;
  
  // Subscribe button handler
  const subscribeBtn = document.getElementById('subscribe-discord-stripe-btn');
  if (subscribeBtn) {
    subscribeBtn.addEventListener('click', async function() {
      if (!supabase) {
        alert('Database connection not available.');
        return;
      }
      
      try {
        // Get current user
        const { data: { session } } = await supabase.auth.getSession();
        if (!session) {
          alert('Please log in first.');
          return;
        }
        
        const { data: { user } } = await supabase.auth.getUser();
        if (!user) return;
        
        // Create Stripe Checkout Session
        // NOTE: You need a backend endpoint for this! See Option 3B below.
        // For now, redirect to Stripe payment link (see Step 3B alternative)
        
        alert('Payment integration requires a backend. See setup instructions for details.');
      } catch (err) {
        console.error('Error:', err);
        alert('Error initiating subscription. Please try again.');
      }
    });
  }
})();
```

**Option 3B: Stripe Payment Links (Simplest - No code needed!)**

1. In Stripe Dashboard, go to **Products**
2. Click on your "Discord Notifications" product
3. Click **Create payment link**
4. Set:
   - Price: Your $5/month price
   - Customer email: Collect email
   - After payment: Redirect to your website success page
5. Copy the payment link URL
6. Add this to your website:

```html
<a href="YOUR_STRIPE_PAYMENT_LINK_HERE" 
   target="_blank"
   style="background: #635BFF; color: white; padding: 12px 24px; border-radius: 8px; text-decoration: none; display: inline-block; font-weight: bold;">
  Subscribe to Discord Notifications - $5/month
</a>
```

#### Step 4: Handle Successful Payment

When payment succeeds, you need to activate the subscription. You have two options:

**Option 4A: Stripe Webhook (Automatic - Recommended)**

1. **Set up a webhook endpoint** (requires a backend server):
   - In Stripe Dashboard → Developers → Webhooks
   - Click **Add endpoint**
   - Endpoint URL: `https://your-domain.com/api/stripe-webhook`
   - Select events: `checkout.session.completed`, `customer.subscription.created`, `invoice.payment_succeeded`
   - Copy the webhook signing secret

2. **Create webhook handler** (Node.js example):

```javascript
// server.js (Node.js/Express example)
const express = require('express');
const stripe = require('stripe')('sk_...'); // Your secret key
const { createClient } = require('@supabase/supabase-js');

const app = express();
const supabase = createClient(SUPABASE_URL, SUPABASE_SERVICE_KEY); // Use service key, not anon key

app.post('/api/stripe-webhook', express.raw({type: 'application/json'}), async (req, res) => {
  const sig = req.headers['stripe-signature'];
  const webhookSecret = 'whsec_...'; // Your webhook secret
  
  let event;
  try {
    event = stripe.webhooks.constructEvent(req.body, sig, webhookSecret);
  } catch (err) {
    return res.status(400).send(`Webhook Error: ${err.message}`);
  }
  
  // Handle the event
  if (event.type === 'checkout.session.completed') {
    const session = event.data.object;
    const customerEmail = session.customer_details.email;
    
    // Find user by email and activate subscription
    const expiresAt = new Date();
    expiresAt.setMonth(expiresAt.getMonth() + 1);
    
    await supabase
      .from('unlocked_users')
      .update({
        discord_subscription_active: true,
        discord_subscription_expires_at: expiresAt.toISOString()
      })
      .eq('email', customerEmail);
  }
  
  res.json({received: true});
});

app.listen(3000);
```

**Option 4B: Manual Activation Page (Simpler - For testing)**

Create a success page that users land on after payment:

```html
<!-- subscription-success.html -->
<!DOCTYPE html>
<html>
<head>
  <title>Subscription Activated</title>
</head>
<body>
  <h1>Thank you for subscribing!</h1>
  <p>Your Discord notifications subscription is being activated...</p>
  
  <script>
    // Get user email from URL or session
    const urlParams = new URLSearchParams(window.location.search);
    const email = urlParams.get('email') || prompt('Enter your email:');
    
    if (email) {
      // Call Supabase to activate (you'll need to expose this via a function or do it manually)
      fetch('YOUR_BACKEND_ENDPOINT/activate-subscription', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email: email })
      }).then(() => {
        alert('Subscription activated! You can now use Discord notifications.');
        window.location.href = '/index.html';
      });
    }
  </script>
</body>
</html>
```

**Option 4C: Supabase Edge Function (No separate server needed!)**

1. In Supabase Dashboard → Edge Functions → Create function
2. Name it `activate-discord-subscription`
3. Code:

```typescript
// supabase/functions/activate-discord-subscription/index.ts
import { serve } from "https://deno.land/std@0.168.0/http/server.ts"
import { createClient } from 'https://esm.sh/@supabase/supabase-js@2'

const corsHeaders = {
  'Access-Control-Allow-Origin': '*',
  'Access-Control-Allow-Headers': 'authorization, x-client-info, apikey, content-type',
}

serve(async (req) => {
  if (req.method === 'OPTIONS') {
    return new Response('ok', { headers: corsHeaders })
  }

  try {
    const { email } = await req.json()
    
    const supabaseClient = createClient(
      Deno.env.get('SUPABASE_URL') ?? '',
      Deno.env.get('SUPABASE_SERVICE_ROLE_KEY') ?? ''
    )
    
    const expiresAt = new Date()
    expiresAt.setMonth(expiresAt.getMonth() + 1)
    
    const { error } = await supabaseClient
      .from('unlocked_users')
      .update({
        discord_subscription_active: true,
        discord_subscription_expires_at: expiresAt.toISOString()
      })
      .eq('email', email)
    
    if (error) throw error
    
    return new Response(
      JSON.stringify({ success: true }),
      { headers: { ...corsHeaders, 'Content-Type': 'application/json' }, status: 200 }
    )
  } catch (error) {
    return new Response(
      JSON.stringify({ error: error.message }),
      { headers: { ...corsHeaders, 'Content-Type': 'application/json' }, status: 400 }
    )
  }
})
```

4. Deploy: `supabase functions deploy activate-discord-subscription`
5. Call it from your success page

### Option B: PayPal Integration

1. **Create PayPal Business Account**
   - Go to https://www.paypal.com/business
   - Sign up and verify your account

2. **Create Subscription Button**
   - PayPal Dashboard → Products and Services → Subscriptions
   - Create new subscription: $5/month
   - Copy the subscription button code
   - Add to your website

3. **Handle IPN (Instant Payment Notification)**
   - Set up IPN listener to receive payment confirmations
   - Update database when payment received

### Option C: Manual Activation (For Testing/Development)

For testing, you can manually activate subscriptions:

```sql
-- Activate subscription for a user
UPDATE unlocked_users 
SET 
  discord_subscription_active = TRUE,
  discord_subscription_expires_at = NOW() + INTERVAL '30 days'
WHERE email = 'user@example.com';
```

## Part 2: Setting Up the Golang Notifier Service

### Step 1: Install Go

1. **Download Go:**
   - Go to https://go.dev/dl/
   - Download for your OS (Windows/Mac/Linux)
   - Install it

2. **Verify installation:**
```bash
go version
```

### Step 2: Create Project Structure

```bash
# Create project directory
mkdir discord-notifier
cd discord-notifier

# Initialize Go module
go mod init discord-notifier

# Create main file
touch main.go
```

### Step 3: Install Dependencies

```bash
go get github.com/supabase-community/supabase-go
go get github.com/joho/godotenv  # For .env file support
```

### Step 4: Create main.go

Create `main.go` with this complete code:

```go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	supabase "github.com/supabase-community/supabase-go"
)

// Notification represents a single notification rule
type Notification struct {
	ID         string   `json:"id"`
	SearchTerm string   `json:"searchTerm"`
	Markets    []string `json:"markets"` // null = all markets
	CreatedAt  string   `json:"createdAt"`
}

// User represents a user with Discord notifications
type User struct {
	AuthUserID            string         `json:"auth_user_id"`
	Email                 string         `json:"email"`
	Username              string         `json:"username"`
	DiscordWebhookURL     string         `json:"discord_webhook_url"`
	DiscordNotifications  []Notification `json:"discord_notifications"`
	SubscriptionActive    bool           `json:"discord_subscription_active"`
	SubscriptionExpiresAt *string        `json:"discord_subscription_expires_at"`
}

// DiscordEmbed represents a Discord embed
type DiscordEmbed struct {
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	URL         string                 `json:"url,omitempty"`
	Color       int                    `json:"color"`
	Fields      []DiscordEmbedField    `json:"fields,omitempty"`
	Timestamp   string                 `json:"timestamp,omitempty"`
	Footer      map[string]interface{} `json:"footer,omitempty"`
}

type DiscordEmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

// DiscordWebhookPayload represents the payload sent to Discord
type DiscordWebhookPayload struct {
	Content string        `json:"content,omitempty"`
	Embeds  []DiscordEmbed `json:"embeds,omitempty"`
}

var (
	supabaseURL  string
	supabaseKey  string
	pollInterval time.Duration = 10 * time.Minute
)

func main() {
	// Load environment variables from .env file
	godotenv.Load()

	// Get configuration from environment
	supabaseURL = os.Getenv("SUPABASE_URL")
	supabaseKey = os.Getenv("SUPABASE_ANON_KEY")

	if supabaseURL == "" {
		supabaseURL = "https://wbpfuuiznsmysbskywdx.supabase.co"
	}
	if supabaseKey == "" {
		log.Fatal("SUPABASE_ANON_KEY environment variable is required")
	}

	// Parse poll interval
	if intervalStr := os.Getenv("POLL_INTERVAL"); intervalStr != "" {
		if duration, err := time.ParseDuration(intervalStr); err == nil {
			pollInterval = duration
		}
	}

	log.Printf("🚀 Starting Discord Notifier Service")
	log.Printf("📡 Supabase URL: %s", supabaseURL)
	log.Printf("⏱️  Poll interval: %v", pollInterval)
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Run immediately on startup
	log.Printf("🔄 Running initial check...")
	processAllNotifications()

	// Then run on interval
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		<-ticker.C
		log.Printf("🔄 Running scheduled check...")
		processAllNotifications()
	}
}

// processAllNotifications fetches all active subscribers and processes their notifications
func processAllNotifications() {
	users, err := fetchActiveSubscribers()
	if err != nil {
		log.Printf("❌ Error fetching subscribers: %v", err)
		return
	}

	if len(users) == 0 {
		log.Printf("ℹ️  No active subscribers found")
		return
	}

	log.Printf("✅ Found %d active subscriber(s)", len(users))

	for _, user := range users {
		if !isSubscriptionActive(user) {
			log.Printf("⏭️  Skipping user %s - subscription expired or inactive", user.Email)
			continue
		}

		log.Printf("👤 Processing notifications for: %s (%s)", user.Username, user.Email)
		processUserNotifications(user)
	}

	log.Printf("✅ Finished processing all notifications")
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

// fetchActiveSubscribers fetches all users with active subscriptions and webhooks
func fetchActiveSubscribers() ([]User, error) {
	client, err := supabase.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Supabase client: %w", err)
	}

	// Query for users with active subscriptions and webhooks
	var users []User
	err = client.DB.From("unlocked_users").
		Select("auth_user_id,email,username,discord_webhook_url,discord_notifications,discord_subscription_active,discord_subscription_expires_at").
		Eq("discord_subscription_active", "true").
		Not("discord_webhook_url", "is", "null").
		ExecuteTo(&users)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	// Parse notifications JSON for each user
	for i := range users {
		if users[i].DiscordNotifications == nil {
			users[i].DiscordNotifications = []Notification{}
		}
	}

	return users, nil
}

// isSubscriptionActive checks if a user's subscription is still valid
func isSubscriptionActive(user User) bool {
	if !user.SubscriptionActive {
		return false
	}

	// If no expiry date, assume lifetime subscription
	if user.SubscriptionExpiresAt == nil || *user.SubscriptionExpiresAt == "" {
		return true
	}

	expiresAt, err := time.Parse(time.RFC3339, *user.SubscriptionExpiresAt)
	if err != nil {
		log.Printf("⚠️  Error parsing expiry date for user %s: %v", user.Email, err)
		return false
	}

	return time.Now().Before(expiresAt)
}

// processUserNotifications processes all notifications for a single user
func processUserNotifications(user User) {
	if len(user.DiscordNotifications) == 0 {
		log.Printf("   ℹ️  No notifications configured for this user")
		return
	}

	log.Printf("   📋 Processing %d notification(s)", len(user.DiscordNotifications))

	for _, notif := range user.DiscordNotifications {
		log.Printf("   🔍 Checking: '%s'", notif.SearchTerm)

		// TODO: Implement your actual search logic here
		// This is where you would:
		// 1. Determine which markets to check (notif.Markets or all markets)
		// 2. Search each market for the search term
		// 3. Compare with previously found items (to avoid duplicates)
		// 4. If new items found, call sendDiscordNotification

		// Example: Mock new items (replace with actual search)
		newItems := []map[string]interface{}{
			{
				"title":       "Vintage Nike Sneakers",
				"description": "Great condition vintage Nike sneakers",
				"url":         "https://example.com/item/123",
				"price":       "$50",
				"market":      "Mercari Japan",
			},
		}

		// Only send if there are new items
		if len(newItems) > 0 {
			log.Printf("   ✅ Found %d new item(s) for '%s'", len(newItems), notif.SearchTerm)
			if err := sendDiscordNotification(user.DiscordWebhookURL, notif, newItems); err != nil {
				log.Printf("   ❌ Error sending notification: %v", err)
			} else {
				log.Printf("   ✅ Notification sent successfully")
			}
		} else {
			log.Printf("   ℹ️  No new items found")
		}
	}
}

// sendDiscordNotification sends a notification to Discord webhook
func sendDiscordNotification(webhookURL string, notification Notification, items []map[string]interface{}) error {
	if webhookURL == "" {
		return fmt.Errorf("webhook URL is empty")
	}

	// Build embeds for each item
	embeds := []DiscordEmbed{}
	for _, item := range items {
		embed := DiscordEmbed{
			Title:       fmt.Sprintf("New: %s", notification.SearchTerm),
			Description: getString(item, "description", "New item found!"),
			URL:         getString(item, "url", ""),
			Color:       3447003, // Blue color
			Timestamp:  time.Now().Format(time.RFC3339),
			Footer: map[string]interface{}{
				"text": "MMCS Notifications",
			},
		}

		// Add fields
		if price := getString(item, "price", ""); price != "" {
			embed.Fields = append(embed.Fields, DiscordEmbedField{
				Name:   "Price",
				Value:  price,
				Inline: true,
			})
		}

		if market := getString(item, "market", ""); market != "" {
			embed.Fields = append(embed.Fields, DiscordEmbedField{
				Name:   "Market",
				Value:  market,
				Inline: true,
			})
		}

		embeds = append(embeds, embed)
	}

	// Build payload
	payload := DiscordWebhookPayload{
		Content: fmt.Sprintf("🔔 **%d new item(s) found for: %s**", len(items), notification.SearchTerm),
		Embeds:  embeds,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Send to Discord
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	return nil
}

// Helper function to safely get string from map
func getString(m map[string]interface{}, key string, defaultValue string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return defaultValue
}
```

### Step 5: Create .env File

Create a `.env` file in your project directory:

```bash
# .env
SUPABASE_URL=https://wbpfuuiznsmysbskywdx.supabase.co
SUPABASE_ANON_KEY=your_anon_key_here
POLL_INTERVAL=10m
```

### Step 6: Test Locally

```bash
# Make sure you're in the project directory
cd discord-notifier

# Run the service
go run main.go
```

You should see:
```
🚀 Starting Discord Notifier Service
📡 Supabase URL: https://wbpfuuiznsmysbskywdx.supabase.co
⏱️  Poll interval: 10m0s
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🔄 Running initial check...
✅ Found 1 active subscriber(s)
👤 Processing notifications for: username (user@example.com)
   📋 Processing 2 notification(s)
   🔍 Checking: 'vintage nike'
   ✅ Found 1 new item(s) for 'vintage nike'
   ✅ Notification sent successfully
```

### Step 7: Implement Your Search Logic

Replace the mock code in `processUserNotifications` with your actual search implementation:

```go
// Example: Search a market
func searchMarket(marketKey string, searchTerm string) ([]map[string]interface{}, error) {
    // Your search logic here
    // This could involve:
    // - Scraping websites
    // - Using APIs
    // - Parsing HTML
    // - etc.
    
    return items, nil
}
```

### Step 8: Add Item Tracking (Avoid Duplicates)

Create a simple file-based or database-based tracker:

```go
// Track notified items to avoid duplicates
type ItemTracker struct {
    items map[string]time.Time // item URL -> first seen time
}

func (t *ItemTracker) IsNew(itemURL string) bool {
    _, exists := t.items[itemURL]
    if !exists {
        t.items[itemURL] = time.Now()
    }
    return !exists
}

// Or use a database table to track notified items
```

### Step 9: Deploy the Service

**Option A: Run on Your Computer (Development)**
```bash
go run main.go
```

**Option B: Build and Run**
```bash
go build -o discord-notifier main.go
./discord-notifier
```

**Option C: Windows Service**
```powershell
# Install NSSM from https://nssm.cc/download
nssm install DiscordNotifier "C:\path\to\discord-notifier.exe"
nssm set DiscordNotifier AppEnvironmentExtra SUPABASE_URL=https://... SUPABASE_ANON_KEY=...
nssm start DiscordNotifier
```

**Option D: Linux Systemd**
```bash
# Create /etc/systemd/system/discord-notifier.service
sudo nano /etc/systemd/system/discord-notifier.service
```

```ini
[Unit]
Description=Discord Notifier Service
After=network.target

[Service]
Type=simple
User=your-user
WorkingDirectory=/path/to/discord-notifier
Environment="SUPABASE_URL=https://wbpfuuiznsmysbskywdx.supabase.co"
Environment="SUPABASE_ANON_KEY=your_key"
ExecStart=/path/to/discord-notifier/discord-notifier
Restart=always

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl enable discord-notifier
sudo systemctl start discord-notifier
```

**Option E: Docker**
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o discord-notifier .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/discord-notifier .
CMD ["./discord-notifier"]
```

```bash
docker build -t discord-notifier .
docker run -d --name discord-notifier --restart unless-stopped \
  -e SUPABASE_URL=https://wbpfuuiznsmysbskywdx.supabase.co \
  -e SUPABASE_ANON_KEY=your_key \
  discord-notifier
```

## Complete Setup Checklist

### Database Setup
- [ ] Run SQL to add columns (discord_webhook_url, discord_notifications, discord_subscription_active, discord_subscription_expires_at)
- [ ] Verify columns exist in Supabase table editor

### Payment Setup
- [ ] Choose payment method (Stripe/PayPal/Manual)
- [ ] Set up payment processor account
- [ ] Create $5/month subscription product
- [ ] Add payment button/link to website
- [ ] Set up webhook/activation handler
- [ ] Test payment flow

### Golang Service Setup
- [ ] Install Go
- [ ] Create project directory
- [ ] Copy main.go code
- [ ] Create .env file with credentials
- [ ] Install dependencies (`go mod tidy`)
- [ ] Test locally (`go run main.go`)
- [ ] Implement search logic
- [ ] Add item tracking
- [ ] Deploy service (choose method)
- [ ] Set up monitoring/logging

### Testing
- [ ] Manually activate test subscription
- [ ] Add test Discord webhook
- [ ] Add test notification
- [ ] Verify service fetches user
- [ ] Verify service sends Discord message
- [ ] Test payment flow end-to-end

## Next Steps After Setup

1. **Implement actual search logic** - Replace mock code with real market searches
2. **Add item tracking** - Store notified items to avoid duplicates
3. **Set up logging** - Log to file or service like Loggly
4. **Monitor service** - Set up alerts if service stops
5. **Handle renewals** - Automatically extend subscriptions on payment
6. **Add admin dashboard** - View/manage subscriptions

## Common Issues & Solutions

**Issue:** Service can't connect to Supabase
- **Solution:** Check API key, verify network connectivity

**Issue:** No users found
- **Solution:** Verify users have `discord_subscription_active = TRUE` and webhook URL set

**Issue:** Discord webhook fails
- **Solution:** Verify webhook URL is correct, check Discord server permissions

**Issue:** Duplicate notifications
- **Solution:** Implement item tracking (store notified item URLs)

**Issue:** Service stops running
- **Solution:** Use systemd/Docker with restart policies, add health checks
