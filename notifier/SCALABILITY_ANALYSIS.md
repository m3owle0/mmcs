# Scalability Analysis: Can It Handle 40 Users?

## Current Architecture

- **Polling Interval**: 5 minutes
- **Processing**: Sequential (one user at a time)
- **Sendico API**: Concurrent requests per user (5 markets in parallel)
- **Item Tracking**: In-memory map
- **Discord Webhooks**: Sequential sending

## Capacity Analysis for 40 Users

### ✅ **What Works Well**

1. **Database Queries**: Supabase can easily handle 40 users
   - Single query fetches all active subscribers
   - No performance issues expected

2. **Memory Usage**: In-memory tracking is fine for 40 users
   - Each user might have 1-5 notifications
   - Each notification tracks ~20-50 items
   - Estimated: ~4,000-10,000 items in memory
   - **Memory**: ~1-5 MB (very manageable)

3. **Concurrent Processing**: Uses errgroup for parallel market searches
   - Each user's markets searched in parallel
   - Good for performance

### ⚠️ **Potential Issues**

1. **Sendico API Rate Limiting** ⚠️
   - **Current**: No delays between requests
   - **Risk**: If 40 users × 2 notifications × 5 markets = 400 API calls every 5 minutes
   - **Sendico Limits**: Unknown, but likely has rate limits
   - **Solution Needed**: Add delays between requests

2. **Sequential User Processing** ⚠️
   - **Current**: Processes users one at a time
   - **Impact**: If each user takes 2-5 seconds, 40 users = 80-200 seconds
   - **Risk**: Could exceed 5-minute poll interval
   - **Solution**: Add delays or process in batches

3. **Discord Webhook Rate Limits** ⚠️
   - **Limit**: 30 requests per 10 seconds per webhook
   - **Current**: Sequential sending with 500ms delay between batches
   - **Risk**: If many users get notifications simultaneously, could hit limits
   - **Solution**: Add exponential backoff

4. **HMAC Refresh** ⚠️
   - Refreshes every 30 minutes
   - Should be fine, but adds overhead

## Recommendations for 40 Users

### **Option 1: Add Rate Limiting (Recommended)**

Add delays to respect API limits:

```go
// In processUserNotifications:
time.Sleep(2 * time.Second) // Delay between users

// In BulkSearch:
time.Sleep(500 * time.Millisecond) // Delay between market searches
```

### **Option 2: Increase Poll Interval**

Change from 5 minutes to 10-15 minutes:
- Reduces API load by 50-66%
- Still provides timely notifications
- Less risk of rate limiting

### **Option 3: Batch Processing**

Process users in smaller batches:
- Process 10 users, wait 1 minute, process next 10
- Spreads load over time
- Reduces peak API usage

### **Option 4: Database Item Tracking**

Replace in-memory tracking with database:
- Better for scaling beyond 40 users
- Prevents memory issues
- Allows persistence across restarts

## Estimated Performance

**Best Case (with optimizations):**
- 40 users × 2 notifications × 5 markets = 400 API calls
- With 500ms delays: ~200 seconds (3.3 minutes)
- **Verdict**: ✅ **Works fine** with 5-minute polling

**Worst Case (current code):**
- No delays, all concurrent
- Could hit rate limits
- **Verdict**: ⚠️ **May have issues** - needs optimization

## Conclusion

**Can it handle 40 users?** 

**YES, but with modifications:**

1. ✅ Add delays between API requests (2 seconds between users, 500ms between markets)
2. ✅ Consider increasing poll interval to 10 minutes
3. ✅ Add error handling for rate limits
4. ✅ Monitor Sendico API responses for 429 (Too Many Requests) errors

**Without modifications:** ⚠️ **Risky** - may hit rate limits

## Quick Fixes Needed

1. Add rate limiting delays
2. Add retry logic for rate limit errors
3. Add logging for API response times
4. Consider database for item tracking (optional but recommended)
