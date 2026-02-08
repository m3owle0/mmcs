# Implementing Real Market Search Notifications

## Current Status: TEST MODE ONLY ⚠️

The notifier currently **does NOT search markets**. It only sends test notifications.

## What You Need to Implement

To get **real notifications**, you need to implement search/scraping logic for each market.

### Option 1: Web Scraping (Most Common)

For each market, you need to:
1. **Make HTTP requests** to the market's search URL
2. **Parse HTML/JSON** to extract item listings
3. **Track new items** (compare against previously seen items)
4. **Send notifications** only for new items

**Example for Mercari Japan:**

```go
func searchMercariJP(searchTerm string) ([]map[string]interface{}, error) {
    // Build search URL (same as frontend)
    url := fmt.Sprintf("https://jp.mercari.com/en/search?keyword=%s&sort=created_time&order=desc", 
        url.QueryEscape(searchTerm))
    
    // Make HTTP request
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    // Parse HTML (use goquery or similar)
    doc, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil {
        return nil, err
    }
    
    // Extract items
    items := []map[string]interface{}{}
    doc.Find(".item-card").Each(func(i int, s *goquery.Selection) {
        title := s.Find(".item-title").Text()
        price := s.Find(".item-price").Text()
        itemURL, _ := s.Find("a").Attr("href")
        
        items = append(items, map[string]interface{}{
            "title": title,
            "price": price,
            "url":  "https://jp.mercari.com" + itemURL,
            "market": "Mercari Japan",
        })
    })
    
    return items, nil
}
```

### Option 2: Official APIs (If Available)

Some markets have APIs:
- **eBay**: Has an official API (requires API key)
- **Facebook Marketplace**: No public API (scraping only)
- **Mercari**: No public API (scraping only)
- Most others: Scraping required

### Option 3: Third-Party Services

Services like:
- **ScraperAPI** - Handles scraping for you
- **Bright Data** - Proxy/scraping service
- **Apify** - Scraping platform

## Challenges You'll Face

### 1. **Anti-Bot Protection**
- Many sites use Cloudflare, CAPTCHAs, rate limiting
- You'll need proxies, user agents, delays between requests
- May need headless browsers (Selenium, Playwright)

### 2. **HTML Structure Changes**
- Sites change their HTML frequently
- Your scrapers will break and need updates
- Requires maintenance

### 3. **Rate Limiting**
- Too many requests = IP ban
- Need delays between requests
- May need multiple IP addresses/proxies

### 4. **Item Tracking**
- Need to store previously seen items (database)
- Compare new results against old ones
- Only notify on truly new items

### 5. **Different Formats**
- Each market has different HTML structure
- Each needs custom parsing logic
- 33 markets = 33 different implementations

## Recommended Approach

### Phase 1: Start Small
1. **Pick 1-2 markets** to implement first (e.g., Mercari JP, eBay)
2. Implement basic scraping
3. Test thoroughly
4. Add item tracking (database)

### Phase 2: Expand Gradually
1. Add more markets one at a time
2. Handle edge cases and errors
3. Add retry logic and error handling

### Phase 3: Production Ready
1. Add proxy rotation
2. Add monitoring/alerting
3. Handle rate limits gracefully
4. Add duplicate detection

## Example Implementation Structure

```go
// Market searcher interface
type MarketSearcher interface {
    Search(term string) ([]Item, error)
    GetMarketName() string
}

// Item structure
type Item struct {
    Title       string
    Price       string
    URL         string
    Market      string
    ImageURL    string
    ListedAt    time.Time
}

// Implement for each market
type MercariJPSearcher struct {}
func (m *MercariJPSearcher) Search(term string) ([]Item, error) {
    // Implementation
}

type EBaySearcher struct {}
func (e *EBaySearcher) Search(term string) ([]Item, error) {
    // Implementation
}

// Main search function
func searchMarkets(searchTerm string, marketKeys []string) ([]Item, error) {
    searchers := map[string]MarketSearcher{
        "mercari-jp": &MercariJPSearcher{},
        "ebay": &EBaySearcher{},
        // ... etc
    }
    
    allItems := []Item{}
    for _, marketKey := range marketKeys {
        if searcher, ok := searchers[marketKey]; ok {
            items, err := searcher.Search(searchTerm)
            if err != nil {
                log.Printf("Error searching %s: %v", marketKey, err)
                continue
            }
            allItems = append(allItems, items...)
        }
    }
    
    return allItems, nil
}
```

## Database Schema for Item Tracking

You'll need a table to track seen items:

```sql
CREATE TABLE seen_items (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    notification_id TEXT NOT NULL,
    market_key TEXT NOT NULL,
    item_url TEXT NOT NULL,
    item_title TEXT,
    first_seen_at TIMESTAMPTZ DEFAULT NOW(),
    last_seen_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, notification_id, item_url)
);

CREATE INDEX idx_seen_items_user_notif ON seen_items(user_id, notification_id);
```

## Libraries You'll Need

```go
// HTML parsing
go get github.com/PuerkitoBio/goquery

// HTTP client with retries
go get github.com/cenkalti/backoff/v4

// Database (PostgreSQL)
go get github.com/lib/pq
// or
go get github.com/jackc/pgx/v5

// JSON handling (built-in)
// Time handling (built-in)
```

## Estimated Effort

- **Per market**: 2-8 hours (depending on complexity)
- **33 markets**: 66-264 hours of development
- **Maintenance**: Ongoing (sites change frequently)
- **Infrastructure**: Proxies, monitoring, error handling

## Alternative: Use Existing Services

Consider using services that already scrape these markets:
- **Price tracking APIs** (if they exist for your markets)
- **RSS feeds** (some markets have them)
- **Email alerts** (some markets send them - parse emails)

## Bottom Line

**Currently: NO, you will NOT receive real notifications.**

To get real notifications, you need to:
1. Implement scraping for each market (significant work)
2. Handle anti-bot protection
3. Track seen items in a database
4. Maintain scrapers as sites change

This is a **major development project** that will take weeks/months to implement properly for all 33 markets.
