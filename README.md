# Multi-Market Secondhand Clothing Search (MMCS)

A web application for searching secondhand clothing across multiple marketplaces.

## Features

- 🔍 Search across multiple secondhand clothing marketplaces
- 💾 Save searches and organize them in folders
- 👤 User authentication with Supabase
- 💳 Discord notifications subscription ($10/month)
- 📱 Mobile-responsive design
- 🔔 Real-time Discord notifications for saved searches

## Setup

### Prerequisites

- A Supabase account and project
- Stripe account (for subscriptions)

### Installation

1. Clone this repository
2. Open `index.html` in a browser or deploy to a static hosting service

### Database Setup

The `unlocked_users` table in Supabase should have the following columns:
- `auth_user_id` (text)
- `email` (text)
- `username` (text)
- `notifications_subscription_active` (boolean)
- `notifications_subscription_expires_at` (timestamp)
- `payment_method` (text)
- `stripe_customer_id` (text, nullable)
- `stripe_subscription_id` (text, nullable)
- `discord_webhook_url` (text, nullable)
- `discord_notifications` (jsonb, nullable)
- `last_active` (timestamp)
- `description` (text, nullable)
- `profile_picture_url` (text, nullable)
- `created_at` (timestamp)

### Configuration

1. **Supabase Configuration**
   - Update the Supabase URL and anon key in `index.html` (search for `SUPABASE_URL` and `SUPABASE_ANON_KEY`)

2. **Stripe Configuration**
   - Update Stripe checkout link in `index.html` (search for `NOTIFICATIONS_STRIPE_LINK`)
   - Replace with your own Stripe payment link for Discord notifications ($10/month)

## Project Structure

```
mmcs/
├── index.html                    # Main application file
├── miku.png                      # Favicon
├── mmcs.png                      # Logo (optional)
├── bg.png                        # Background image (optional)
├── notifier/                     # Discord notifier service (optional)
│   ├── main.go
│   ├── sendico.go
│   ├── hmac.go
│   └── go.mod
├── supabase/                     # Supabase functions (optional)
│   ├── config.toml
│   └── functions/
│       └── stripe-webhook/
│           └── index.ts
└── README.md                     # This file
```

### Optional Components

- **supabase/functions/**: Stripe webhook handler for automatic subscription activation
- **notifier/**: Go-based Discord notification service for subscribed users

## Usage

1. Open `index.html` in a web browser
2. Sign up or log in with your email
3. Start searching for clothing items
4. Save searches and organize them in folders
5. Subscribe to Discord notifications ($10/month) for real-time alerts on your saved searches

## Pricing

- **Free**: Unlimited searches, all features, forever - completely free!
- **Discord Notifications ($10/month)**: Real-time alerts when new items match your saved searches

## License

MIT License
