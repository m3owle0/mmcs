# Optimizations for 40 Users

## Changes Made

### 1. **Rate Limiting Between Users**
- Added 2-second delay between processing users
- **Impact**: 40 users × 2 seconds = 80 seconds processing time
- **Benefit**: Reduces API load and prevents rate limiting

### 2. **Concurrent Request Limiting**
- Limited to 3 concurrent Sendico API requests at a time
- Added 500ms delay between market searches
- **Impact**: Prevents overwhelming Sendico API
- **Benefit**: More reliable, respects rate limits

### 3. **Translation Rate Limiting**
- Added 300ms delay before translation requests
- **Impact**: Prevents translation API rate limits
- **Benefit**: More reliable translations

### 4. **Increased Poll Interval**
- Changed from 5 minutes to 10 minutes
- **Impact**: Reduces API load by 50%
- **Benefit**: Less risk of hitting rate limits, still timely notifications

### 5. **Discord Webhook Delays**
- Increased delay between batches from 500ms to 1 second
- **Impact**: Better compliance with Discord's 30 requests/10 seconds limit
- **Benefit**: Prevents Discord rate limiting

## Performance Estimate

**With 40 users:**
- Average: 2 notifications per user = 80 notifications
- Each notification: 1 translation + 5 market searches = 6 API calls
- Total API calls: ~480 calls per poll cycle
- Processing time: ~80 seconds (users) + ~240 seconds (API calls) = **~5.3 minutes**
- **Poll interval**: 10 minutes ✅ (plenty of time)

## Capacity

✅ **Can handle 40 users comfortably** with these optimizations

**Scaling beyond 40 users:**
- Consider increasing poll interval to 15 minutes
- Consider database for item tracking (currently in-memory)
- Monitor Sendico API for 429 (rate limit) errors
- Consider processing users in batches (e.g., 20 at a time)

## Monitoring

Watch for these in logs:
- `⚠️ Sendico API error: status 429` - Rate limit hit, increase delays
- `⚠️ Discord returned status 429` - Discord rate limit, increase delays
- Processing time exceeding poll interval - increase poll interval

## Recommended Settings

**For 40 users:**
- ✅ Poll interval: 10 minutes (current)
- ✅ User delay: 2 seconds (current)
- ✅ Market delay: 500ms (current)
- ✅ Max concurrent: 3 (current)

**For 100+ users:**
- Poll interval: 15 minutes
- User delay: 3 seconds
- Market delay: 1 second
- Max concurrent: 2
