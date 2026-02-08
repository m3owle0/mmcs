# Supported Markets for Discord Notifications

The notifier only processes notifications for markets that are supported. Custom markets (user-added websites) are **not supported** and will be automatically skipped.

## Supported Markets

### Asia Markets
- `mercari-jp` - Mercari Japan
- `paypay-fleamarket` - Yahoo PayPay Flea Market
- `rakuma` - Rakuten Rakuma
- `rakuten-jp` - Rakuten Japan
- `xianyu` - Xianyu (闲鱼)
- `yahoo-auctions` - Yahoo Japan Auctions
- `2ndstreet-jp` - 2nd Street JP
- `carousell-sg` - Carousell Singapore
- `carousell-hk` - Carousell Hong Kong
- `carousell-id` - Carousell Indonesia
- `carousell-my` - Carousell Malaysia
- `carousell-ph` - Carousell Philippines
- `carousell-tw` - Carousell Taiwan
- `fruits-family` - Fruits Family
- `kindal` - Kindal

### International Markets
- `depop` - Depop
- `ebay` - eBay
- `facebook` - Facebook Marketplace
- `gem` - Gem
- `grailed` - Grailed
- `mercari-us` - Mercari US
- `poshmark` - Poshmark
- `shopgoodwill` - ShopGoodwill
- `vinted` - Vinted
- `automated-searches` - Automated Searches
- `avito` - Avito
- `ebay-global` - eBay Global
- `google-images-past-month` - Google Images (Past month)
- `instagram` - Instagram

### Designer Markets
- `secondstreet` - 2nd STREET
- `therealreal` - The RealReal
- `vestiaire` - Vestiaire Collective

## Not Supported

- **Custom markets** (user-added websites starting with `custom-`) are **NOT supported**
- Any market not listed above will be automatically skipped
- If a notification has no supported markets after filtering, it will be skipped entirely

## How It Works

1. When processing a notification, the notifier checks each market in the notification's `markets` array
2. Only markets that exist in the `supportedMarkets` map are processed
3. Custom markets and unsupported markets are logged as warnings and skipped
4. If all markets are filtered out, the notification is skipped
5. If no markets are specified (empty array), all supported markets are used

## Adding Support for New Markets

To add support for a new market:

1. Add the market key to the `supportedMarkets` map in `main.go`
2. Ensure the market exists in the `marketUrls` object in `index.html`
3. Rebuild and restart the notifier service
