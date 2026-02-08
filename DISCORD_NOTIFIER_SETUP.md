# Discord Notifier Integration Guide

## Database Setup (Supabase SQL)

Run this SQL in your Supabase SQL Editor to add the required columns:

```sql
-- Add Discord notification columns to unlocked_users table
ALTER TABLE unlocked_users 
ADD COLUMN IF NOT EXISTS discord_webhook_url TEXT,
ADD COLUMN IF NOT EXISTS discord_notifications JSONB;

-- Add index for faster queries (optional but recommended)
CREATE INDEX IF NOT EXISTS idx_discord_notifications_active 
ON unlocked_users (discord_webhook_url) 
WHERE discord_webhook_url IS NOT NULL;
```

## API Endpoint for Golang Service

Your Golang service should query Supabase to fetch active notifications. Here's the API structure:

### Endpoint: Supabase REST API
**URL:** `https://wbpfuuiznsmysbskywdx.supabase.co/rest/v1/unlocked_users`

**Method:** `GET`

**Headers:**
```
apikey: YOUR_SUPABASE_ANON_KEY
Authorization: Bearer YOUR_SUPABASE_ANON_KEY
Content-Type: application/json
Prefer: return=representation
```

**Query Parameters:**
```
select=auth_user_id,discord_webhook_url,discord_notifications,email,username
discord_webhook_url=not.is.null
```

**Full Example Request:**
```bash
curl -X GET \
  'https://wbpfuuiznsmysbskywdx.supabase.co/rest/v1/unlocked_users?select=auth_user_id,discord_webhook_url,discord_notifications,email,username&discord_webhook_url=not.is.null' \
  -H 'apikey: YOUR_SUPABASE_ANON_KEY' \
  -H 'Authorization: Bearer YOUR_SUPABASE_ANON_KEY' \
  -H 'Content-Type: application/json'
```

### Response Format

The API will return an array of users with active Discord notifications:

```json
[
  {
    "auth_user_id": "user-uuid-here",
    "email": "user@example.com",
    "username": "username",
    "discord_webhook_url": "https://discord.com/api/webhooks/123456/abcdef",
    "discord_notifications": [
      {
        "id": "notification-id",
        "searchTerm": "vintage nike",
        "markets": ["mercari-jp", "grailed", "depop"],
        "createdAt": "2026-02-05T12:00:00Z"
      },
      {
        "id": "notification-id-2",
        "searchTerm": "supreme hoodie",
        "markets": null,
        "createdAt": "2026-02-05T13:00:00Z"
      }
    ]
  }
]
```

**Note:** If `markets` is `null`, the notification should monitor ALL markets.

## Golang Service Integration

### Recommended Service Structure

Your Golang service should:

1. **Poll Supabase periodically** (every 5-15 minutes recommended)
2. **Fetch all users with active notifications** using the API above
3. **For each user:**
   - Parse their `discord_notifications` JSON array
   - For each notification:
     - Check the specified markets (or all markets if `markets` is null)
     - Search for the `searchTerm`
     - If new items found, send Discord webhook notification
4. **Track which items have already been notified** (to avoid duplicates)

### Example Golang Code Structure

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type Notification struct {
    ID          string   `json:"id"`
    SearchTerm  string   `json:"searchTerm"`
    Markets     []string `json:"markets"` // null = all markets
    CreatedAt   string   `json:"createdAt"`
}

type User struct {
    AuthUserID          string         `json:"auth_user_id"`
    Email               string         `json:"email"`
    Username            string         `json:"username"`
    DiscordWebhookURL   string         `json:"discord_webhook_url"`
    DiscordNotifications []Notification `json:"discord_notifications"`
}

func fetchNotifications() ([]User, error) {
    // Your Supabase API call here
    // Use the endpoint structure above
}

func processNotifications(users []User) {
    for _, user := range users {
        for _, notif := range user.DiscordNotifications {
            // Process each notification
            // Check markets, search for term, send webhook if new items found
        }
    }
}

func main() {
    ticker := time.NewTicker(10 * time.Minute) // Check every 10 minutes
    defer ticker.Stop()
    
    for {
        users, err := fetchNotifications()
        if err != nil {
            fmt.Printf("Error fetching notifications: %v\n", err)
        } else {
            processNotifications(users)
        }
        <-ticker.C
    }
}
```

## How to Run the Golang Service

### Option 1: Direct Execution
```bash
# Navigate to your Golang project directory
cd /path/to/your/golang-notifier

# Run the service
go run main.go
```

### Option 2: Build and Run
```bash
# Build the binary
go build -o discord-notifier main.go

# Run the binary
./discord-notifier
```

### Option 3: Run as a Service (Linux)
```bash
# Create a systemd service file: /etc/systemd/system/discord-notifier.service
[Unit]
Description=MMCS Discord Notifier Service
After=network.target

[Service]
Type=simple
User=your-user
WorkingDirectory=/path/to/your/golang-notifier
ExecStart=/path/to/your/golang-notifier/discord-notifier
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target

# Enable and start the service
sudo systemctl enable discord-notifier
sudo systemctl start discord-notifier
sudo systemctl status discord-notifier
```

### Option 4: Run in Background (Windows)
```powershell
# Run in background using Start-Process
Start-Process -FilePath ".\discord-notifier.exe" -WindowStyle Hidden

# Or use Task Scheduler to run on startup
```

### Option 5: Docker (Recommended for Production)
```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o discord-notifier .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/discord-notifier .
CMD ["./discord-notifier"]
```

```bash
# Build and run
docker build -t discord-notifier .
docker run -d --name discord-notifier --restart unless-stopped discord-notifier
```

## Environment Variables (Recommended)

Create a `.env` file or set environment variables:

```bash
SUPABASE_URL=https://wbpfuuiznsmysbskywdx.supabase.co
SUPABASE_ANON_KEY=your_anon_key_here
POLL_INTERVAL=10m  # How often to check for new items
LOG_LEVEL=info
```

## Discord Webhook Format

When sending notifications to Discord, use this format:

```json
{
  "content": "🔔 **New Item Found!**",
  "embeds": [
    {
      "title": "Search Term: vintage nike",
      "description": "New item found on Mercari Japan",
      "url": "https://jp.mercari.com/item/123456",
      "color": 3447003,
      "fields": [
        {
          "name": "Price",
          "value": "¥5,000",
          "inline": true
        },
        {
          "name": "Market",
          "value": "Mercari Japan",
          "inline": true
        }
      ],
      "timestamp": "2026-02-05T12:00:00Z",
      "footer": {
        "text": "MMCS Notifications"
      }
    }
  ]
}
```

## Testing

1. **Test Webhook in UI:** Users can click "Test Webhook" button in the Discord Notifications modal
2. **Test Golang Service:** Run your service and check logs to ensure it's fetching notifications correctly
3. **Monitor:** Check your Discord channel for test notifications

## Troubleshooting

- **No notifications appearing:** Check that `discord_webhook_url` is not null in database
- **Webhook not working:** Verify webhook URL is correct and Discord server has webhook permissions
- **Service not fetching:** Check Supabase API key and network connectivity
- **Duplicate notifications:** Implement item tracking/hashing to avoid notifying same items twice

## Security Notes

- Store Supabase API key securely (environment variables, not in code)
- Consider using Row Level Security (RLS) policies in Supabase
- Rate limit your Discord webhook calls to avoid hitting Discord rate limits
- Monitor for abuse (users creating too many notifications)
