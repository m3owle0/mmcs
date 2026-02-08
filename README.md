# Multi-Market Secondhand Clothing Search (MMCS)

A web application for searching secondhand clothing across multiple marketplaces.

## Features

- 🔍 Search across multiple secondhand clothing marketplaces
- 💾 Save searches and organize them in folders
- 👤 User authentication with Supabase
- 💳 Subscription tiers (Trial, Basic, Pro)
- 📱 Mobile-responsive design
- 🔔 Discord notifications (Pro tier)

## Setup

### Prerequisites

- A Supabase account and project
- Stripe account (for subscriptions)

### Installation

1. Clone this repository
2. Open `index.html` in a browser or deploy to a static hosting service

### Database Setup

Run the SQL in `ADD_SUBSCRIPTION_COLUMNS.sql` in your Supabase SQL Editor to add subscription columns to the `unlocked_users` table.

### Configuration

1. **Supabase Configuration**
   - Update the Supabase URL and anon key in `index.html` (search for `SUPABASE_URL` and `SUPABASE_ANON_KEY`)

2. **Stripe Configuration**
   - Update Stripe checkout links in `index.html` (search for `STRIPE_CHECKOUT_LINKS`)
   - Replace with your own Stripe payment links

## Project Structure

```
mmcs/
├── index.html                    # Main application file
├── miku.png                      # Favicon
├── mmcs.png                      # Logo
├── bg.png                        # Background image
├── ADD_SUBSCRIPTION_COLUMNS.sql  # Database setup SQL
├── notifier/                     # Discord notifier service (optional)
├── supabase/                     # Supabase functions (optional)
└── README.md                     # This file
```

### Optional Components

- **supabase/functions/**: Stripe webhook handler for automatic subscription activation
- **notifier/**: Go-based Discord notification service for Pro tier users

## Usage

1. Open `index.html` in a web browser
2. Sign up or log in with your email
3. Start searching for clothing items
4. Save searches and organize them in folders
5. Upgrade to Basic or Pro tier for unlimited searches and premium features

## Subscription Tiers

- **Free Trial**: 75 free searches
- **Basic ($5/month)**: Unlimited searches, premium badge, priority support
- **Pro ($10/month)**: Everything in Basic + Discord notifications, advanced filters

## License

MIT License
