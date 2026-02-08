# Backend Setup - Discord Notifier Service

## Step 1: Install Go

**Windows:**
1. Download: https://go.dev/dl/
2. Run installer
3. Open PowerShell, verify: `go version`

**Mac:**
```bash
brew install go
```

**Linux:**
```bash
sudo apt-get install golang-go
```

## Step 2: Create Project Folder

```powershell
# Windows PowerShell
mkdir C:\discord-notifier
cd C:\discord-notifier
```

```bash
# Mac/Linux
mkdir ~/discord-notifier
cd ~/discord-notifier
```

## Step 3: Get Your Supabase API Key

1. Go to: https://supabase.com/dashboard
2. Select your project
3. Click **"Project Settings"** (gear icon) in the left sidebar, OR click the **"Connect"** button
4. Go to **"API Keys"** section
5. Copy the **"Publishable key"** (starts with `sb_publishable_...`) - this is the new format
   - OR if you see a **"Legacy API Keys"** tab, copy the **"anon"** key from there
6. **Important:** Use the publishable/anon key, NOT the service_role key (which has admin privileges)

## Step 4: Create main.go

Create a file called `main.go` in your project folder and paste this code:

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
	supabaseKey  = ""
	pollInterval = 5 * time.Minute // Check every 5 minutes
)

func main() {
	// Get API key from environment or prompt
	if key := os.Getenv("SUPABASE_ANON_KEY"); key != "" {
		supabaseKey = key
	} else {
		fmt.Print("Enter your Supabase Anon Key: ")
		fmt.Scanln(&supabaseKey)
		if supabaseKey == "" {
			log.Fatal("❌ API key is required")
		}
	}

	log.Printf("🚀 Starting Discord Notifier")
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
		log.Printf("❌ Error fetching users: %v", err)
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
	url := fmt.Sprintf("%s/rest/v1/unlocked_users?select=auth_user_id,email,username,discord_webhook_url,discord_notifications,discord_subscription_active,discord_subscription_expires_at&discord_subscription_active=eq.true&discord_webhook_url=not.is.null", supabaseURL)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("apikey", supabaseKey)
	req.Header.Set("Authorization", "Bearer "+supabaseKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
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
		return true // Lifetime subscription
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

		// TODO: Implement your actual search logic here
		// For now, this is a placeholder that sends a test notification
		
		// Example: Search markets and find new items
		// items := searchMarkets(notif.SearchTerm, notif.Markets)
		// if len(items) > 0 {
		//     sendDiscordNotification(user.DiscordWebhookURL, notif, items)
		// }

		// TEST MODE: Send test notification
		testItems := []map[string]interface{}{
			{
				"title":       fmt.Sprintf("Test Item: %s", notif.SearchTerm),
				"description": "This is a test notification. Replace this with your actual search results.",
				"url":         "https://example.com/test",
				"price":       "$50",
				"market":      "Test Market",
			},
		}

		log.Printf("   ✅ Sending notification...")
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
			Title:       fmt.Sprintf("🔔 New: %s", notification.SearchTerm),
			Description: getString(item, "description", "New item found!"),
			URL:         getString(item, "url", ""),
			Color:       3447003, // Blue color
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
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		return fmt.Errorf("Discord returned status %d", resp.StatusCode)
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

## Step 5: Run It

### Option A: Run Once (Testing)

**Windows PowerShell:**
```powershell
$env:SUPABASE_ANON_KEY="paste_your_key_here"
go run main.go
```

**Mac/Linux:**
```bash
export SUPABASE_ANON_KEY="paste_your_key_here"
go run main.go
```

**Or enter key when prompted:**
```powershell
go run main.go
# Then paste your key when asked
```

### Option B: Build Executable

**Windows:**
```powershell
go build -o discord-notifier.exe main.go
$env:SUPABASE_ANON_KEY="your_key"
.\discord-notifier.exe
```

**Mac/Linux:**
```bash
go build -o discord-notifier main.go
export SUPABASE_ANON_KEY="your_key"
./discord-notifier
```

## Step 6: Keep It Running 24/7

### Windows: Run as Service (NSSM)

1. **Download NSSM:** https://nssm.cc/download
2. **Extract and run:**
```powershell
# In NSSM folder
.\nssm.exe install DiscordNotifier
# Set these:
# Path: C:\discord-notifier\discord-notifier.exe
# Startup directory: C:\discord-notifier
# Environment: SUPABASE_ANON_KEY=your_key_here
.\nssm.exe start DiscordNotifier
```

### Windows: Task Scheduler

1. Open Task Scheduler
2. Create Basic Task
3. Name: "Discord Notifier"
4. Trigger: When computer starts
5. Action: Start a program
6. Program: `C:\discord-notifier\discord-notifier.exe`
7. Start in: `C:\discord-notifier`
8. Add argument: (leave empty)
9. Check "Run whether user is logged on or not"

**Set environment variable:**
- In Task Scheduler → Properties → General → "Run with highest privileges"
- Actions → Edit → Add to "Start in (optional)": Set environment variable in a batch file instead

**Or create `start-notifier.bat`:**
```batch
@echo off
set SUPABASE_ANON_KEY=your_key_here
cd C:\discord-notifier
discord-notifier.exe
```

Then point Task Scheduler to `start-notifier.bat`

### Mac/Linux: systemd Service

Create `/etc/systemd/system/discord-notifier.service`:

```ini
[Unit]
Description=Discord Notifier Service
After=network.target

[Service]
Type=simple
User=your_username
WorkingDirectory=/home/your_username/discord-notifier
Environment="SUPABASE_ANON_KEY=your_key_here"
ExecStart=/home/your_username/discord-notifier/discord-notifier
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Then:
```bash
sudo systemctl daemon-reload
sudo systemctl enable discord-notifier
sudo systemctl start discord-notifier
sudo systemctl status discord-notifier
```

### Mac/Linux: Screen/Tmux (Simple)

```bash
# Install screen
sudo apt-get install screen  # Linux
brew install screen          # Mac

# Run in screen
screen -S notifier
export SUPABASE_ANON_KEY="your_key"
./discord-notifier

# Detach: Ctrl+A then D
# Reattach: screen -r notifier
```

## Step 7: Test It

1. **Make sure you have:**
   - Active subscription in database
   - Discord webhook URL saved
   - At least one notification configured

2. **Run the service:**
   ```powershell
   go run main.go
   ```

3. **Expected output:**
   ```
   🚀 Starting Discord Notifier
   📡 Supabase URL: https://wbpfuuiznsmysbskywdx.supabase.co
   ⏱️  Poll interval: 5m0s
   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
   ✅ Found 1 subscriber(s)
   👤 Processing: yourname (your@email.com)
      🔍 Checking: 'nike'
      ✅ Sending notification...
      ✅ Notification sent!
   ```

4. **Check Discord:** You should see a notification!

## Next Steps: Add Real Search Logic

Replace the test code in `processUserNotifications()` with your actual market search implementation. The function receives:
- `user`: User object with webhook URL
- `notif`: Notification object with search term and markets
- You need to: Search markets → Find new items → Send to Discord

## Troubleshooting

**"API returned status 401":**
- Check your Supabase API Key is correct
- Make sure you copied the "Publishable key" (sb_publishable_...) or Legacy "anon" key, NOT the service_role key
- Verify the key hasn't expired or been regenerated

**"No active subscribers found":**
- Check database: `discord_subscription_active = TRUE`
- Check: `discord_webhook_url` is not null
- Verify user has notifications configured

**"Discord returned status 404":**
- Webhook URL is invalid or deleted
- Check webhook in Discord channel settings

**Service stops running:**
- Check logs for errors
- Make sure environment variable is set correctly
- Verify Go executable has proper permissions

## Quick Reference

**Supabase URL:** `https://wbpfuuiznsmysbskywdx.supabase.co`

**Get API Key:** 
1. Supabase Dashboard → Select Project
2. Project Settings (gear) or Connect button → API Keys
3. Copy "Publishable key" (sb_publishable_...) OR Legacy "anon" key

**Test Run:**
```powershell
$env:SUPABASE_ANON_KEY="your_key"
go run main.go
```

**Build:**
```powershell
go build -o discord-notifier.exe main.go
```

That's it! Your backend is ready to run. 🚀
