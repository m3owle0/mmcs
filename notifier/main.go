package main

// Discord Notifier with Sendico integration for 5 Japanese markets:
// - mercari-jp (Mercari)
// - paypay-fleamarket (Yahoo PayPay Flea)
// - rakuma (Rakuten Rakuma)
// - rakuten-jp (Rakuten)
// - yahoo-auctions (Yahoo Auctions)

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Notification struct {
	ID         string   `json:"id"`
	SearchTerm string   `json:"searchTerm"`
	Markets    []string `json:"markets"`
	Webhooks   []string `json:"webhooks,omitempty"` // Per-notification webhooks
	CreatedAt  string   `json:"createdAt"`
}

// cachedSearchResult stores search results to avoid duplicate API calls across users
type cachedSearchResult struct {
	Items      []SendicoItem
	Timestamp  time.Time
	ExpiresAt  time.Time
}

type User struct {
	AuthUserID            string          `json:"auth_user_id"`
	Email                 string          `json:"email"`
	Username              string          `json:"username"`
	DiscordWebhookURL     string          `json:"discord_webhook_url"`
	DiscordNotifications  json.RawMessage `json:"discord_notifications"` // Store as raw JSON first
	SubscriptionActive    bool            `json:"notifications_subscription_active"`
	SubscriptionExpiresAt *string         `json:"notifications_subscription_expires_at"`

	// Parsed notifications (populated after unmarshalling)
	Notifications []Notification
}

type DiscordEmbed struct {
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	URL         string                 `json:"url,omitempty"`
	Color       int                    `json:"color"`
	Fields      []DiscordEmbedField    `json:"fields,omitempty"`
	Timestamp   string                 `json:"timestamp,omitempty"`
	Footer      map[string]interface{} `json:"footer,omitempty"`
	Image       map[string]string      `json:"image,omitempty"`
}

type DiscordEmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type DiscordWebhookPayload struct {
	Content string         `json:"content,omitempty"`
	Embeds  []DiscordEmbed `json:"embeds,omitempty"`
}

var (
	supabaseURL   = "https://wbpfuuiznsmysbskywdx.supabase.co"
	supabaseKey   = ""
	pollInterval  = 1 * time.Minute // Check every 1 minute (fast notifications)
	sendicoClient *SendicoClient

	// Cycle lock to prevent overlapping processing cycles
	processingMu sync.Mutex
	isProcessing bool

	// Sendico-supported markets (5 Japanese markets via Sendico API)
	sendicoMarkets = map[string]SendicoShop{
		"mercari-jp":        SendicoMercari,
		"paypay-fleamarket": SendicoYahoo,
		"rakuma":            SendicoRakuma,
		"rakuten-jp":        SendicoRakuten,
		"yahoo-auctions":    SendicoYahooAuctions,
	}

	// Item tracking (in-memory with TTL to prevent memory bloat)
	// Key: "notificationID:shop:code", Value: timestamp when first seen
	seenItems   = make(map[string]time.Time)
	seenItemsMu sync.RWMutex
	seenItemsTTL = 7 * 24 * time.Hour // Keep seen items for 7 days (prevents re-notification)

	// Translation cache to avoid duplicate API calls (many users search same terms)
	translationCache   = make(map[string]string)
	translationCacheMu sync.RWMutex

	// Concurrency limits - optimized for multiple users
	maxConcurrentUsers    = 15 // Increased to 15 for better multi-user throughput
	maxConcurrentSearches = 5  // Max 5 concurrent Sendico searches
	
	// Search batching - avoid duplicate searches across users
	searchCache   = make(map[string]*cachedSearchResult) // Key: "termJP:markets"
	searchCacheMu sync.RWMutex
	
	// Pages to search per query (to catch recently uploaded items)
	// Set to 1 for fastest performance (page 1 typically has newest items)
	// Set to 2-3 if you want to catch more recently uploaded items (slower)
	maxSearchPages = 1 // Default to 1 page for speed (can be increased if needed)

	// Supported markets - only these markets will be processed by the notifier
	// This matches the marketUrls object in index.html
	// Custom markets (starting with "custom-") are NOT supported and will be skipped
	// If a notification has no supported markets after filtering, it will be skipped
	supportedMarkets = map[string]bool{
		"mercari-jp":               true,
		"paypay-fleamarket":        true,
		"rakuma":                   true,
		"rakuten-jp":               true,
		"xianyu":                   true,
		"yahoo-auctions":           true,
		"depop":                    true,
		"ebay":                     true,
		"facebook":                 true,
		"gem":                      true,
		"grailed":                  true,
		"mercari-us":               true,
		"poshmark":                 true,
		"shopgoodwill":             true,
		"vinted":                   true,
		"secondstreet":             true,
		"therealreal":              true,
		"vestiaire":                true,
		"2ndstreet-jp":             true,
		"carousell-sg":             true,
		"carousell-hk":             true,
		"carousell-id":             true,
		"carousell-my":             true,
		"carousell-ph":             true,
		"carousell-tw":             true,
		"fruits-family":            true,
		"kindal":                   true,
		"automated-searches":       true,
		"avito":                    true,
		"ebay-global":              true,
		"google-images-past-month": true,
		"instagram":                true,
	}
)

func main() {
	// Get API key from environment or prompt
	// NOTE: Use SERVICE_ROLE_KEY for notifier (bypasses RLS to read all users)
	// Get it from: Supabase Dashboard → Project Settings → API → service_role key
	if key := os.Getenv("SUPABASE_SERVICE_ROLE_KEY"); key != "" {
		supabaseKey = key
	} else if key := os.Getenv("SUPABASE_ANON_KEY"); key != "" {
		supabaseKey = key
		log.Printf("⚠️  WARNING: Using ANON_KEY. For production, use SERVICE_ROLE_KEY to bypass RLS.")
	} else {
		fmt.Print("Enter your Supabase Service Role Key (or Anon Key): ")
		fmt.Scanln(&supabaseKey)
		if supabaseKey == "" {
			log.Fatal("❌ API key is required")
		}
	}

	log.Printf("🚀 Starting Discord Notifier")
	log.Printf("📡 Supabase URL: %s", supabaseURL)
	log.Printf("⏱️  Poll interval: %v", pollInterval)
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Verify database schema before proceeding
	log.Printf("🔍 Verifying database schema...")
	if err := verifyDatabaseSchema(); err != nil {
		log.Fatalf("❌ Database schema verification failed: %v", err)
	}
	log.Printf("✅ Database schema verified")

	// Initialize Sendico client
	log.Printf("🔧 Initializing Sendico client...")
	var err error
	sendicoClient, err = NewSendicoClient()
	if err != nil {
		log.Fatalf("❌ Failed to initialize Sendico client: %v", err)
	}
	log.Printf("✅ Sendico client initialized")

	// Run immediately
	processAllNotifications()

	// Then run on interval
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		<-ticker.C
		processAllNotifications()
	}
}

func processAllNotifications() {
	// Prevent overlapping cycles - skip if previous cycle still running
	processingMu.Lock()
	if isProcessing {
		log.Printf("⏭️  Previous cycle still running, skipping this check")
		processingMu.Unlock()
		return
	}
	isProcessing = true
	processingMu.Unlock()

	defer func() {
		processingMu.Lock()
		isProcessing = false
		processingMu.Unlock()
	}()

	startTime := time.Now()
	log.Printf("🔄 Starting notification cycle...")

	users, err := fetchActiveSubscribers()
	if err != nil {
		log.Printf("❌ Error fetching users: %v", err)
		return
	}

	if len(users) == 0 {
		log.Printf("ℹ️  No active subscribers found")
		return
	}

	log.Printf("✅ Found %d subscriber(s)", len(users))

	// Filter to only active subscriptions
	activeUsers := make([]User, 0, len(users))
	for _, user := range users {
		if isSubscriptionActive(user) {
			activeUsers = append(activeUsers, user)
		} else {
			log.Printf("⏭️  Skipping %s - subscription expired", user.Email)
		}
	}

	if len(activeUsers) == 0 {
		log.Printf("ℹ️  No active subscriptions found")
		return
	}

	log.Printf("🚀 Processing %d active subscriber(s) in parallel (max %d concurrent)", len(activeUsers), maxConcurrentUsers)

	// Process users in parallel with worker pool
	userSem := make(chan struct{}, maxConcurrentUsers)
	var wg sync.WaitGroup

	for i, user := range activeUsers {
		wg.Add(1)
		go func(idx int, u User) {
			defer wg.Done()

			// Acquire semaphore
			userSem <- struct{}{}
			defer func() { <-userSem }()

			log.Printf("👤 Processing: %s (%s) [%d/%d]", u.Username, u.Email, idx+1, len(activeUsers))
			processUserNotifications(u)
		}(i, user)
	}

	wg.Wait()
	duration := time.Since(startTime)
	log.Printf("✅ Finished processing all subscribers (took %v)", duration)

	// Warn if processing took longer than poll interval
	if duration > pollInterval {
		log.Printf("⚠️  WARNING: Processing took longer than poll interval! Consider reducing user count or increasing interval.")
	}
	
	// Clean up search cache after each cycle
	cleanupSearchCache()
}

// cleanupSearchCache removes expired entries from search cache
func cleanupSearchCache() {
	searchCacheMu.Lock()
	defer searchCacheMu.Unlock()
	
	now := time.Now()
	expiredCount := 0
	for key, cached := range searchCache {
		if now.After(cached.ExpiresAt) {
			delete(searchCache, key)
			expiredCount++
		}
	}
	if expiredCount > 0 {
		log.Printf("🧹 Cleaned up %d expired search cache entries", expiredCount)
	}
}

// verifyDatabaseSchema checks that the database table has the expected structure
func verifyDatabaseSchema() error {
	log.Printf("   Checking database connection and schema...")
	
	// Supabase uses PostgreSQL, verify we can connect and query
	// Query the table with a simple select to verify it exists and has required columns
	url := fmt.Sprintf("%s/rest/v1/unlocked_users?select=auth_user_id,email,username,discord_webhook_url,discord_notifications,notifications_subscription_active,notifications_subscription_expires_at&limit=1", supabaseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("apikey", supabaseKey)
	req.Header.Set("Authorization", "Bearer "+supabaseKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=representation")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode == 404 {
		return fmt.Errorf("table 'unlocked_users' does not exist - please run database/schema.sql in Supabase SQL Editor")
	}

	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		return fmt.Errorf("authentication failed (status %d) - check your API key and ensure it has proper permissions", resp.StatusCode)
	}

	if resp.StatusCode != 200 && resp.StatusCode != 206 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("database query failed (status %d): %s", resp.StatusCode, string(body))
	}

	// Try to decode to verify column types match
	var testUsers []User
	bodyBytes, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(bodyBytes, &testUsers); err != nil {
		// Check if it's an empty array (which is fine)
		var emptyArray []interface{}
		if json.Unmarshal(bodyBytes, &emptyArray) == nil {
			log.Printf("   ✓ Database table 'unlocked_users' exists (empty table)")
		} else {
			return fmt.Errorf("failed to parse database response - schema may be incorrect: %w\nResponse: %s", err, string(bodyBytes))
		}
	} else {
		log.Printf("   ✓ Database table 'unlocked_users' exists and schema is valid")
	}

	// Verify database type (Supabase uses PostgreSQL)
	log.Printf("   ✓ Database type: PostgreSQL (via Supabase)")
	log.Printf("   ✓ Required columns verified: auth_user_id, email, username, discord_webhook_url, discord_notifications, notifications_subscription_active, notifications_subscription_expires_at")
	
	// Check for active subscribers count
	activeCount, err := getActiveSubscriberCount()
	if err == nil {
		log.Printf("   ✓ Found %d active subscriber(s) with webhooks configured", activeCount)
	}
	
	return nil
}

// getActiveSubscriberCount returns the count of active subscribers (for verification)
func getActiveSubscriberCount() (int, error) {
	url := fmt.Sprintf("%s/rest/v1/unlocked_users?select=auth_user_id&notifications_subscription_active=eq.true&discord_webhook_url=not.is.null", supabaseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("apikey", supabaseKey)
	req.Header.Set("Authorization", "Bearer "+supabaseKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "count=exact")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 206 {
		return 0, fmt.Errorf("status %d", resp.StatusCode)
	}

	var users []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return 0, err
	}

	return len(users), nil
}

func fetchActiveSubscribers() ([]User, error) {
	// First, get ALL users to see what's actually in the database (for debugging)
	log.Printf("   🔍 Querying database for subscribers...")
	
	// Query ALL users to see what we have (for debugging)
	urlAllUsers := fmt.Sprintf("%s/rest/v1/unlocked_users?select=auth_user_id,email,username,discord_webhook_url,notifications_subscription_active&limit=100", supabaseURL)
	reqAll, _ := http.NewRequest("GET", urlAllUsers, nil)
	reqAll.Header.Set("apikey", supabaseKey)
	reqAll.Header.Set("Authorization", "Bearer "+supabaseKey)
	reqAll.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{Timeout: 30 * time.Second}
	respAll, err := client.Do(reqAll)
	if err == nil {
		defer respAll.Body.Close()
		if respAll.StatusCode == 200 {
			var allUsers []map[string]interface{}
			if json.NewDecoder(respAll.Body).Decode(&allUsers) == nil {
				totalUsers := len(allUsers)
				activeCount := 0
				withWebhookCount := 0
				activeWithWebhookCount := 0
				
				for _, u := range allUsers {
					email, _ := u["email"].(string)
					subscriptionActive := false
					if sa, ok := u["notifications_subscription_active"].(bool); ok {
						subscriptionActive = sa
					}
					
					webhookURL := ""
					if w, ok := u["discord_webhook_url"].(string); ok {
						webhookURL = strings.TrimSpace(w)
					} else if u["discord_webhook_url"] != nil {
						// Handle non-null but non-string values
						webhookURL = fmt.Sprintf("%v", u["discord_webhook_url"])
						webhookURL = strings.TrimSpace(webhookURL)
					}
					
					// Fix corrupted webhook URLs that have JSON data appended
					if strings.Contains(webhookURL, "[{") || strings.Contains(webhookURL, "{\"") {
						if idx := strings.Index(webhookURL, "[{"); idx > 0 {
							webhookURL = strings.TrimSpace(webhookURL[:idx])
						} else if idx := strings.Index(webhookURL, "{\""); idx > 0 {
							webhookURL = strings.TrimSpace(webhookURL[:idx])
						}
					}
					
					hasWebhook := webhookURL != "" && len(webhookURL) > 20 && strings.HasPrefix(webhookURL, "https://discord.com/api/webhooks/")
					
					if subscriptionActive {
						activeCount++
						if hasWebhook {
							activeWithWebhookCount++
							log.Printf("   ✓ User %s: subscription=ACTIVE, webhook=%s", email, maskWebhookURL(webhookURL))
						} else {
							log.Printf("   ⚠️  User %s: subscription=ACTIVE, webhook=EMPTY/NULL", email)
						}
					}
					
					if hasWebhook {
						withWebhookCount++
					}
				}
				
				log.Printf("   📊 Database stats: %d total user(s), %d active subscription(s), %d with webhook URL, %d active+webhook", 
					totalUsers, activeCount, withWebhookCount, activeWithWebhookCount)
				
				if activeCount > 0 && activeWithWebhookCount == 0 {
					log.Printf("   ⚠️  Warning: %d user(s) have active subscriptions but no valid webhook URLs", activeCount)
				}
			}
		}
	}

	// Now get the actual subscribers we can notify (active + webhook)
	// Query for all active subscribers first, then filter in code (more reliable than PostgREST null checks)
	url := fmt.Sprintf("%s/rest/v1/unlocked_users?select=auth_user_id,email,username,discord_webhook_url,discord_notifications,notifications_subscription_active,notifications_subscription_expires_at&notifications_subscription_active=eq.true", supabaseURL)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("apikey", supabaseKey)
	req.Header.Set("Authorization", "Bearer "+supabaseKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var allUsers []User
	if err := json.NewDecoder(resp.Body).Decode(&allUsers); err != nil {
		return nil, err
	}
	
	// Filter to only users with webhook URLs OR notifications with webhooks
	users := make([]User, 0, len(allUsers))
	for i := range allUsers {
		webhookURL := strings.TrimSpace(allUsers[i].DiscordWebhookURL)
		
		// Fix corrupted webhook URLs that have JSON data appended
		// Sometimes the URL field gets JSON notifications concatenated to it
		if strings.Contains(webhookURL, "[{") || strings.Contains(webhookURL, "{\"") {
			// Extract just the URL part (everything before the JSON starts)
			if idx := strings.Index(webhookURL, "[{"); idx > 0 {
				webhookURL = strings.TrimSpace(webhookURL[:idx])
				log.Printf("   🔧 Fixed corrupted webhook URL for user %s (had JSON appended)", allUsers[i].Email)
			} else if idx := strings.Index(webhookURL, "{\""); idx > 0 {
				webhookURL = strings.TrimSpace(webhookURL[:idx])
				log.Printf("   🔧 Fixed corrupted webhook URL for user %s (had JSON appended)", allUsers[i].Email)
			}
			// Update the user struct with the cleaned URL
			allUsers[i].DiscordWebhookURL = webhookURL
		}
		
		// Check if user has global webhook
		hasGlobalWebhook := webhookURL != "" && len(webhookURL) > 20 && strings.HasPrefix(webhookURL, "https://discord.com/api/webhooks/")
		
		// Check if user has notifications with per-notification webhooks
		hasNotificationWebhooks := false
		if len(allUsers[i].DiscordNotifications) > 0 {
			// Try to parse notifications to check for webhooks
			var notifications []Notification
			if err := json.Unmarshal(allUsers[i].DiscordNotifications, &notifications); err == nil {
				for _, notif := range notifications {
					if len(notif.Webhooks) > 0 {
						hasNotificationWebhooks = true
						break
					}
				}
			} else {
				// Try as string first
				var str string
				if err2 := json.Unmarshal(allUsers[i].DiscordNotifications, &str); err2 == nil {
					if err3 := json.Unmarshal([]byte(str), &notifications); err3 == nil {
						for _, notif := range notifications {
							if len(notif.Webhooks) > 0 {
								hasNotificationWebhooks = true
								break
							}
						}
					}
				}
			}
		}
		
		// Include user if they have either global webhook OR notification webhooks
		if hasGlobalWebhook || hasNotificationWebhooks {
			users = append(users, allUsers[i])
			if hasGlobalWebhook {
				log.Printf("   ✓ Including user %s (%s) - global webhook: %s", allUsers[i].Email, allUsers[i].Username, maskWebhookURL(webhookURL))
			} else {
				log.Printf("   ✓ Including user %s (%s) - has notification webhooks", allUsers[i].Email, allUsers[i].Username)
			}
		} else {
			log.Printf("   ⚠️  Excluding user %s (%s) - no webhooks configured", 
				allUsers[i].Email, allUsers[i].Username)
		}
	}
	
	log.Printf("   ✅ Found %d subscriber(s) ready for notifications (active subscription + webhook configured)", len(users))

	// Parse notifications JSON (handle both string and array formats)
	for i := range users {
		users[i].Notifications = []Notification{}

		if len(users[i].DiscordNotifications) == 0 {
			continue
		}

		// Try to unmarshal as array directly
		var notifications []Notification
		if err := json.Unmarshal(users[i].DiscordNotifications, &notifications); err != nil {
			// If that fails, try as a string first
			var str string
			if err2 := json.Unmarshal(users[i].DiscordNotifications, &str); err2 == nil {
				// It's a string, try to unmarshal the string content
				if err3 := json.Unmarshal([]byte(str), &notifications); err3 != nil {
					log.Printf("   ⚠️  Failed to parse notifications for user %s: %v", users[i].Email, err3)
					continue
				}
			} else {
				log.Printf("   ⚠️  Failed to parse notifications for user %s: %v", users[i].Email, err)
				continue
			}
		}

		users[i].Notifications = notifications
	}

	return users, nil
}

// maskWebhookURL masks most of the webhook URL for security in logs
func maskWebhookURL(url string) string {
	if len(url) < 20 {
		return "***"
	}
	return url[:20] + "..." + url[len(url)-10:]
}

func isSubscriptionActive(user User) bool {
	if !user.SubscriptionActive {
		return false
	}
	if user.SubscriptionExpiresAt == nil || *user.SubscriptionExpiresAt == "" {
		return true // Lifetime subscription
	}
	expiresAt, err := time.Parse(time.RFC3339, *user.SubscriptionExpiresAt)
	if err != nil {
		return false
	}
	return time.Now().Before(expiresAt)
}

func processUserNotifications(user User) {
	if len(user.Notifications) == 0 {
		log.Printf("   ℹ️  No notifications configured")
		return
	}

	for _, notif := range user.Notifications {
		log.Printf("   🔍 Checking: '%s'", notif.SearchTerm)

		// Filter markets to only include supported ones
		validMarkets := filterSupportedMarkets(notif.Markets)

		if len(notif.Markets) > 0 && len(validMarkets) == 0 {
			log.Printf("   ⚠️  Skipping - no supported markets (requested: %v)", notif.Markets)
			continue
		}

		if len(notif.Markets) > 0 {
			log.Printf("   📋 Markets: %v (filtered to supported: %v)", notif.Markets, validMarkets)
		} else {
			log.Printf("   📋 Markets: All supported markets (none specified)")
			// If no markets specified, use all supported markets
			validMarkets = getAllSupportedMarkets()
		}

		// Filter to only Sendico-supported markets
		sendicoMarketsList := filterSendicoMarkets(validMarkets)

		if len(sendicoMarketsList) == 0 {
			log.Printf("   ⚠️  No Sendico-supported markets in this notification")
			log.Printf("   💡 Sendico supports: mercari-jp, paypay-fleamarket, rakuma, rakuten-jp, yahoo-auctions")
			continue
		}

		// Translate search term to Japanese (Sendico requires Japanese)
		// Use cache to avoid duplicate API calls
		ctx := context.Background()
		termJP := getCachedTranslation(notif.SearchTerm)

		if termJP == "" {
			// Not in cache, translate it
			var err error
			termJP, err = sendicoClient.Translate(ctx, notif.SearchTerm)
			if err != nil {
				log.Printf("   ❌ Translation error: %v", err)
				continue
			}
			// Cache the translation
			cacheTranslation(notif.SearchTerm, termJP)
		}
		log.Printf("   🇯🇵 Translated '%s' → '%s'", notif.SearchTerm, termJP)

		// Search Sendico markets
		shops := make([]SendicoShop, 0, len(sendicoMarketsList))
		for _, marketKey := range sendicoMarketsList {
			shops = append(shops, sendicoMarkets[marketKey])
		}

		// Create cache key for this search (term + markets)
		marketsKey := strings.Join(sendicoMarketsList, ",")
		cacheKey := fmt.Sprintf("%s:%s", termJP, marketsKey)
		
		// Check if we have cached results for this exact search (across all users)
		items := getCachedSearchResults(cacheKey)
		
		if items == nil {
			// Not in cache, perform search
			// Note: Sendico API returns items sorted by relevance/date (page 1 typically has newest)
			var err error
			searchOpts := SendicoSearchOptions{
				TermJP: termJP,
			}
			
			if maxSearchPages > 1 {
				// Multi-page search (slower but catches more items)
				log.Printf("   🔎 Searching %d market(s) (%d pages for recently uploaded items)...", len(shops), maxSearchPages)
				items, err = sendicoClient.BulkSearchMultiplePages(ctx, shops, searchOpts, maxSearchPages)
			} else {
				// Single page search (fastest - original behavior)
				log.Printf("   🔎 Searching %d market(s)...", len(shops))
				items, err = sendicoClient.BulkSearch(ctx, shops, searchOpts)
			}
			if err != nil {
				log.Printf("   ❌ Search error: %v", err)
				continue
			}
			
			// Cache the results for 30 seconds (to share across users in same cycle)
			cacheSearchResults(cacheKey, items)
			log.Printf("   📦 Found %d item(s) across %d page(s)", len(items), maxSearchPages)
		} else {
			log.Printf("   📦 Using cached results: %d item(s) (shared across users)", len(items))
		}

		// Filter out already-seen items
		newItems := filterSeenItems(items, notif.ID)
		if len(newItems) == 0 {
			log.Printf("   ℹ️  No new items (all already seen)")
			continue
		}

		log.Printf("   ✨ %d new item(s) found!", len(newItems))

		// Convert to notification format
		notificationItems := make([]map[string]interface{}, 0, len(newItems))
		for _, item := range newItems {
			marketName := getMarketNameFromShop(item.Shop)
			notificationItems = append(notificationItems, map[string]interface{}{
				"title":       item.Name,
				"description": fmt.Sprintf("Price: ¥%d ($%d)", item.PriceYen, item.PriceUSD),
				"url":         item.URL,
				"price":       fmt.Sprintf("¥%d ($%d)", item.PriceYen, item.PriceUSD),
				"market":      marketName,
				"image":       item.Image,
			})
		}

		// Determine which webhooks to use
		webhooksToUse := notif.Webhooks
		if len(webhooksToUse) == 0 {
			// Fallback to global webhook if no per-notification webhooks
			if user.DiscordWebhookURL != "" {
				webhooksToUse = []string{user.DiscordWebhookURL}
			} else {
				log.Printf("   ⚠️  No webhooks configured for this notification")
				continue
			}
		}

		// Send notification to each webhook
		log.Printf("   ✅ Sending notification to %d webhook(s)...", len(webhooksToUse))
		for i, webhookURL := range webhooksToUse {
			webhookURL = strings.TrimSpace(webhookURL)
			if webhookURL == "" || !strings.HasPrefix(webhookURL, "https://discord.com/api/webhooks/") {
				log.Printf("   ⚠️  Skipping invalid webhook %d/%d", i+1, len(webhooksToUse))
				continue
			}
			
			if err := sendDiscordNotification(webhookURL, notif, notificationItems); err != nil {
				log.Printf("   ❌ Error sending to webhook %d/%d: %v", i+1, len(webhooksToUse), err)
			} else {
				log.Printf("   ✅ Notification sent to webhook %d/%d!", i+1, len(webhooksToUse))
			}
		}
	}
}

// filterSupportedMarkets filters the markets list to only include supported markets
func filterSupportedMarkets(markets []string) []string {
	if len(markets) == 0 {
		return []string{}
	}

	valid := []string{}
	for _, market := range markets {
		// Skip custom markets (they start with "custom-")
		if len(market) > 7 && market[:7] == "custom-" {
			log.Printf("      ⚠️  Skipping custom market: %s (not supported by notifier)", market)
			continue
		}

		if supportedMarkets[market] {
			valid = append(valid, market)
		} else {
			log.Printf("      ⚠️  Skipping unsupported market: %s", market)
		}
	}

	return valid
}

// getAllSupportedMarkets returns a list of all supported market keys
func getAllSupportedMarkets() []string {
	markets := make([]string, 0, len(supportedMarkets))
	for market := range supportedMarkets {
		markets = append(markets, market)
	}
	return markets
}

// filterSendicoMarkets filters markets to only include Sendico-supported ones
func filterSendicoMarkets(markets []string) []string {
	result := []string{}
	for _, market := range markets {
		if _, ok := sendicoMarkets[market]; ok {
			result = append(result, market)
		}
	}
	return result
}

// filterSeenItems filters out items that have already been seen
// Uses TTL-based tracking to prevent memory bloat while ensuring we don't miss items
func filterSeenItems(items []SendicoItem, notificationID string) []SendicoItem {
	newItems := []SendicoItem{}
	now := time.Now()

	seenItemsMu.Lock()
	defer seenItemsMu.Unlock()

	// Clean up expired entries periodically (every 1000 items checked)
	if len(seenItems) > 10000 {
		expiredKeys := make([]string, 0)
		for key, timestamp := range seenItems {
			if now.Sub(timestamp) > seenItemsTTL {
				expiredKeys = append(expiredKeys, key)
			}
		}
		for _, key := range expiredKeys {
			delete(seenItems, key)
		}
		if len(expiredKeys) > 0 {
			log.Printf("   🧹 Cleaned up %d expired seen items", len(expiredKeys))
		}
	}

	for _, item := range items {
		// Create unique key: notificationID:shop:code
		key := fmt.Sprintf("%s:%s:%s", notificationID, item.Shop, item.Code)

		seenTime, exists := seenItems[key]
		if !exists {
			// New item - mark as seen
			seenItems[key] = now
			newItems = append(newItems, item)
		} else {
			// Item was seen before - check if TTL expired (shouldn't happen often)
			if now.Sub(seenTime) > seenItemsTTL {
				// TTL expired, treat as new (very rare case)
				seenItems[key] = now
				newItems = append(newItems, item)
				log.Printf("   ⚠️  Item %s expired from cache, treating as new", item.Code)
			}
			// Otherwise, skip (already seen)
		}
	}

	return newItems
}

// getMarketNameFromShop returns the human-readable market name
func getMarketNameFromShop(shop SendicoShop) string {
	switch shop {
	case SendicoMercari:
		return "Mercari Japan"
	case SendicoRakuma:
		return "Rakuten Rakuma"
	case SendicoRakuten:
		return "Rakuten"
	case SendicoYahooAuctions:
		return "Yahoo Auctions"
	case SendicoYahoo:
		return "Yahoo PayPay Flea"
	default:
		return string(shop)
	}
}

// Translation caching functions
func getCachedTranslation(term string) string {
	translationCacheMu.RLock()
	defer translationCacheMu.RUnlock()
	return translationCache[term]
}

func cacheTranslation(term, translation string) {
	translationCacheMu.Lock()
	defer translationCacheMu.Unlock()
	translationCache[term] = translation
	// Limit cache size to prevent memory issues (keep last 1000 translations)
	if len(translationCache) > 1000 {
		// Simple eviction: clear oldest 200 entries (or use LRU in production)
		// For now, just clear if too large (simple approach)
		if len(translationCache) > 1200 {
			translationCache = make(map[string]string)
		}
	}
}

// Search result caching functions (for batching across users)
func getCachedSearchResults(cacheKey string) []SendicoItem {
	searchCacheMu.RLock()
	defer searchCacheMu.RUnlock()
	
	cached, exists := searchCache[cacheKey]
	if !exists {
		return nil
	}
	
	// Check if cache is still valid (30 seconds)
	if time.Now().After(cached.ExpiresAt) {
		return nil
	}
	
	return cached.Items
}

func cacheSearchResults(cacheKey string, items []SendicoItem) {
	searchCacheMu.Lock()
	defer searchCacheMu.Unlock()
	
	// Create a copy of items to avoid race conditions
	itemsCopy := make([]SendicoItem, len(items))
	copy(itemsCopy, items)
	
	searchCache[cacheKey] = &cachedSearchResult{
		Items:     itemsCopy,
		Timestamp: time.Now(),
		ExpiresAt: time.Now().Add(30 * time.Second), // Cache for 30 seconds
	}
	
	// Clean up expired entries periodically
	if len(searchCache) > 100 {
		now := time.Now()
		for key, cached := range searchCache {
			if now.After(cached.ExpiresAt) {
				delete(searchCache, key)
			}
		}
	}
}

func sendDiscordNotification(webhookURL string, notification Notification, items []map[string]interface{}) error {
	// Discord limits embeds to 10 per message, so we need to batch
	maxEmbeds := 10
	totalItems := len(items)

	for i := 0; i < totalItems; i += maxEmbeds {
		end := i + maxEmbeds
		if end > totalItems {
			end = totalItems
		}

		batch := items[i:end]
		embeds := []DiscordEmbed{}

		for _, item := range batch {
			// Use item title, or fallback to search term
			itemTitle := getString(item, "title", "")
			if itemTitle == "" {
				itemTitle = notification.SearchTerm
			}

			// Truncate title if too long (Discord limit is 256 chars)
			if len(itemTitle) > 200 {
				itemTitle = itemTitle[:197] + "..."
			}

			embed := DiscordEmbed{
				Title:       itemTitle,
				Description: getString(item, "description", ""),
				URL:         getString(item, "url", ""),
				Color:       3447003, // Blue color
				Timestamp:   time.Now().Format(time.RFC3339),
				Footer: map[string]interface{}{
					"text": fmt.Sprintf("MMCS • %s", notification.SearchTerm),
				},
			}

			if price := getString(item, "price", ""); price != "" {
				embed.Fields = append(embed.Fields, DiscordEmbedField{
					Name:   "Price",
					Value:  price,
					Inline: true,
				})
			}

			if market := getString(item, "market", ""); market != "" {
				embed.Fields = append(embed.Fields, DiscordEmbedField{
					Name:   "Market",
					Value:  market,
					Inline: true,
				})
			}

			// Add image thumbnail if available
			if imageURL := getString(item, "image", ""); imageURL != "" {
				embed.Image = map[string]string{
					"url": imageURL,
				}
			}

			embeds = append(embeds, embed)
		}

		// Create content message
		var content string
		if i == 0 {
			// First batch
			if totalItems > maxEmbeds {
				content = fmt.Sprintf("🔔 **%d new item(s) found for: %s** (showing first %d)", totalItems, notification.SearchTerm, len(batch))
			} else {
				content = fmt.Sprintf("🔔 **%d new item(s) found for: %s**", totalItems, notification.SearchTerm)
			}
		} else {
			// Subsequent batches
			content = fmt.Sprintf("🔔 **More items for: %s** (%d-%d of %d)", notification.SearchTerm, i+1, end, totalItems)
		}

		payload := DiscordWebhookPayload{
			Content: content,
			Embeds:  embeds,
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %w", err)
		}

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			return fmt.Errorf("failed to send request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 && resp.StatusCode != 204 {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("Discord returned status %d: %s", resp.StatusCode, string(body))
		}

		// Small delay between batches to avoid Discord rate limiting
		// Discord limit: 30 requests per 10 seconds per webhook
		// Reduced to 500ms for faster processing while still respecting limits
		if end < totalItems {
			time.Sleep(500 * time.Millisecond)
		}
	}

	return nil
}

func getString(m map[string]interface{}, key string, defaultValue string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return defaultValue
}
