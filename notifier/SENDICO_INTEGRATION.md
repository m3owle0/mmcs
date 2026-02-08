# Sendico Integration - 5 Markets Working! ✅

## What Was Implemented

The notifier now has **real search functionality** for **5 Japanese markets** via the Sendico API:

1. **mercari-jp** - Mercari Japan
2. **paypay-fleamarket** - Yahoo PayPay Flea Market  
3. **rakuma** - Rakuten Rakuma
4. **rakuten-jp** - Rakuten
5. **yahoo-auctions** - Yahoo Japan Auctions

## How It Works

1. **Translation**: Search terms are automatically translated from English to Japanese (Sendico requires Japanese)
2. **Search**: Uses Sendico API to search across the selected markets
3. **Deduplication**: Tracks seen items in memory to avoid duplicate notifications
4. **Notifications**: Sends Discord webhooks with real item data (title, price, image, URL)

## Setup

### 1. Install Dependencies

```powershell
cd notifier
go mod download
```

### 2. Run

```powershell
$env:SUPABASE_ANON_KEY="your_key"
go run .
```

Or use `start.bat` (edit it first to add your key).

## Features

- ✅ **Real searches** - Actually queries Sendico API
- ✅ **Auto-translation** - English → Japanese automatically
- ✅ **Item tracking** - Prevents duplicate notifications
- ✅ **Rich Discord embeds** - Includes images, prices, market names
- ✅ **Error handling** - Continues if one market fails
- ✅ **Logging** - Detailed logs for debugging

## Limitations

- **In-memory tracking**: Seen items are stored in memory (lost on restart)
  - For production, consider using a database (SQLite/PostgreSQL)
- **5 markets only**: Only the 5 Sendico-supported markets work
  - Other markets will be skipped with a warning
- **Sendico dependency**: Relies on Sendico API working
  - If Sendico changes their API, this may break

## Example Output

```
🚀 Starting Discord Notifier
📡 Supabase URL: https://wbpfuuiznsmysbskywdx.supabase.co
⏱️  Poll interval: 5m0s
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🔧 Initializing Sendico client...
✅ Sendico client initialized
✅ Found 1 subscriber(s)
👤 Processing: username (user@email.com)
   🔍 Checking: 'rick owens'
   📋 Markets: [mercari-jp rakuma] (filtered to supported: [mercari-jp rakuma])
   🇯🇵 Translated 'rick owens' → 'リックオーウェンズ'
   🔎 Searching 2 market(s)...
   📦 Found 15 item(s)
   ✨ 3 new item(s) found!
   ✅ Sending notification...
   ✅ Notification sent!
```

## Next Steps (Optional)

1. **Database tracking**: Replace in-memory `seenItems` with database
2. **More markets**: Implement scraping for other markets
3. **Price filtering**: Add min/max price support
4. **Image caching**: Cache item images locally
5. **Rate limiting**: Add delays between searches

## Files Added

- `sendico.go` - Sendico API client
- `hmac.go` - HMAC signature generation for Sendico
- `go.mod` - Go dependencies
- Updated `main.go` - Integrated Sendico search

## Testing

1. Set up a notification in your website for one of the 5 markets
2. Run the notifier
3. Check Discord for notifications with real items!
