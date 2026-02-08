# Complete Setup Instructions for Discord Notifications ($5/month)

## Step 1: Database Setup (Supabase)

1. **Go to your Supabase Dashboard**
   - Navigate to: https://supabase.com/dashboard
   - Select your project

2. **Open SQL Editor**
   - Click on "SQL Editor" in the left sidebar
   - Click "New query"

3. **Run this SQL:**

```sql
-- Add Discord notification and subscription columns
ALTER TABLE unlocked_users 
ADD COLUMN IF NOT EXISTS discord_webhook_url TEXT,
ADD COLUMN IF NOT EXISTS discord_notifications JSONB,
ADD COLUMN IF NOT EXISTS discord_subscription_active BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS discord_subscription_expires_at TIMESTAMPTZ;

-- Add index for faster queries
CREATE INDEX IF NOT EXISTS idx_discord_subscription_active 
ON unlocked_users (discord_subscription_active, discord_subscription_expires_at) 
WHERE discord_subscription_active = TRUE;

-- Add comment for documentation
COMMENT ON COLUMN unlocked_users.discord_webhook_url IS 'Discord webhook URL for notifications';
COMMENT ON COLUMN unlocked_users.discord_notifications IS 'JSON array of notification rules';
COMMENT ON COLUMN unlocked_users.discord_subscription_active IS 'Whether user has active $5/month subscription';
COMMENT ON COLUMN unlocked_users.discord_subscription_expires_at IS 'When subscription expires (null = lifetime)';
```

4. **Verify the columns were added:**
   - Go to "Table Editor" → `unlocked_users`
   - You should see the new columns: `discord_webhook_url`, `discord_notifications`, `discord_subscription_active`, `discord_subscription_expires_at`

## Step 2: Set Up Payment Processing (Choose One)

### Option A: Stripe (Recommended)

1. **Create Stripe Account**
   - Go to: https://stripe.com
   - Sign up and get your API keys

2. **Add Stripe to your website:**
   - You'll need to add a payment button/page
   - When payment succeeds, update `discord_subscription_active = TRUE` and set `discord_subscription_expires_at` to 30 days from now

3. **Stripe Webhook Handler:**
   - Set up a webhook endpoint to handle subscription renewals
   - Update `discord_subscription_expires_at` when payment is received

### Option B: PayPal

1. **Create PayPal Business Account**
   - Go to: https://www.paypal.com/business
   - Set up recurring payments

2. **Add PayPal Button:**
   - Create a subscription button for $5/month
   - Handle IPN (Instant Payment Notification) to update subscription status

### Option C: Manual Activation (For Testing)

For testing, you can manually activate subscriptions in Supabase:

```sql
-- Activate subscription for a user (replace USER_ID with actual auth_user_id)
UPDATE unlocked_users 
SET 
  discord_subscription_active = TRUE,
  discord_subscription_expires_at = NOW() + INTERVAL '30 days'
WHERE auth_user_id = 'USER_ID_HERE';
```

## Step 3: Update Website Code (Payment Integration)

You need to add a payment page/button. Here's what to add:

1. **Create a subscription page** (or add to existing page):
   - Add a "Subscribe to Discord Notifications" button
   - Link to your payment processor (Stripe/PayPal)
   - After successful payment, update the database

2. **Example: Manual activation link** (for testing):
   Add this somewhere in your code (temporary, for testing):

```javascript
// Temporary admin function to activate subscription
async function activateDiscordSubscription(userId, months = 1) {
  const expiresAt = new Date();
  expiresAt.setMonth(expiresAt.getMonth() + months);
  
  const { error } = await supabase
    .from('unlocked_users')
    .update({
      discord_subscription_active: true,
      discord_subscription_expires_at: expiresAt.toISOString()
    })
    .eq('auth_user_id', userId);
  
  return !error;
}
```

## Step 4: Set Up Your Golang Notifier Service

### 4.1 Create Project Structure

```bash
mkdir discord-notifier
cd discord-notifier
go mod init discord-notifier
```

### 4.2 Install Dependencies

```bash
go get github.com/supabase-community/supabase-go
go get github.com/bwmarrin/discordgo  # If you need Discord SDK
```

### 4.3 Create main.go

Create a file `main.go` with this structure:

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"
    
    "github.com/supabase-community/supabase-go"
)

type Notification struct {
    ID         string   `json:"id"`
    SearchTerm string   `json:"searchTerm"`
    Markets    []string `json:"markets"`
    CreatedAt  string   `json:"createdAt"`
}

type User struct {
    AuthUserID            string         `json:"auth_user_id"`
    Email                 string         `json:"email"`
    Username              string         `json:"username"`
    DiscordWebhookURL     string         `json:"discord_webhook_url"`
    DiscordNotifications  []Notification `json:"discord_notifications"`
    SubscriptionActive    bool           `json:"discord_subscription_active"`
    SubscriptionExpiresAt string         `json:"discord_subscription_expires_at"`
}

var (
    supabaseURL    = os.Getenv("SUPABASE_URL")
    supabaseKey    = os.Getenv("SUPABASE_ANON_KEY")
    pollInterval   = 10 * time.Minute // Default: check every 10 minutes
)

func main() {
    // Load environment variables
    if supabaseURL == "" {
        supabaseURL = "https://wbpfuuiznsmysbskywdx.supabase.co"
    }
    if supabaseKey == "" {
        log.Fatal("SUPABASE_ANON_KEY environment variable is required")
    }
    
    // Parse poll interval from env
    if intervalStr := os.Getenv("POLL_INTERVAL"); intervalStr != "" {
        if duration, err := time.ParseDuration(intervalStr); err == nil {
            pollInterval = duration
        }
    }
    
    log.Printf("Starting Discord Notifier Service")
    log.Printf("Poll interval: %v", pollInterval)
    
    // Run immediately, then on interval
    processAllNotifications()
    
    ticker := time.NewTicker(pollInterval)
    defer ticker.Stop()
    
    for {
        <-ticker.C
        processAllNotifications()
    }
}

func processAllNotifications() {
    users, err := fetchActiveSubscribers()
    if err != nil {
        log.Printf("Error fetching subscribers: %v", err)
        return
    }
    
    log.Printf("Found %d active subscribers", len(users))
    
    for _, user := range users {
        if !isSubscriptionActive(user) {
            log.Printf("Skipping user %s - subscription expired", user.Email)
            continue
        }
        
        processUserNotifications(user)
    }
}

func fetchActiveSubscribers() ([]User, error) {
    client := supabase.CreateClient(supabaseURL, supabaseKey)
    
    // Query for users with active subscriptions and webhooks
    response, err := client.From("unlocked_users").
        Select("auth_user_id,email,username,discord_webhook_url,discord_notifications,discord_subscription_active,discord_subscription_expires_at").
        Eq("discord_subscription_active", "true").
        Not("discord_webhook_url", "is", "null").
        Execute()
    
    if err != nil {
        return nil, fmt.Errorf("failed to fetch: %w", err)
    }
    
    var users []User
    if err := json.Unmarshal(response, &users); err != nil {
        return nil, fmt.Errorf("failed to parse response: %w", err)
    }
    
    // Parse notifications JSON
    for i := range users {
        if users[i].DiscordNotifications == nil {
            users[i].DiscordNotifications = []Notification{}
        }
    }
    
    return users, nil
}

func isSubscriptionActive(user User) bool {
    if !user.SubscriptionActive {
        return false
    }
    
    if user.SubscriptionExpiresAt == "" {
        return true // Lifetime subscription
    }
    
    expiresAt, err := time.Parse(time.RFC3339, user.SubscriptionExpiresAt)
    if err != nil {
        log.Printf("Error parsing expiry date for user %s: %v", user.Email, err)
        return false
    }
    
    return time.Now().Before(expiresAt)
}

func processUserNotifications(user User) {
    for _, notif := range user.DiscordNotifications {
        // TODO: Implement your search logic here
        // Check markets, search for term, find new items
        // Then send Discord webhook if new items found
        
        log.Printf("Processing notification: %s for user %s", notif.SearchTerm, user.Email)
        
        // Example: Send test notification (replace with actual search logic)
        // sendDiscordNotification(user.DiscordWebhookURL, notif, newItems)
    }
}

func sendDiscordNotification(webhookURL string, notification Notification, items []map[string]interface{}) error {
    // Build Discord embed
    embeds := []map[string]interface{}{}
    
    for _, item := range items {
        embed := map[string]interface{}{
            "title":       fmt.Sprintf("New: %s", notification.SearchTerm),
            "description": item["description"].(string),
            "url":         item["url"].(string),
            "color":       3447003, // Blue color
            "fields": []map[string]interface{}{
                {
                    "name":   "Price",
                    "value":  item["price"].(string),
                    "inline": true,
                },
                {
                    "name":   "Market",
                    "value":  item["market"].(string),
                    "inline": true,
                },
            },
            "timestamp": time.Now().Format(time.RFC3339),
            "footer": map[string]interface{}{
                "text": "MMCS Notifications",
            },
        }
        embeds = append(embeds, embed)
    }
    
    payload := map[string]interface{}{
        "content": fmt.Sprintf("🔔 **%d new item(s) found for: %s**", len(items), notification.SearchTerm),
        "embeds":  embeds,
    }
    
    jsonData, _ := json.Marshal(payload)
    
    resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return fmt.Errorf("failed to send webhook: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != 204 {
        return fmt.Errorf("webhook returned status %d", resp.StatusCode)
    }
    
    return nil
}
```

### 4.4 Create .env file

```bash
# .env
SUPABASE_URL=https://wbpfuuiznsmysbskywdx.supabase.co
SUPABASE_ANON_KEY=your_anon_key_here
POLL_INTERVAL=10m
LOG_LEVEL=info
```

### 4.5 Create go.mod dependencies

```bash
go mod tidy
```

## Step 5: Run the Golang Service

### Option A: Local Development

```bash
# Load environment variables
export SUPABASE_URL=https://wbpfuuiznsmysbskywdx.supabase.co
export SUPABASE_ANON_KEY=your_anon_key_here
export POLL_INTERVAL=10m

# Run
go run main.go
```

### Option B: Build and Run

```bash
# Build
go build -o discord-notifier main.go

# Run
./discord-notifier
```

### Option C: Run as Windows Service

1. **Install NSSM (Non-Sucking Service Manager):**
   - Download from: https://nssm.cc/download
   - Extract to a folder (e.g., `C:\nssm`)

2. **Install the service:**
```powershell
cd C:\nssm\win64
.\nssm.exe install DiscordNotifier "C:\path\to\discord-notifier.exe"
.\nssm.exe set DiscordNotifier AppEnvironmentExtra SUPABASE_URL=https://wbpfuuiznsmysbskywdx.supabase.co SUPABASE_ANON_KEY=your_key_here
.\nssm.exe start DiscordNotifier
```

### Option D: Run as Linux Systemd Service

1. **Create service file:** `/etc/systemd/system/discord-notifier.service`

```ini
[Unit]
Description=MMCS Discord Notifier Service
After=network.target

[Service]
Type=simple
User=your-user
WorkingDirectory=/path/to/discord-notifier
Environment="SUPABASE_URL=https://wbpfuuiznsmysbskywdx.supabase.co"
Environment="SUPABASE_ANON_KEY=your_key_here"
Environment="POLL_INTERVAL=10m"
ExecStart=/path/to/discord-notifier/discord-notifier
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

2. **Enable and start:**
```bash
sudo systemctl daemon-reload
sudo systemctl enable discord-notifier
sudo systemctl start discord-notifier
sudo systemctl status discord-notifier
```

### Option E: Docker

1. **Create Dockerfile:**
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o discord-notifier .

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/discord-notifier .
ENV TZ=UTC
CMD ["./discord-notifier"]
```

2. **Create docker-compose.yml:**
```yaml
version: '3.8'
services:
  discord-notifier:
    build: .
    container_name: discord-notifier
    restart: unless-stopped
    environment:
      - SUPABASE_URL=https://wbpfuuiznsmysbskywdx.supabase.co
      - SUPABASE_ANON_KEY=${SUPABASE_ANON_KEY}
      - POLL_INTERVAL=10m
    volumes:
      - ./logs:/root/logs
```

3. **Run:**
```bash
docker-compose up -d
docker-compose logs -f
```

## Step 6: Test the Integration

1. **Test Database:**
   - Manually activate a subscription for your test user
   - Add a Discord webhook URL
   - Add a test notification

2. **Test Website:**
   - Log in as test user
   - Click profile → Discord Notifications
   - Should see subscription status
   - Add a notification
   - Test webhook button

3. **Test Golang Service:**
   - Run the service
   - Check logs to see if it fetches users
   - Verify it processes notifications

## Step 7: Payment Integration (Choose One)

### Stripe Integration Example:

```javascript
// When user clicks "Subscribe" button
async function subscribeToDiscordNotifications() {
  // Redirect to Stripe Checkout
  const response = await fetch('/api/create-stripe-session', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      priceId: 'price_xxxxx', // Your Stripe $5/month price ID
      userId: currentUserId,
      successUrl: window.location.origin + '/subscription-success',
      cancelUrl: window.location.origin + '/subscription-cancel'
    })
  });
  
  const { sessionId } = await response.json();
  window.location.href = session.url;
}

// After successful payment (webhook handler or success page)
async function activateSubscription(userId) {
  const expiresAt = new Date();
  expiresAt.setMonth(expiresAt.getMonth() + 1);
  
  await supabase
    .from('unlocked_users')
    .update({
      discord_subscription_active: true,
      discord_subscription_expires_at: expiresAt.toISOString()
    })
    .eq('auth_user_id', userId);
}
```

## Step 8: Monitor and Maintain

1. **Set up logging** for your Golang service
2. **Monitor subscription expirations** - set up alerts
3. **Handle subscription renewals** - update expiry dates automatically
4. **Rate limiting** - Don't spam Discord webhooks
5. **Error handling** - Log failed webhook sends

## Quick Start Checklist

- [ ] Run SQL to add database columns
- [ ] Set up payment processor (Stripe/PayPal) OR use manual activation for testing
- [ ] Create Golang service project
- [ ] Add your Supabase credentials to .env
- [ ] Implement search logic in Golang service
- [ ] Test with one user
- [ ] Deploy Golang service (choose method above)
- [ ] Set up monitoring/logging
- [ ] Test end-to-end flow

## Troubleshooting

- **Service not fetching users:** Check Supabase API key and network
- **Webhooks not sending:** Verify webhook URLs are valid
- **Subscriptions not working:** Check database columns exist and subscription status
- **Service crashes:** Check logs, ensure proper error handling

## Next Steps

1. Implement your actual search logic in the Golang service
2. Add item tracking to avoid duplicate notifications
3. Set up automated subscription renewal handling
4. Add admin dashboard to manage subscriptions
5. Set up monitoring/alerting for service health
