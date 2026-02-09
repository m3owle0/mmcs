# Multi-Market Clothing Search (MMCS)

A web application for searching secondhand clothing across multiple marketplaces with Discord notification support.

## Features

- **Multi-Market Search**: Search across 20+ marketplaces including Mercari, Rakuten, eBay, Grailed, and more
- **Discord Notifications**: Get notified when new items matching your searches are found
- **User Authentication**: Supabase-powered authentication with verification system
- **Custom Markets**: Add your own custom market URLs
- **Wishlist**: Save favorite searches and URLs
- **Theme Customization**: Customize colors and background images

## Setup Instructions

### 1. Supabase Database Setup

1. Create a new Supabase project at [supabase.com](https://supabase.com)
2. Go to **SQL Editor** in your Supabase dashboard
3. Run the SQL schema from `database/schema.sql`:
   - This creates the `unlocked_users` table
   - Sets up Row Level Security (RLS) policies
   - Creates triggers for automatic user record creation

### 2. Configure Supabase Credentials

1. In Supabase Dashboard, go to **Project Settings** → **API**
2. Copy your **Project URL** and **anon/public key**
3. Update `index.html` with your credentials:
   ```javascript
   const SUPABASE_URL = 'https://your-project.supabase.co';
   const SUPABASE_ANON_KEY = 'your-anon-key-here';
   ```
   (Update both occurrences in the file)

### 3. Discord Notifier Setup

The notifier is a Go service that polls the database and sends Discord notifications.

#### Prerequisites
- Go 1.23+ installed
- Supabase Service Role Key (for bypassing RLS)

#### Configuration

1. Get your **Service Role Key** from Supabase:
   - Go to **Project Settings** → **API**
   - Copy the **service_role** key (keep this secret!)

2. Set up the notifier:
   ```bash
   cd notifier
   go mod download
   ```

3. Run the notifier:
   ```bash
   # Using environment variable (recommended)
   export SUPABASE_SERVICE_ROLE_KEY="your-service-role-key"
   go run .

   # Or it will prompt you for the key
   go run .
   ```

4. For production, build and run:
   ```bash
   go build -o notifier
   ./notifier
   ```

#### Running as a Service (Linux)

Create `/etc/systemd/system/mmcs-notifier.service`:
```ini
[Unit]
Description=MMCS Discord Notifier
After=network.target

[Service]
Type=simple
User=your-user
WorkingDirectory=/path/to/mmcs/notifier
Environment="SUPABASE_SERVICE_ROLE_KEY=your-service-role-key"
ExecStart=/path/to/mmcs/notifier/notifier
Restart=always

[Install]
WantedBy=multi-user.target
```

Then:
```bash
sudo systemctl enable mmcs-notifier
sudo systemctl start mmcs-notifier
```

### 4. User Setup

#### For Users:

1. **Sign Up**: Create an account on the website
2. **Get Verified**: Contact the admin to verify your account (sets `verified = true` in database)
3. **Configure Discord Webhook**:
   - Go to your Discord server
   - Server Settings → Integrations → Webhooks → New Webhook
   - Copy the webhook URL
   - Click "🔔 Notifications" button on the website
   - Paste your webhook URL and save
4. **Create Notifications**:
   - Click "🔔 Notifications" → "+ Add Notification"
   - Enter search term (e.g., "Supreme hoodie")
   - Select markets (or leave empty for all)
   - Save

#### For Admins (Verifying Users):

Run this SQL in Supabase SQL Editor:
```sql
-- Verify a user by email
UPDATE unlocked_users 
SET verified = true 
WHERE email = 'user@example.com';

-- Activate notifications subscription (lifetime)
UPDATE unlocked_users 
SET notifications_subscription_active = true,
    notifications_subscription_expires_at = NULL
WHERE email = 'user@example.com';

-- Activate notifications subscription (with expiration)
UPDATE unlocked_users 
SET notifications_subscription_active = true,
    notifications_subscription_expires_at = '2025-12-31T23:59:59Z'
WHERE email = 'user@example.com';
```

### 5. Deploy Website

#### Option 1: Static Hosting (Netlify, Vercel, GitHub Pages)

1. Upload `index.html` and image files (`miku.svg`, `mmcs.svg`, `bg.svg`)
2. Configure your hosting platform
3. The website will work as-is (no backend needed for the frontend)

#### Option 2: Local Development

1. Open `index.html` in a web browser
2. Or use a local server:
   ```bash
   # Python
   python -m http.server 8000
   
   # Node.js
   npx http-server
   ```

## File Structure

```
.
├── index.html                          # Main website file
├── multi-market-clothing-search.html   # Alternative version
├── database/
│   └── schema.sql                      # Database schema
├── notifier/
│   ├── main.go                         # Notifier service
│   ├── sendico.go                      # Sendico API client
│   ├── hmac.go                         # HMAC signing
│   └── go.mod                          # Go dependencies
├── miku.svg                            # Favicon
├── mmcs.svg                            # Logo
├── bg.svg                              # Background image
└── README.md                           # This file
```

## Supported Markets

### Japanese Markets (via Sendico API)
- Mercari Japan
- Yahoo PayPay Flea Market
- Rakuten Rakuma
- Rakuten
- Yahoo Auctions

### Other Markets (direct links)
- Xianyu (China)
- Depop
- eBay
- Facebook Marketplace
- Gem
- Grailed
- Mercari US
- Poshmark
- ShopGoodwill
- Vinted
- 2nd Street
- The RealReal
- Vestiaire Collective
- And more...

## Notifications

The notifier service:
- Polls the database every 1 minute
- Searches Sendico API for Japanese markets
- Translates search terms to Japanese automatically
- Sends Discord webhook notifications for new items
- Tracks seen items to avoid duplicates
- Supports up to 10 concurrent users and 5 concurrent searches

## Security Notes

- **Service Role Key**: Keep this secret! Never commit it to version control
- **Discord Webhooks**: Users should keep their webhook URLs private
- **RLS Policies**: The database uses Row Level Security - users can only access their own data
- **Notifier**: Uses service role key to bypass RLS (required for reading all users)

## Troubleshooting

### Notifications not working?
1. Check that the notifier service is running
2. Verify user has `notifications_subscription_active = true`
3. Check that Discord webhook URL is set
4. Check notifier logs for errors

### User can't log in?
1. Verify user exists in `unlocked_users` table
2. Check that `verified = true` for the user
3. Verify Supabase credentials in `index.html`

### Database errors?
1. Ensure schema.sql has been run
2. Check RLS policies are enabled
3. Verify service role key is correct (for notifier)

## License

[Add your license here]

## Support

For issues or questions, please contact [your contact info]
