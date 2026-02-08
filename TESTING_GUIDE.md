# Testing Guide: Discord Notifications System

## Quick Test Setup (5 Minutes)

### Step 1: Set Up Database (One-Time)

1. **Go to Supabase Dashboard:**
   - https://supabase.com/dashboard
   - Select your project

2. **Open SQL Editor:**
   - Click "SQL Editor" in left sidebar
   - Click "New query"

3. **Run this SQL:**

```sql
-- Add columns if they don't exist
ALTER TABLE unlocked_users 
ADD COLUMN IF NOT EXISTS discord_webhook_url TEXT,
ADD COLUMN IF NOT EXISTS discord_notifications JSONB,
ADD COLUMN IF NOT EXISTS discord_subscription_active BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS discord_subscription_expires_at TIMESTAMPTZ;
```

4. **Verify:** Go to Table Editor → `unlocked_users` → Check columns exist

### Step 2: Get Your Discord Webhook URL

1. **Open Discord:**
   - Go to your Discord server
   - Right-click on the channel where you want notifications
   - Click "Edit Channel"

2. **Create Webhook:**
   - Go to "Integrations" → "Webhooks"
   - Click "New Webhook"
   - Name it (e.g., "MMCS Notifications")
   - Copy the "Webhook URL" (looks like: `https://discord.com/api/webhooks/123456/abcdef...`)
   - Click "Save Changes"

### Step 3: Activate Your Test Subscription

1. **Find Your User ID:**
   - Log into your website
   - Click your profile picture/username
   - Note your User ID (shown in profile modal)

2. **Activate Subscription in Supabase:**

Run this SQL (replace `YOUR_EMAIL` with your actual email):

```sql
-- Activate subscription for yourself
UPDATE unlocked_users 
SET 
  discord_subscription_active = TRUE,
  discord_subscription_expires_at = NOW() + INTERVAL '30 days'
WHERE email = 'YOUR_EMAIL';
```

Or if you know your `auth_user_id`:

```sql
UPDATE unlocked_users 
SET 
  discord_subscription_active = TRUE,
  discord_subscription_expires_at = NOW() + INTERVAL '30 days'
WHERE auth_user_id = 'YOUR_USER_ID';
```

### Step 4: Add Your Webhook and Notifications

1. **On Your Website:**
   - Click your profile → "🔔 Discord Notifications ($5/mo)"
   - You should see "✅ Active subscription"
   - Paste your Discord webhook URL
   - Click "+ Add Search Term"
   - Enter a test search term (e.g., "nike")
   - Optionally select specific markets (or leave empty for all)
   - Click "Add"
   - Click "Save"

2. **Test Webhook:**
   - Click "Test Webhook" button
   - Check your Discord channel - you should see a test message!

### Step 5: Set Up Golang Notifier (For Real Notifications)

#### Option A: Quick Test (Windows)

1. **Install Go:**
   - Download from: https://go.dev/dl/
   - Install it
   - Open PowerShell/Command Prompt
   - Verify: `go version`

2. **Create Project:**

```powershell
# Create folder
mkdir C:\discord-notifier
cd C:\discord-notifier

# Create main.go file (copy code from below)
notepad main.go
```

3. **Copy This Code to main.go:**

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
	SubscriptionExpiresAt *string        `json:"discord_subscription_expires_at"`
}

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

type DiscordWebhookPayload struct {
	Content string        `json:"content,omitempty"`
	Embeds  []DiscordEmbed `json:"embeds,omitempty"`
}

var (
	supabaseURL  = "https://wbpfuuiznsmysbskywdx.supabase.co"
	supabaseKey  = "" // You'll set this
	pollInterval = 1 * time.Minute // Check every minute for testing
)

func main() {
	// Get API key from environment or prompt
	if key := os.Getenv("SUPABASE_ANON_KEY"); key != "" {
		supabaseKey = key
	} else {
		fmt.Print("Enter your Supabase Anon Key: ")
		fmt.Scanln(&supabaseKey)
	}

	log.Printf("🚀 Starting Discord Notifier (TEST MODE)")
	log.Printf("📡 Supabase URL: %s", supabaseURL)
	log.Printf("⏱️  Poll interval: %v", pollInterval)
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Run immediately
	processAllNotifications()

	// Then run on interval
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
		log.Printf("❌ Error: %v", err)
		return
	}

	if len(users) == 0 {
		log.Printf("ℹ️  No active subscribers found")
		return
	}

	log.Printf("✅ Found %d subscriber(s)", len(users))

	for _, user := range users {
		if !isSubscriptionActive(user) {
			log.Printf("⏭️  Skipping %s - subscription expired", user.Email)
			continue
		}

		log.Printf("👤 Processing: %s (%s)", user.Username, user.Email)
		processUserNotifications(user)
	}
}

func fetchActiveSubscribers() ([]User, error) {
	// Use Supabase REST API
	url := fmt.Sprintf("%s/rest/v1/unlocked_users?select=auth_user_id,email,username,discord_webhook_url,discord_notifications,discord_subscription_active,discord_subscription_expires_at&discord_subscription_active=eq.true&discord_webhook_url=not.is.null", supabaseURL)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("apikey", supabaseKey)
	req.Header.Set("Authorization", "Bearer "+supabaseKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var users []User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
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
	if user.SubscriptionExpiresAt == nil || *user.SubscriptionExpiresAt == "" {
		return true
	}
	expiresAt, err := time.Parse(time.RFC3339, *user.SubscriptionExpiresAt)
	if err != nil {
		return false
	}
	return time.Now().Before(expiresAt)
}

func processUserNotifications(user User) {
	if len(user.DiscordNotifications) == 0 {
		log.Printf("   ℹ️  No notifications configured")
		return
	}

	for _, notif := range user.DiscordNotifications {
		log.Printf("   🔍 Checking: '%s'", notif.SearchTerm)

		// TEST MODE: Send a test notification immediately
		testItems := []map[string]interface{}{
			{
				"title":       fmt.Sprintf("Test Item for: %s", notif.SearchTerm),
				"description": "This is a test notification from MMCS",
				"url":         "https://example.com/test",
				"price":       "$50",
				"market":      "Test Market",
			},
		}

		log.Printf("   ✅ Sending test notification...")
		if err := sendDiscordNotification(user.DiscordWebhookURL, notif, testItems); err != nil {
			log.Printf("   ❌ Error: %v", err)
		} else {
			log.Printf("   ✅ Notification sent!")
		}
	}
}

func sendDiscordNotification(webhookURL string, notification Notification, items []map[string]interface{}) error {
	embeds := []DiscordEmbed{}
	for _, item := range items {
		embed := DiscordEmbed{
			Title:       fmt.Sprintf("New: %s", notification.SearchTerm),
			Description: getString(item, "description", "New item found!"),
			URL:         getString(item, "url", ""),
			Color:       3447003,
			Timestamp:   time.Now().Format(time.RFC3339),
			Footer: map[string]interface{}{
				"text": "MMCS Notifications",
			},
		}

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

	payload := DiscordWebhookPayload{
		Content: fmt.Sprintf("🔔 **%d new item(s) found for: %s**", len(items), notification.SearchTerm),
		Embeds:  embeds,
	}

	jsonData, _ := json.Marshal(payload)
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		return fmt.Errorf("status %d", resp.StatusCode)
	}

	return nil
}

func getString(m map[string]interface{}, key string, defaultValue string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return defaultValue
}
```

4. **Get Your Supabase Anon Key:**
   - Go to Supabase Dashboard → Select your project
   - Click "Project Settings" (gear icon) or "Connect" button
   - Go to "API Keys" section
   - Copy the "Publishable key" (sb_publishable_...) OR Legacy "anon" key

5. **Run the Service:**

```powershell
# Set environment variable (replace with your key)
$env:SUPABASE_ANON_KEY="your_anon_key_here"

# Run
go run main.go
```

**Expected Output:**
```
🚀 Starting Discord Notifier (TEST MODE)
📡 Supabase URL: https://wbpfuuiznsmysbskywdx.supabase.co
⏱️  Poll interval: 1m0s
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✅ Found 1 subscriber(s)
👤 Processing: yourusername (your@email.com)
   🔍 Checking: 'nike'
   ✅ Sending test notification...
   ✅ Notification sent!
```

6. **Check Discord:** You should see a notification in your Discord channel!

### Step 6: Test the Full Flow

1. **Add a Real Notification:**
   - Go to Discord Notifications modal
   - Add a search term you know exists (e.g., "vintage")
   - Save

2. **Modify the Golang Code:**
   - Replace the test notification code with your actual search logic
   - Or keep test mode to verify the system works

3. **Run Again:**
   ```powershell
   go run main.go
   ```

## Quick Test Checklist

- [ ] Database columns added
- [ ] Discord webhook created and copied
- [ ] Subscription activated for your account
- [ ] Webhook URL saved in website
- [ ] Test notification added
- [ ] "Test Webhook" button works (sends message to Discord)
- [ ] Golang service runs without errors
- [ ] Golang service sends notification to Discord

## Troubleshooting

**"No active subscribers found":**
- Check: `discord_subscription_active = TRUE` in database
- Check: `discord_webhook_url` is not null
- Verify: Your email matches in database

**"Webhook test fails":**
- Verify webhook URL is correct
- Check Discord server permissions
- Make sure webhook wasn't deleted

**"Golang service can't connect":**
- Check Supabase API key is correct
- Verify network connectivity
- Check Supabase project is active

**"Notifications not sending":**
- Check service logs for errors
- Verify webhook URL is valid
- Test webhook manually in Discord

## Next Steps

Once testing works:
1. Implement your actual search logic (replace test code)
2. Add item tracking to avoid duplicates
3. Deploy service to run 24/7 (see SETUP_INSTRUCTIONS.md)
4. Set up automatic subscription activation via Stripe webhook
