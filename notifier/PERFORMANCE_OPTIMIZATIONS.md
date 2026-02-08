# Performance Optimizations for 100+ Users

## Changes Made

### 1. **Parallel User Processing** ⚡
- **Before**: Sequential processing (one user at a time)
- **After**: Worker pool with 10 concurrent users
- **Impact**: 10x faster user processing
- **Benefit**: Can handle 100 users in ~2-3 minutes instead of 20+ minutes

### 2. **Translation Caching** 💾
- **Before**: Every search term translated individually
- **After**: Cache translations (many users search same terms like "nike", "rick owens")
- **Impact**: Eliminates duplicate translation API calls
- **Benefit**: If 20 users search "nike", only 1 API call instead of 20

### 3. **Optimized Rate Limiting** 🎯
- **Before**: 2-second delay between users, 500ms between markets
- **After**: Parallel processing with smart delays
  - 200ms between market searches (reduced from 500ms)
  - 500ms between Discord batches (reduced from 1s)
- **Impact**: Faster while still respecting API limits
- **Benefit**: Better throughput without hitting rate limits

### 4. **Increased Concurrency** 🚀
- **Market searches**: Up to 5 concurrent (increased from 3)
- **User processing**: 10 concurrent users
- **Impact**: Better CPU/network utilization
- **Benefit**: Processes more users simultaneously

### 5. **Reduced Poll Interval** ⏱️
- **Before**: 10 minutes
- **After**: 8 minutes
- **Impact**: More frequent checks
- **Benefit**: Faster notifications while still efficient

## Performance Estimates

### **100 Users Scenario**

**Assumptions:**
- 100 active subscribers
- Average 2 notifications per user
- Average 3 markets per notification
- 50% translation cache hit rate (common terms)

**Processing Time:**
- **Users**: 100 users ÷ 10 concurrent = 10 batches
- **Per batch**: ~15-20 seconds (includes API calls)
- **Total**: ~2.5-3 minutes per cycle
- **Poll interval**: 8 minutes ✅ (plenty of time)

**API Calls:**
- **Translations**: ~100 calls (50% cache hit saves 100 calls)
- **Market searches**: ~600 calls (100 users × 2 notifications × 3 markets)
- **Discord webhooks**: Variable (only when new items found)

**Memory Usage:**
- Translation cache: ~1-2 MB (1000 entries max)
- Seen items: ~5-10 MB (estimated 10,000 items)
- **Total**: ~10-15 MB (very manageable)

## Capacity

✅ **Can handle 100+ users comfortably** with these optimizations

**Scaling beyond 100 users:**
- **200 users**: Increase `maxConcurrentUsers` to 15-20
- **500+ users**: Consider database for item tracking
- **1000+ users**: May need multiple instances or queue system

## Configuration

Current settings (optimized for 100 users):
```go
maxConcurrentUsers = 10      // Process 10 users in parallel
maxConcurrentSearches = 5    // 5 concurrent market searches
pollInterval = 8 minutes     // Check every 8 minutes
```

**For 200+ users**, adjust:
```go
maxConcurrentUsers = 15-20
pollInterval = 10 minutes
```

## Monitoring

Watch for:
- **Rate limit errors (429)**: Increase delays if seen
- **Processing time**: Should complete in < poll interval
- **Memory usage**: Should stay under 50-100 MB
- **API response times**: Should be < 2 seconds per call

## Benefits

1. ⚡ **10x faster** user processing
2. 💾 **50-80% fewer** translation API calls (cache)
3. 🚀 **Better throughput** with parallel processing
4. ⏱️ **Faster notifications** (8-minute checks)
5. 📈 **Scalable** to 100+ users on single machine

## Next Steps (Optional)

For even better performance:
1. **Database item tracking**: Replace in-memory map
2. **Redis cache**: For translation cache (shared across instances)
3. **Message queue**: For very high scale (1000+ users)
4. **Multiple instances**: Horizontal scaling
